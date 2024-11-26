package data

import "gorm.io/gorm"

type DataFactory interface {
	User() UserStore

	Begin() *gorm.DB
}
