package database

import (
	entities "github.com/Dominikk7/PhotoBomb/backend/pkg/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DB_USERNAME = "root"
const DB_NAME_I = "images"
const DB_NAME_U = "users"
const DB_HOST = "172.17.0.2" //172.17.0.2 127.0.0.1
const DB_PORT = "3306"

var ImageInstance *gorm.DB
var UserInstance *gorm.DB
var err error

func Connect(DB_PASSWORD string) {
	//Image DB
	dsn := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME_I + "?" + "parseTime=true&loc=Local"
	ImageInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Image Database...")

	//User DB
	dsnU := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME_U + "?" + "parseTime=true&loc=Local"
	UserInstance, err = gorm.Open(mysql.Open(dsnU), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to User Database...")
}

func Migrate() {
	ImageInstance.AutoMigrate(&entities.Image{})
	UserInstance.AutoMigrate(&entities.User{})
	log.Println("Database Migration Completed...")
}
