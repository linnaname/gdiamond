package service

import (
	"gdiamond/server/common"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	notifyURL = "/diamond-server/notify"
	protocol  = "http://"
)

func notifyConfigInfoChange(dataID, group string) {
	nodes := getNodeAddress()
	for _, addr := range nodes {
		urlString := generateNotifyConfigInfoPath(dataID, group, addr)
		result, err := invokeURL(urlString)
		log.Println("notify node and result", addr, result, err)
	}
}

func generateNotifyConfigInfoPath(dataID, group, address string) string {
	urlString := protocol + address + notifyURL
	urlString += "?dataId=" + dataID + "&group=" + group
	return urlString
}

func getNodeAddress() []string {
	nameServerAddressList := common.NameServerAddressList
	if nameServerAddressList == nil || nameServerAddressList.Empty() {
		return []string{"127.0.0.1:1210"}
	}
	nodeAddress := make([]string, nameServerAddressList.Size())
	for i := range nodeAddress {
		value, _ := nameServerAddressList.Get(i)
		nodeAddress[i] = value.(string) + ":1210"
	}
	return nodeAddress
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
