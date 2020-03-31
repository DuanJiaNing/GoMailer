package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"

	"GoMailer/app"
	"GoMailer/conf"
	"GoMailer/log"
)

var (
	engine *xorm.Engine
)

func NewClient() (*xorm.Engine, error) {
	if engine == nil {
		err := prepareEngine()
		if err != nil {
			return nil, err
		}
	}

	err := engine.Ping()
	if err != nil {
		err := engine.Close()
		if err != nil {
			return nil, err
		}
		log.Warningf("xorm engine closed due to ping failed, will recreate soon.")

		// Recreate engine.
		err = prepareEngine()
		if err != nil {
			return nil, err
		}
	}
	return engine, nil
}

func prepareEngine() error {
	var err error
	engine, err = xorm.NewEngine("mysql", conf.DataSource().URL)
	if err != nil {
		log.Errorf("fail to create db engine: %v", err)
		return err
	}
	err = engine.Ping()
	if err != nil {
		log.Errorf("fail to ping db: %v", err)
		return err
	}

	engine.SetMapper(core.SnakeMapper{})

	engine.ShowSQL(app.IsDevAppServer())
	engine.ShowExecTime(true)

	return nil
}
