package service

import (
	"gdiamond/server/internal/model"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

//PageSize  page size dump to disk
const PageSize = 1000

//Init schedule dump disk task
func SetupDumpTask() error {
	err := DumpAll2Disk()
	if err != nil {
		return err
	}

	log.Println("Finished first dump config from database to disk")
	Logger.WithFields(logrus.Fields{}).Info("Finished first dump config from database to disk")

	ticker := time.NewTicker(time.Minute * 2)
	go func() {
		defer ticker.Stop()
		for {
			<-ticker.C
			err := DumpAll2Disk()
			if err != nil {
				continue
			}
		}
	}()
	return nil
}

//DumpAll2Disk  dump all config info from database to disk
func DumpAll2Disk() error {
	page, err := findAllConfigInfo(1, PageSize)
	if err != nil {
		return err
	}

	if page != nil {
		totalPages := page.PageAvailable
		updateConfigInfo2CacheAndDisk(page)
		if totalPages > 1 {
			for pageNo := 2; pageNo <= totalPages; pageNo++ {
				page, err := findAllConfigInfo(pageNo, PageSize)
				if err != nil {
					Logger.WithFields(logrus.Fields{
						"err":    err.Error(),
						"pageNo": pageNo,
					}).Warn("findAllConfigInfo failed")
					return err
				}
				if page != nil {
					updateConfigInfo2CacheAndDisk(page)
				}
			}
		}
	}
	return nil
}

func updateConfigInfo2CacheAndDisk(page *model.Page) error {
	for _, item := range page.PageItems {
		if item == nil {
			continue
		}
		configInfo, ok := item.(*model.ConfigInfo)
		if ok {
			UpdateMD5Cache(configInfo)
			err := SaveToDisk(configInfo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
