package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID           uint           // Standard field for the primary key
	Name         string         // 一个常规字符串字段
	Email        *string        // 一个指向字符串的指针, allowing for null values
	Age          uint8          // 一个未签名的8位整数
	Birthday     *time.Time     // A pointer to time.Time, can be null
	MemberNumber sql.NullString // Uses sql.NullString to handle nullable strings
	ActivatedAt  sql.NullTime   // Uses sql.NullTime for nullable time fields
	CreatedAt    time.Time      // 创建时间（由GORM自动管理）
	UpdatedAt    time.Time      // 最后一次更新时间（由GORM自动管理）
}

func main() {
	dsn := "root:123456@tcp(172.16.89.130:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

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

	_ = db.AutoMigrate(&User{})

	user := User{
		Name: "ali",
	}
	// db.Create(&User{
	// 	Name: "ali233",
	// })
	result := db.Create(&user) // 插入后会自动生成 id
	fmt.Println(user.ID)
	fmt.Println(result.RowsAffected)

	// 可以成功更新零值
	// db.Model(&User{ID: 1}).Update("Name", "")

	// Updates 不会更新空值
	// empty := "" // 奇淫技巧解决更新为零值
	// db.Model(&User{ID: 1}).Updates(
	// 	User{
	// 		Email: &empty,
	// 	},
	// )

}
