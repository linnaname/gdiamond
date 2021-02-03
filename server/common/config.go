package common

import (
	"flag"
	"gdiamond/common/namesrv"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/spf13/viper"
	"strings"
)

//GMySQLConfig  mysql sql config holder
var GMySQLConfig *MySQLConfig

//RegisterRequestConfig name server register config
var RegisterRequestConfig *namesrv.RegisterRequest

//NameServerAddressList name server address from flag args
var NameServerAddressList *singlylinkedlist.List

//MySQLConfig struct of mysql config
type MySQLConfig struct {
	DBUrl        string
	MaxIdleConns int
	MaxOpenConns int
}

//InitConfig read mysql connect config from local file
func InitConfig() error {
	nameSrvAdders := flag.String("n", "", "name server address,spe with ; when more than one")
	c := flag.String("c", "", "config/etc")
	if !flag.Parsed() {
		flag.Parse()
	}
	configPath := *c
	v := viper.New()
	v.SetConfigName("gdiamond")
	v.SetConfigType("toml")
	v.AddConfigPath(configPath)
	err := v.ReadInConfig()
	if err != nil {
		return err
	}
	mysqlViper := v.Sub("mysql")
	GMySQLConfig = &MySQLConfig{}
	err = mysqlViper.Unmarshal(GMySQLConfig)
	if err != nil {
		return err
	}
	serverViper := v.Sub("server")
	RegisterRequestConfig = &namesrv.RegisterRequest{}
	err = serverViper.Unmarshal(RegisterRequestConfig)
	if err != nil {
		return err
	}
	NameServerAddressList = getNameServerAddressList(nameSrvAdders)
	return nil
}

//getNameServerAddressList get name server address from flag
func getNameServerAddressList(nameSrvAdders *string) *singlylinkedlist.List {
	if nameSrvAdders != nil {
		if *nameSrvAdders == "" {
			return nil
		}
		nameSrvAdderArr := strings.Split(*nameSrvAdders, ";")
		if len(nameSrvAdderArr) == 0 {
			return nil
		}
		nameServerAddressList := singlylinkedlist.New()
		for _, nameServerAddress := range nameSrvAdderArr {
			nameServerAddressList.Add(nameServerAddress)
		}
		return nameServerAddressList
	}
	return nil
}
