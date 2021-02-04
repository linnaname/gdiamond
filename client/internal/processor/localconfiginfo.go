package processor

import (
	"gdiamond/client/internal/configinfo"
	"gdiamond/util/fileutil"
	"gdiamond/util/filewatcher"
	"github.com/fsnotify/fsnotify"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

//LocalConfigInfoProcessor local file processor
type LocalConfigInfoProcessor struct {
	sync.Mutex
	rootPath   string
	isRun      bool
	existFiles map[string] /*filePath*/ int64 /*version*/
}

//BaseDir local file base dir
const BaseDir = "config-data"

//NewLocalConfigInfoProcessor new
func NewLocalConfigInfoProcessor() *LocalConfigInfoProcessor {
	p := &LocalConfigInfoProcessor{isRun: false, existFiles: make(map[string]int64)}
	return p
}

//Start setup processor and create local file dir
func (p *LocalConfigInfoProcessor) Start(rootPath string) {
	p.Lock()
	defer p.Unlock()
	if p.isRun {
		return
	}
	p.rootPath = rootPath
	p.isRun = true

	initDataDir(p.rootPath)
	p.startCheckLocalDir(p.rootPath)

}

//Stop stop processor
func (p *LocalConfigInfoProcessor) Stop() {
	p.Lock()
	if !p.isRun {
		return
	}
	p.isRun = false
	p.Unlock()
}

//GetLocalConfigureInformation get config info from local file
func (p *LocalConfigInfoProcessor) GetLocalConfigureInformation(cacheData *configinfo.CacheData, force bool) (string, error) {
	filePath := p.getFilePath(cacheData.DataId(), cacheData.Group())
	_, ok := p.existFiles[filePath]
	if !ok {
		if cacheData.UseLocalConfigInfo() {
			cacheData.SetLastModifiedHeader("")
			cacheData.SetMD5("")
			cacheData.SetLocalConfigInfoFile("")
			cacheData.SetLocalConfigInfoVersion(int64(0))
			cacheData.SetUseLocalConfigInfo(false)
		}
		return "", nil
	}

	if force {
		log.Println("主动从本地获取配置数据, dataId:" + cacheData.DataId() + ", group:" + cacheData.Group())
		return fileutil.GetFileContent(filePath)
	}

	// 判断是否变更，没有变更，返回null
	if filePath == cacheData.GetLocalConfigInfoFile() || p.existFiles[filePath] != cacheData.GetLocalConfigInfoVersion() {
		content, err := fileutil.GetFileContent(filePath)
		if err != nil {
			return "", err
		}
		cacheData.SetLocalConfigInfoFile(filePath)
		cacheData.SetLocalConfigInfoVersion(p.existFiles[filePath])
		cacheData.SetUseLocalConfigInfo(true)
		log.Println("本地配置数据发生变化, dataId:" + cacheData.DataId() + ", group:" + cacheData.Group())
		return content, nil
	}
	cacheData.SetUseLocalConfigInfo(true)
	log.Println("本地配置数据没有发生变化, dataId:" + cacheData.DataId() + ", group:" + cacheData.Group())
	return "", nil
}

func initDataDir(rootPath string) {
	fileutil.CreateDirIfNecessary(rootPath)
}

func (p *LocalConfigInfoProcessor) startCheckLocalDir(filePath string) error {
	rw, err := filewatcher.NewWatcher()
	if err != nil {
		return err
	}
	err = rw.AddRecursive(filePath)
	if err != nil {
		return err
	}
	//TODO 对于已经创建的fsnotify是无能为力的，是不是要把已经创建的文件和目录做一次创建事件的恢复？这里是不是要先做一次主动的check?
	go func() {
		for {
			select {
			case e := <-rw.Events:
				p.processEvents(e)
			case e := <-rw.Errors:
				log.Println("file watch err:", e)
			}
		}
	}()
	return nil
}

func (p *LocalConfigInfoProcessor) processEvents(e fsnotify.Event) {
	grandpaDir, err := fileutil.GetGrandpaDir(e.Name)
	log.Println("file watch GetGrandpaDir:", err)

	if e.Op&fsnotify.Create == fsnotify.Create || e.Op&fsnotify.Write == fsnotify.Write {
		log.Println("创建或写入文件 : ", e.Name)
		if BaseDir != grandpaDir {
			log.Println("无效的文件进入监控目录: " + e.Name)
			return
		}
		p.existFiles[e.Name] = time.Now().Unix()
	} else if e.Op&fsnotify.Remove == fsnotify.Remove {
		log.Println("delete file : ", e.Name)
		if BaseDir == grandpaDir {
			// 删除的是文件
			delete(p.existFiles, e.Name)
		} else {
			// 删除的是目录
			if fileutil.IsDir(e.Name) {
				for k := range p.existFiles {
					if strings.HasPrefix(k, e.Name) {
						delete(p.existFiles, k)
					}
				}
			}
		}
	} else if e.Op&fsnotify.Rename == fsnotify.Rename {
		log.Println("rename file : ", e.Name)
	}
}

func (p *LocalConfigInfoProcessor) getFilePath(dataId, group string) string {
	filePathBuilder := strings.Builder{}
	filePathBuilder.WriteString(p.rootPath)
	filePathBuilder.WriteString("/")
	filePathBuilder.WriteString(BaseDir)
	filePathBuilder.WriteString("/")
	filePathBuilder.WriteString(group)
	filePathBuilder.WriteString("/")
	filePathBuilder.WriteString(dataId)
	abs, _ := filepath.Abs(filePathBuilder.String())
	return abs
}
