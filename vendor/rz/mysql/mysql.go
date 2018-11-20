package mysql

import (
	"rz/config"
	"rz/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var instance *gorm.DB

func InitMysql() {
	var err error
	var connectStr string
	if (config.IS_TEST) {
		connectStr = config.TEST_MYSQL_USER + ":" + config.TEST_MYSQL_PASSWORD+"@tcp("+config.TEST_MYSQL_HOST+")/"+config.TEST_MYSQL_DB+"?charset=utf8&parseTime=True&loc=Local"
	} else {
		connectStr = config.MYSQL_USER + ":" + config.MYSQL_PASSWORD+"@tcp("+config.MYSQL_HOST+")/"+config.MYSQL_DB+"?charset=utf8&parseTime=True&loc=Local"
	}
	instance, err = gorm.Open("mysql", connectStr)
	util.CheckErr(err)
}

func Instance() *gorm.DB {
	if instance == nil {
		InitMysql()
	}
	return instance
}