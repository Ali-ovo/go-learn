package service

import (
	"context"
	"shop/app/shop_srv/order/srv/internal/domain/dto"
	metav1 "shop/gmicro/pkg/common/meta/v1"
)

type OrderSrv interface {
	Get(ctx context.Context, orderSn string) (*dto.OrderDTO, error)                                            // 获取订单信息
	List(ctx context.Context, userID int64, meta metav1.ListMeta, orderby []string) (*dto.OrderDTOList, error) // 获取订单列表
	Submit(ctx context.Context, order *dto.OrderDTO) error                                                     // 创建订单流程(扣减库存->创建订单)
	Create(ctx context.Context, order *dto.OrderDTO) error                                                     // 创建订单
	CreateCom(ctx context.Context, order *dto.OrderDTO) error                                                  // 创建订单 补偿操作
	RollBack(ctx context.Context, order *dto.OrderDTO) error                                                   // 回滚订单
	RollBackCom(ctx context.Context, order *dto.OrderDTO) error                                                // 回滚订单 补偿操作
	Update(ctx context.Context, order *dto.OrderDTO) error                                                     // 更新订单状态
}
