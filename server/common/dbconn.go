package common

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var GDBConn *sql.DB

func openDBConn(dbUrl string, maxIdleConns, maxOpenConns int) (err error) {
	if GDBConn, err = sql.Open("mysql", dbUrl); err != nil {
		return err
	}
	//Configuring sql.DB for Better Performance:https://www.alexedwards.net/blog/configuring-sqldb
	GDBConn.SetMaxIdleConns(maxIdleConns)
	GDBConn.SetMaxOpenConns(maxOpenConns)

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

func CloseConn() {
	if GDBConn != nil {
		GDBConn.Close()
	}
}
