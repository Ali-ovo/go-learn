package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	gorm.Model
	CreditCards []CreditCard `gorm:"foreignKey:UserRefer"`
}

type CreditCard struct {
	gorm.Model
	Name      string
	UserRefer uint
}

func main() {
	dsn := "root:123456@tcp(192.168.189.128:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger

			Colorful: true, // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	// db.AutoMigrate(&CreditCard{})

	// user := User{}
	// db.Create(&user)

	// db.Create(&CreditCard{
	// 	Name:      "12",
	// 	UserRefer: user.ID,
	// })

	// db.Create(&CreditCard{
	// 	Name:      "34",
	// 	UserRefer: user.ID,
	// })

	var user User
	db.Preload("CreditCards").First(&user)

	for _, cart := range user.CreditCards {
		fmt.Println(cart)
	}

}
