package utils

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql package required for gorm
	"github.com/wantedly/gorm-zap"
	"go.uber.org/zap"
)

func NewMySQLConnector(address string, port int, dbName, user, password string) MySQLConnector {
	return MySQLConnector{address: address, port: port, dbName: dbName, user: user, password: password}
}

type MySQLConnector struct {
	address  string
	port     int
	dbName   string
	user     string
	password string
	db       *gorm.DB
}

func (msc *MySQLConnector) Connect() (*gorm.DB, error) {
	zap.S().Debug("Connecting to DB")
	var err error
	msc.db, err = gorm.Open("mysql", msc.connectionString())
	if err != nil {
		zap.S().Error("Generic error on MySQLConnector.Connect", err)
		return nil, err
	}

	// Config stuff
	msc.db.LogMode(true) // TODO: remove this magic number, use environment variable
	msc.db.SetLogger(gormzap.New(zap.L()))
	msc.db.DB().SetMaxIdleConns(10)           // TODO: remove this magic number, use environment variable
	msc.db.DB().SetMaxOpenConns(100)          // TODO: remove this magic number, use environment variable
	msc.db.DB().SetConnMaxLifetime(time.Hour) // TODO: remove this magic number, use environment variable

	return msc.db, nil
}

func (msc MySQLConnector) connectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True", msc.user, msc.password, msc.address, msc.port, msc.dbName)
}
