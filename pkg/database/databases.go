package database

import (
	"ecommerce/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
	cfg := config.Get()

	db, err := gorm.Open(sqlite.Open(cfg.DBHost), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	DB = db

	err = DB.Exec("CREATE TABLE IF NOT EXISTS users (" +
		"`id` integer NOT NULL primary key AUTOINCREMENT, " +
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
		"`id` integer NOT NULL primary key AUTOINCREMENT, " +
		"name varchar(50) not null default ``, " +
		"price integer not null default ``, " +
		"qty integer not null default ``, " +
		"description varchar(50) not null default ``, " +
		"image varchar(50) not null default ``, " +
		"created_at datetime not null default current_timestamp, " +
		"updated_at datetime not null default current_timestamp " +
		")").Error

	if err != nil {
		panic(err)
	}

	err = DB.Exec("CREATE TABLE IF NOT EXISTS orders (" +
		"`id` integer NOT NULL primary key AUTOINCREMENT, " +
		"customer_id integer not null , " +
		"status varchar(50) not null default ``, " +
		"total_qty integer not null, " +
		"total_amount integer not null, " +
		"created_at datetime not null default current_timestamp, " +
		"updated_at datetime not null default current_timestamp, " +
		"FOREIGN KEY (customer_id) REFERENCES customers(id)" +
		")").Error

	if err != nil {
		panic(err)
	}

	err = DB.Exec("CREATE TABLE IF NOT EXISTS order_units (" +
		"`id` integer NOT NULL primary key AUTOINCREMENT, " +
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
