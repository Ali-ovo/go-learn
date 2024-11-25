package controller

import (
	inventory_pb "shop/api/inventory/v1"
	"shop/app/shop_srv/inventory/srv/internal/service"
)

type InventoryServer struct {
	inventory_pb.UnimplementedInventoryServer
	srv service.ServiceFactory
}

func NewInventoryServer(srv service.ServiceFactory) inventory_pb.InventoryServer {
	return &InventoryServer{srv: srv}
}
