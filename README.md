# gRPC Order Management Service

A microservice-based order management system built with Go, featuring gRPC for inter-service communication and HTTP for external APIs.

## Architecture

```
┌──────────────────────────────────────────┐       gRPC (:9000)       ┌─────────────────────┐
│            Orders Service                │◄────────────────────────►│   Kitchen Service   │
│                                          │                          │                     │
│  HTTP API (:8000)    gRPC Server (:9000) │                          │  HTTP UI (:3000)    │
└──────────────────────────────────────────┘                          └─────────────────────┘
```

**Orders Service** — The core service that manages order data. Exposes both an HTTP REST API and a gRPC server.

**Kitchen Service** — A frontend that connects to the Orders Service via gRPC and renders orders in an HTML table.

## Project Structure

```
.
├── protobuf/                        # Protobuf definitions
│   └── orders.proto
├── services/
│   ├── common/
│   │   ├── genproto/orders/         # Generated Go code from proto
│   │   └── utils/                   # Shared HTTP utilities
│   ├── kitchen/                     # Kitchen service (gRPC client + HTTP UI)
│   │   ├── main.go
│   │   └── http.go
│   └── orders/                      # Orders service (HTTP + gRPC server)
│       ├── main.go
│       ├── http.go
│       ├── grpc.go
│       ├── handlers/orders/         # HTTP and gRPC request handlers
│       ├── service/                 # Business logic layer
│       └── types/                   # Interface definitions
├── Makefile
├── go.mod
└── go.sum
```

## Prerequisites

- **Go** 1.25+
- **protoc** (Protocol Buffers compiler) — only needed if regenerating proto files
- **protoc-gen-go** and **protoc-gen-go-grpc** — Go protobuf plugins

## Getting Started

### 1. Run the Orders Service

```sh
make run-orders
```

This starts both the HTTP server on `:8000` and the gRPC server on `:9000`.

### 2. Run the Kitchen Service

In a separate terminal:

```sh
make run-kitchen
```

This starts the Kitchen UI on `:3000`. Open [http://localhost:3000](http://localhost:3000) in a browser to see the orders table.

## API Reference

### HTTP (Orders Service — port 8000)

#### Create Order

```
POST /orders
Content-Type: application/json

{
  "customerID": 1,
  "productID": 100,
  "quantity": 3
}
```

Response:

```json
{"status": "success"}
```

### gRPC (Orders Service — port 9000)

Defined in `protobuf/orders.proto`:

| RPC | Request | Response |
|-----|---------|----------|
| `CreateOrder` | `CreateOrderRequest` | `CreateOrderResponse` |
| `GetOrders` | `GetOrdersRequest` | `GetOrderResponse` |

`GetOrders` accepts a `customerID` field. Pass `0` to retrieve all orders.

## Regenerating Protobuf Code

If you modify `protobuf/orders.proto`, regenerate the Go code:

```sh
make gen
```

## Running Tests

```sh
go test ./...
```

With race detection:

```sh
go test -race ./...
```

## Tech Stack

- **Go** — primary language
- **gRPC** / **Protocol Buffers** — inter-service communication
- **net/http** — HTTP REST API and Kitchen UI
