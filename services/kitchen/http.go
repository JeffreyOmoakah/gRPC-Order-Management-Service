package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/JeffreyOmoakah/gRPC-Order-Management-Service/services/common/genproto/orders"
)

type httpServer struct {
	addr string
}

func NewHttpServer(addr string) *httpServer {
	return &httpServer{addr: addr}
}

func (s *httpServer) Run() error {
	router := http.NewServeMux()

	conn := NewGRPCClient(":9000")
	defer conn.Close()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := orders.NewOrderServiceClient(conn)

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
		defer cancel()

		_, err := c.CreateOrder(ctx, &orders.CreateOrderRequest{
			CustomerID: 24,
			ProductID:  3123,
			Quantity:   2,
		})
		if err != nil {
			log.Printf("client error: %v", err)
			http.Error(w, "failed to create order", http.StatusInternalServerError)
			return
		}

		res, err := c.GetOrders(ctx, &orders.GetOrdersRequest{
			CustomerID: 42,
		})
		if err != nil {
			log.Printf("client error: %v", err)
			http.Error(w, "failed to get orders", http.StatusInternalServerError)
			return
		}

		t := template.Must(template.New("orders").Parse(ordersTemplate))

		if err := t.Execute(w, res.GetOrders()); err != nil {
			log.Printf("template error: %v", err)
		}
	})

	log.Println("Starting server on", s.addr)
	return http.ListenAndServe(s.addr, router)
}

var ordersTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Kitchen Orders</title>
</head>
<body>
    <h1>Orders List</h1>
    <table border="1">
        <tr>
            <th>Order ID</th>
            <th>Customer ID</th>
            <th>Quantity</th>
        </tr>
        {{range .}}
        <tr>
            <td>{{.OrderID}}</td>
            <td>{{.CustomerID}}</td>
            <td>{{.Quantity}}</td>
        </tr>
        {{end}}
    </table>
</body>
</html>`
