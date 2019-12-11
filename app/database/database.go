package database

import (
	"github.com/RedHatInsights/vmaas-go/app/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	Db    *gorm.DB
)

func Configure() {
	db, err := gorm.Open("sqlite3", config.SQLiteFilePath)
	if err != nil {
		panic(err)
	}
	Db = db
}
