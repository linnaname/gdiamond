package filewatcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"
)

const (
	TestDir = "test"
)

func TestRWatcher(t *testing.T) {
	rw, err := NewWatcher()
	assert.NoError(t, err)
	assert.NotNil(t, rw)
	err = rw.AddRecursive(TestDir)
	assert.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case e := <-rw.fsnotify.Events:
				if e.Op&fsnotify.Create == fsnotify.Create {
					log.Println("创建文件 : ", e.Name)
				}
				if e.Op&fsnotify.Write == fsnotify.Write {
					log.Println("写入文件 : ", e.Name)
				}
				if e.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("删除文件 : ", e.Name)
				}
				if e.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("重命名文件 : ", e.Name)
				}
				if e.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Println("修改权限 : ", e.Name)
				}
			case e := <-rw.fsnotify.Errors:
				log.Println("err:", e)
			}
		}
		wg.Done()
	}()
	wg.Wait()
}
