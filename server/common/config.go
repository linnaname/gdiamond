package common

import (
	"encoding/json"
	"io/ioutil"
)

var GMySQLConfig *MySQLConfig

type MySQLConfig struct {
	DBUrl        string `json:"dbUrl"`
	MaxIdleConns int    `json:"maxIdleConns"`
	MaxOpenConns int    `json:"maxOpenConns"`
}

const CONFIG_PATH = "/Users/goranka/linnana/go/gdiamond/server/etc/etc.json"

/**
read mysql connect config from local file
*/
func InitConfig() error {
	content, err := ioutil.ReadFile(CONFIG_PATH)
	if err != nil {
		return err
	}
	GMySQLConfig = &MySQLConfig{}
	err = json.Unmarshal(content, GMySQLConfig)
	if err != nil {
		return err
	}
	return nil
}
