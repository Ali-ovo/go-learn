package data

import (
	"context"
	"shop/gmicro/pkg/common/time"
)

type UserDO struct {
	ID       int64     `json:"id"`
	NickName string    `json:"nick_name"`
	Birthday time.Time `json:"birthday"`
	Gender   string    `json:"gender"`
	Role     int32     `json:"role"`
	Mobile   string    `json:"mobile"`
	PassWord string    `json:"password"`
}

type UserDOList struct {
	TotalCount int64     `json:"total_count,omitempty"` // 总数
	Items      []*UserDO `json:"items"`                 // 用户数据
	//Items []*upbv1.UserInfoResponse `json:"items"`
}

type UserData interface {
	Create(ctx context.Context, user *UserDO) error
	Update(ctx context.Context, user *UserDO) error
	Get(ctx context.Context, userID int64) (*UserDO, error)
	GetByMobile(ctx context.Context, mobile string) (*UserDO, error)
	CheckPassWord(ctx context.Context, password string, encryptedPwd string) error
}
