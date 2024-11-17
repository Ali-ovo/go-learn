package model

import "time"
import "gorm.io/gorm"

type BaseModel struct {
	ID        int32          `gorm:"primarykey;type:int comment 'ID'" json:"id"`
	CreatedAt time.Time      `gorm:"column:add_time;type:datetime comment '创建时间'" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:update_time;type:datetime comment '更新时间'" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime comment '删除时间'" json:"deleted_at"`
	IsDeleted bool
}

/*
1. 密文
 1. 对称加密
 2. 非对称加密
 3. md5 信息摘要算法
2. 密文不可反解
 密码如果不可反解, 用户找回密码
*/

type User struct {
	BaseModel
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11) comment '手机号';not null" json:"mobile"`
	Password string     `gorm:"type:varchar(100) comment '密码';not null" json:"password"`
	NickName string     `gorm:"type:varchar(20) comment '用户名'" json:"nick_name"`
	Birthday *time.Time `gorm:"type:datetime comment '生日'" json:"birthday"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment 'female 表示女, male 表示男'" json:"gender"`
	Role     int        `gorm:"column:role;default:1; type:int comment '1 表示普通用户, 2 表示管理员'" json:"role"`
}
