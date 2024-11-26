package do

import (
	"shop/pkg/gorm"
	"time"
)

/*
1. 密文
 1. 对称加密
 2. 非对称加密
 3. md5 信息摘要算法
2. 密文不可反解
 密码如果不可反解, 用户找回密码
*/

type UserDO struct {
	gorm.BaseModel
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11) comment '手机号';not null" json:"mobile"`
	Password string     `gorm:"type:varchar(100) comment '密码';not null" json:"password"`
	NickName string     `gorm:"type:varchar(20) comment '用户名'" json:"nick_name"`
	Birthday *time.Time `gorm:"type:datetime comment '生日'" json:"birthday"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment 'female 表示女, male 表示男'" json:"gender"`
	Role     int        `gorm:"column:role;default:1; type:int comment '1 表示普通用户, 2 表示管理员'" json:"role"`
}

func (u *UserDO) TableName() string {
	return "user"
}

type UserDOList struct {
	TotalCount int64     `json:"totalCount,omitempty"` // 总数
	Items      []*UserDO `json:"data"`                 // 数据
}
