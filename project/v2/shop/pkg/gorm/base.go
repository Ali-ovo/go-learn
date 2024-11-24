package gorm

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int64          `gorm:"primarykey;type:int comment 'ID'" json:"id"`
	CreatedAt time.Time      `gorm:"column:add_time;type:datetime comment '创建时间'" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:update_time;type:datetime comment '更新时间'" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime comment '删除时间'" json:"deleted_at"`
	IsDeleted bool
}
