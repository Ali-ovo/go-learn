package do

import (
	"shop/pkg/gorm"
	"time"
)

type OrderGoods struct {
	gorm.BaseModel
	Order int32 `gorm:"type:int;index"`
	Goods int32 `gorm:"type:int;index"`

	// 把商品的信息保存下来了	字段冗余, 实际高并发系统中我们一般都不会遵循三范式	做镜像 比如那个时候 商品的价格是多少
	GoodsName  string  `gorm:"type:varchar(100);index"`
	GoodsImage string  `gorm:"type:varchar(200)"`
	GoodsPrice float32 `gorm:"type:float comment '商品价格'"`
	Nums       int32   `gorm:"type:int"`
}

func (OrderGoods) TableName() string {
	return "ordergoods"
}

type OrderInfoDO struct {
	gorm.BaseModel
	OrderGoods   []*OrderGoods `gorm:"foreignKey:Order;references:ID" json:"goods"`
	User         int32         `gorm:"type:int;index"`
	OrderSn      string        `gorm:"type:varchar(30) comment '订单号';index"` // 自己生成的订单号
	PayType      string        `gorm:"type:varchar(20) comment 'alipay(支付宝)， wechat(微信)'"`
	Status       string        `gorm:"type:varchar(20)  comment 'PAYING(待支付), TRADE_SUCCESS(成功)， TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)'"` //status大家可以考虑使用iota来做
	TradeNo      string        `gorm:"type:varchar(100) comment '交易号'"`                                                                                            // 交易号就是支付宝的订单号 查账
	OrderMount   float32       `gorm:"type:float comment '订单金额'"`
	PayTime      *time.Time    `gorm:"type:datetime comment '订单时间'"`
	Address      string        `gorm:"type:varchar(100) comment '收件人地址'"`
	SignerName   string        `gorm:"type:varchar(20) comment '收件人名字'"`
	SingerMobile string        `gorm:"type:varchar(11) comment '收件人手机号'"`
	Post         string        `gorm:"type:varchar(20) comment '留言信息'"`
}

func (OrderInfoDO) TableName() string {
	return "orderinfo"
}

type OrderInfoDOList struct {
	TotalCount int64          `json:"totalCount,omitempty"`
	Items      []*OrderInfoDO `json:"items"`
}
