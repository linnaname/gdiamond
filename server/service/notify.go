package service

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	URL_PREFIX = "/diamond-server/notify"
	PROTOCOL   = "http://"
)

func notifyConfigInfoChange(dataId, group string) {
	nodes := getNodeAddress()
	for _, addr := range nodes {
		urlString := generateNotifyConfigInfoPath(dataId, group, addr)
		result, err := invokeURL(urlString)
		log.Println("notify node and result", addr, result, err)
	}
}

func generateNotifyConfigInfoPath(dataId, group, address string) string {
	urlString := PROTOCOL + address + URL_PREFIX
	urlString += "?method=notifyConfigInfo&dataId=" + dataId + "&group=" + group
	return urlString
}

func getNodeAddress() []string {
	//TODO get node address from namesrv
	return []string{"127.0.0.1:1210"}
}

func invokeURL(urlString string) (string, error) {
	client := &http.Client{Timeout: time.Minute * 5}
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
