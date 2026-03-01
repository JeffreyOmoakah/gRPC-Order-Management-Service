package service

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/JeffreyOmoakah/gRPC-Order-Management-Service/services/common/genproto/orders"
)

type OrderService struct {
	mu       sync.RWMutex
	orders   []*orders.Order
	nextID   atomic.Int32
}

func NewOrderService() *OrderService {
	s := &OrderService{
		orders: make([]*orders.Order, 0),
	}
	s.nextID.Store(1)
	return s
}

func (s *OrderService) NextID() int32 {
	return s.nextID.Add(1) - 1
}

func (s *OrderService) CreateOrder(ctx context.Context, order *orders.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders = append(s.orders, order)
	return nil
}

func (s *OrderService) GetOrders(ctx context.Context, customerID int32) []*orders.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if customerID == 0 {
		result := make([]*orders.Order, len(s.orders))
		copy(result, s.orders)
		return result
	}

	var filtered []*orders.Order
	for _, o := range s.orders {
		if o.CustomerID == customerID {
			filtered = append(filtered, o)
		}
	}
	return filtered
}
