package service

import (
	"gdiamond/server/model"
	"log"
	"time"
)

const PAGE_SIZE = 1000

func Init() {
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
}

func DumpAll2Disk() error {
	page, err := findAllConfigInfo(1, PAGE_SIZE)
	if err != nil {
		log.Println("errs:", err)
		return err
	}

	if page != nil {
		totalPages := page.PageAvailable
		updateConfigInfo2CacheAndDisk(page)
		if totalPages > 1 {
			for pageNo := 2; pageNo <= totalPages; pageNo++ {
				page, err := findAllConfigInfo(pageNo, PAGE_SIZE)
				if err != nil {
					log.Println("errs:", err)
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
