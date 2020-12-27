package model

import (
	"log"
	"strings"
	"time"

	"github.com/hprose/hprose-golang/io"
	"github.com/jinzhu/gorm"
	// for db SQL
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"snack.com/xiyanxiyan10/stocktrader/config"
)

var (
	// DB Database
	DB *gorm.DB
)

func InitModel() {
	io.Register((*User)(nil), "User", "json")
	io.Register((*Exchange)(nil), "Exchange", "json")
	io.Register((*Algorithm)(nil), "Algorithm", "json")
	io.Register((*Trader)(nil), "Trader", "json")
	io.Register((*Log)(nil), "Log", "json")

	dbType := config.String("dbtype")
	dbURL := config.String("dburl")
	var err error
	DB, err = gorm.Open(strings.ToLower(dbType), dbURL)
	if err != nil {
		log.Panicf("Connect to %v database error: %v\n", dbType, err)
		dbType = "sqlite3"
		dbURL = "custom/data.db"
		DB, err = gorm.Open(dbType, dbURL)
		if err != nil {
			log.Panicln("Connect to database error:", err)
		}
	}
	DB.AutoMigrate(&User{}, &Exchange{}, &Algorithm{}, &TraderExchange{}, &Trader{}, &Log{})
	users := []User{}
	DB.Find(&users)
	if len(users) == 0 {
		admin := User{
			Username: "admin",
			Password: "admin",
			Level:    99,
		}
		if err := DB.Create(&admin).Error; err != nil {
			log.Fatalln("Create admin error:", err)
		}
	}
	DB.LogMode(false)
	go ping()
}

func ping() {

	dbType := config.String("dbtype")
	dbURL := config.String("dburl")
	for {
		if err := DB.Exec("SELECT 1").Error; err != nil {
			log.Println("Database ping error:", err)
			if DB, err = gorm.Open(strings.ToLower(dbType), dbURL); err != nil {
				log.Println("Retry connect to database error:", err)
			}
		}
		time.Sleep(time.Minute)
	}
}

// NewOrm ...
func NewOrm() (*gorm.DB, error) {
	dbType := config.String("dbtype")
	dbURL := config.String("dburl")
	return gorm.Open(strings.ToLower(dbType), dbURL)
}
