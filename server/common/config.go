package common

import (
	"encoding/json"
	"io/ioutil"
)

//GMySQLConfig  mysql sql config holder
var GMySQLConfig *MySQLConfig

//MySQLConfig struct of mysql config
type MySQLConfig struct {
	DBUrl        string `json:"dbUrl"`
	MaxIdleConns int    `json:"maxIdleConns"`
	MaxOpenConns int    `json:"maxOpenConns"`
}

const configPath = "/Users/goranka/linnana/go/gdiamond/server/etc/etc.json"

//InitConfig read mysql connect config from local file
func InitConfig() error {
	content, err := ioutil.ReadFile(configPath)
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
