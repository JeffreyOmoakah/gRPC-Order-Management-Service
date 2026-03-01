package service

import (
	"context"
	"sync"
	"testing"

	"github.com/JeffreyOmoakah/gRPC-Order-Management-Service/services/common/genproto/orders"
)

func TestNextID_AutoIncrements(t *testing.T) {
	s := NewOrderService()

	id1 := s.NextID()
	id2 := s.NextID()
	id3 := s.NextID()

	if id1 != 1 || id2 != 2 || id3 != 3 {
		t.Errorf("expected IDs 1,2,3 got %d,%d,%d", id1, id2, id3)
	}
}

func TestCreateOrder(t *testing.T) {
	s := NewOrderService()
	ctx := context.Background()

	order := &orders.Order{
		OrderID:    s.NextID(),
		CustomerID: 10,
		ProductID:  20,
		Quantity:   5,
	}

	err := s.CreateOrder(ctx, order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := s.GetOrders(ctx, 0)
	if len(result) != 1 {
		t.Fatalf("expected 1 order, got %d", len(result))
	}
	if result[0].OrderID != 1 {
		t.Errorf("expected OrderID 1, got %d", result[0].OrderID)
	}
	if result[0].CustomerID != 10 {
		t.Errorf("expected CustomerID 10, got %d", result[0].CustomerID)
	}
}

func TestGetOrders_FiltersByCustomerID(t *testing.T) {
	s := NewOrderService()
	ctx := context.Background()

	s.CreateOrder(ctx, &orders.Order{OrderID: s.NextID(), CustomerID: 1, ProductID: 10, Quantity: 1})
	s.CreateOrder(ctx, &orders.Order{OrderID: s.NextID(), CustomerID: 2, ProductID: 20, Quantity: 2})
	s.CreateOrder(ctx, &orders.Order{OrderID: s.NextID(), CustomerID: 1, ProductID: 30, Quantity: 3})

	// Filter for customer 1
	result := s.GetOrders(ctx, 1)
	if len(result) != 2 {
		t.Fatalf("expected 2 orders for customer 1, got %d", len(result))
	}

	// Filter for customer 2
	result = s.GetOrders(ctx, 2)
	if len(result) != 1 {
		t.Fatalf("expected 1 order for customer 2, got %d", len(result))
	}

	// No filter (customerID=0) returns all
	result = s.GetOrders(ctx, 0)
	if len(result) != 3 {
		t.Fatalf("expected 3 orders for no filter, got %d", len(result))
	}

	// Non-existent customer
	result = s.GetOrders(ctx, 999)
	if len(result) != 0 {
		t.Fatalf("expected 0 orders for customer 999, got %d", len(result))
	}
}

func TestConcurrentAccess(t *testing.T) {
	s := NewOrderService()
	ctx := context.Background()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			order := &orders.Order{
				OrderID:    s.NextID(),
				CustomerID: 1,
				ProductID:  1,
				Quantity:   1,
			}
			s.CreateOrder(ctx, order)
		}()
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.GetOrders(ctx, 0)
		}()
	}

	wg.Wait()

	result := s.GetOrders(ctx, 0)
	if len(result) != 100 {
		t.Errorf("expected 100 orders after concurrent writes, got %d", len(result))
	}
}
