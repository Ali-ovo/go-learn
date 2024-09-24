package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Language struct {
	// gorm.Model
	Name    string
	AddTime time.Time
}

// func (Language) TableName() string {
// 	return "my_languages"
// }

func (l *Language) BeforeCreate(tx *gorm.DB) (err error) {
	l.AddTime = time.Now()
	return
}

func main() {
	dsn := "root:123456@tcp(192.168.189.128:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{TablePrefix: "ali_"},
		Logger:         newLogger,
	})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Language{})

	db.Create(&Language{
		Name: "python",
	})

}
