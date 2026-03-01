package types

import (
	"context"

	"github.com/JeffreyOmoakah/gRPC-Order-Management-Service/services/common/genproto/orders"
)

type OrderService interface {
	CreateOrder(context.Context, *orders.Order) error
	GetOrders(context.Context, int32) []*orders.Order
	NextID() int32
}
