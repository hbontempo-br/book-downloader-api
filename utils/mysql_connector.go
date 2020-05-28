package utils

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wantedly/gorm-zap"
	"go.uber.org/zap"
	"time"
)

func NewMySQLConnector(address string, port int, dbName string, user string, password string) MySQLConnector {
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
	msc.db.LogMode(true)
	msc.db.SetLogger(gormzap.New(zap.L()))
	msc.db.DB().SetMaxIdleConns(10)
	msc.db.DB().SetMaxOpenConns(100)
	msc.db.DB().SetConnMaxLifetime(time.Hour)

	return msc.db, nil
}

func (msc MySQLConnector) connectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True", msc.user, msc.password, msc.address, msc.port, msc.dbName)
}
