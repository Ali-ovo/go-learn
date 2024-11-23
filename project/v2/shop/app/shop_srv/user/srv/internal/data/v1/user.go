package data

import (
	"context"
	metav1 "shop/gmicro/pkg/common/meta/v1"
	"time"

	"shop/pkg/gorm"
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

type UserStore interface {
	/*
		有数据访问的方法, 一定要有 error
		参数中最好有 ctx 后期便于管理 比如 cancel掉 链路追踪等
	*/
	// List 用户列表
	List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*UserDOList, error)
	// GetByMobile 通过手机号码查询用户
	GetByMobile(ctx context.Context, mobile string) (*UserDO, error)
	// GetByID 通过 id 查询用户
	GetByID(ctx context.Context, id uint64) (*UserDO, error)
	// Create 添加用户
	Create(ctx context.Context, user *UserDO) error
	// Update 更新用户
	Update(ctx context.Context, user *UserDO) error
}
