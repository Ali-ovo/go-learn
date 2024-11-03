package mysql

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	umock "go-learn/ch19/mock"
)

func TestGetUserByMobile(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("init sqlmock failed: %v", err)
	}

	defer func(db *sql.DB) {
		_ = db.Close()
		if err != nil {
			t.Fatalf("close db failed: %v", err)
		}

	}(db)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("init gorm failed: %v", err)
	}
	mobile := "18"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`mobile` = ? ORDER BY `users`.`mobile` LIMIT 1")).
		WithArgs(mobile).WillReturnRows(
		sqlmock.NewRows([]string{"mobile", "password", "nick_name"}).AddRow("183757139090",
			"123456", "ali_18"))

	mock.ExpectClose()

	userData := NewUser(gormDB)
	user, err := userData.GetUserByMobile(context.Background(), mobile)
	assert.Nil(t, err)

	expUser := umock.User{
		Mobile:   "183757139090",
		Password: "123456",
		NickName: "ali_18",
	}
	assert.Equal(t, expUser, user)
}
