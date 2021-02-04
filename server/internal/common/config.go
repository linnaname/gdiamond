package common

import (
	"errors"
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

//ParseCmdAndInitConfig read config file path from cmd
//and read config from file
func ParseCmdAndInitConfig() error {
	nameSrvAdders := flag.String("n", "", "name server address,spe with ; when more than one")
	c := flag.String("c", "", "need a config file dir where contain a config file name gdiamond.toml")
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
	err = setupGMySQLConfig(v)
	if err != nil {
		return err
	}
	err = setupRegisterRequestConfig(v)
	if err != nil {
		return err
	}
	//parse name server list from cmd
	return setupNameServerAddressList(nameSrvAdders)
}

func setupGMySQLConfig(v *viper.Viper) error {
	mysqlViper := v.Sub("mysql")
	if mysqlViper == nil {
		return errors.New("can't find mysql config")
	}
	GMySQLConfig = &MySQLConfig{}
	err := mysqlViper.Unmarshal(GMySQLConfig)
	return err
}

func setupRegisterRequestConfig(v *viper.Viper) error {
	serverViper := v.Sub("server")
	RegisterRequestConfig = &namesrv.RegisterRequest{}
	err := serverViper.Unmarshal(RegisterRequestConfig)
	return err
}

func setupNameServerAddressList(nameSrvAdders *string) error {
	NameServerAddressList = getNameServerAddressListFromCmd(nameSrvAdders)
	if NameServerAddressList == nil || NameServerAddressList.Size() == 0 {
		return errors.New("can't get name server address from cmd")
	}
	return nil
}

//getNameServerAddressList get name server address from flag
func getNameServerAddressListFromCmd(nameSrvAdders *string) *singlylinkedlist.List {
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
