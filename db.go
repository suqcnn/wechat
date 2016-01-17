package wechat

import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var Orm *xorm.Engine

func InitDB(dbtype string, dataSourceName string, objs ...interface{}) error {

	var err error
	//创建表
	Orm, err = xorm.NewEngine(dbtype, dataSourceName)
	if err != nil {
		return err
	}

	Orm.ShowSQL = true
	err = Orm.Sync2(objs...)
	if err != nil {
		return err
	}

	return nil
}
