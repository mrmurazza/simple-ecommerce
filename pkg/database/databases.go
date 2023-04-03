package database

import (
	"ecommerce/config"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var DB *gorm.DB

func InitDatabase() {
	cfg := config.Get()

	dbResourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	if cfg.DBDriver == "sqlite3" {
		dbResourceName = cfg.DBHost
	}

	db, err := gorm.Open(cfg.DBDriver, dbResourceName)
	if err != nil {
		panic(err)
	}

	DB = db

	// set gorm configuration
	DB.LogMode(true)
	DB.SingularTable(false)

	err = DB.Exec("CREATE TABLE IF NOT EXISTS customers (" +
		"`id` integer NOT NULL primary key, " +
		"email varchar(50) not null default ``, " +
		"name varchar(50) not null default ``, " +
		"password varchar(50) not null default ``, " +
		"role varchar(50) not null, " +
		"created_at datetime not null default current_timestamp, " +
		"updated_at datetime not null default current_timestamp " +
		")").Error

	if err != nil {
		panic(err)
	}

	err = DB.Exec("CREATE TABLE IF NOT EXISTS products (" +
		"`id` integer NOT NULL primary key, " +
		"name varchar(50) not null default ``, " +
		"price integer not null default ``, " +
		"description varchar(50) not null default ``, " +
		"image varchar(50) not null default ``, " +
		"created_at datetime not null default current_timestamp, " +
		"updated_at datetime not null default current_timestamp " +
		")").Error

	if err != nil {
		panic(err)
	}

	err = DB.Exec("CREATE TABLE IF NOT EXISTS orders (" +
		"`id` integer NOT NULL primary key, " +
		"customer_id integer not null , " +
		"status varchar(50) not null default ``, " +
		"total_qty integer not null, " +
		"total_price integer not null, " +
		"created_at datetime not null default current_timestamp, " +
		"updated_at datetime not null default current_timestamp, " +
		"FOREIGN KEY (customer_id) REFERENCES customers(id)" +
		")").Error

	if err != nil {
		panic(err)
	}

	err = DB.Exec("CREATE TABLE IF NOT EXISTS order_units (" +
		"`id` integer NOT NULL primary key, " +
		"order_id integer not null, " +
		"product_id integer not null, " +
		"qty integer not null default 1, " +
		"price integer not null default 0, " +
		"name varchar(50) not null default ``, " +
		"description varchar(50) not null default ``, " +
		"image varchar(50) not null default ``, " +
		"created_at datetime not null default current_timestamp, " +
		"FOREIGN KEY (order_id) REFERENCES orders(id)" +
		")").Error

	if err != nil {
		panic(err)
	}

}
