package common

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var GDBConn *gorm.DB

func openDBConn(dbUrl string, maxIdleConns, maxOpenConns int) (err error) {
	if GDBConn, err = gorm.Open("mysql", dbUrl); err != nil {
		return err
	}
	//Configuring sql.DB for Better Performance:https://www.alexedwards.net/blog/configuring-sqldb
	GDBConn.DB().SetMaxIdleConns(maxIdleConns)
	GDBConn.DB().SetMaxOpenConns(maxOpenConns)

	GDBConn.LogMode(true)
	GDBConn.SingularTable(true)

	return nil
}

/**
just init open db
*/
func InitDBConn() error {
	err := openDBConn(GMySQLConfig.DBUrl, GMySQLConfig.MaxIdleConns, GMySQLConfig.MaxOpenConns)
	if err != nil {
		return err
	}
	return nil
}
