<<<<<<< HEAD
# Enterprise Modular Monolith with Full Observability

A robust, enterprise-grade Go application implementing clean architecture, modular monolith patterns, and full observability.
=======
# Enterprise Go Observability & Modular Monolith API

A high-performance enterprise-grade Go application built with a **Modular Monolith / Clean Layered Architecture**, integrated with a **Full Observability Stack** (OpenTelemetry, Elasticsearch, Kibana, Prometheus, Grafana, and Jaeger).
>>>>>>> ad52c4f (feat: complete modular enterprise architecture and README)

## Tech Stack
- **Language**: Go (Golang)
- **Framework**: Gin Gonic
- **Database**: PostgreSQL (External)
- **Observability**: OpenTelemetry (Traces), Prometheus (Metrics), Elasticsearch/Kibana (Logs), Jaeger (Tracing)
- **Infrastructure**: Docker, Docker Compose

<<<<<<< HEAD
## Modular Architecture (`/internal/modules`)
The application is organized into domain-specific modules:
- `tenant`: Tenant management
- `user`: User management
- `order`: Order processing
- `product`: Product catalog
- `payment`: Payment gateway integration
- `shipment`: Logistics and tracking
- `inventory`: Inventory stock management
- `warehouse`: Warehouse operations
- `audit`: System-wide audit logging

## Getting Started
1. **Database**: Ensure PostgreSQL is running externally on port 5432.
2. **Configuration**: Update `internal/infrastructure/db.go` with your database credentials.
3. **Build & Run**:
   ```bash
   docker compose up --build
   ```

## API Endpoints (v1)
All modules expose REST endpoints grouped under `/api/v1/`:
- `/tenant`, `/users`, `/orders`, `/products`, `/payments`, `/shipments`, `/inventory`, `/warehouses`, `/audit`
=======
## 🚀 Architecture & Modules

The application is structured cleanly into independent domains under `internal/modules/`:
- **Tenant**: Multi-tenant isolation and configuration.
- **User**: User lifecycle and registration management.
- **Order**: Order placement and transactional processing.
- **Product**: Product catalog and inventory metadata.
- **Payment**: Payment gateway integration and transaction logging.
- **Shipment**: Courier tracking and logistics status.
- **Inventory & Warehouse**: Stock level tracking across warehouses.
- **Audit Log**: System-wide event tracking and security logging.

---

## 📡 Available API Endpoints

All endpoints are prefixed with `/api/v1` and support standard JSON request/response payloads.

### 1. Tenant Management
- **GET `/api/v1/tenants`** — Retrieve all tenants.
- **POST `/api/v1/tenants`** — Create a new tenant.

### 2. User Management
- **GET `/api/v1/users`** — List all registered users.
- **POST `/api/v1/users`** — Register a new user.
  - *Payload example:*
    ```json
    {
      "tenant_id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    }
    ```

### 3. Order Management
- **GET `/api/v1/orders`** — List all orders.
- **POST `/api/v1/orders`** — Place a new order.
  - *Payload example:*
    ```json
    {
      "user_id": 1,
      "total_amount": 150.00,
      "shipping_address": "123 Enterprise St."
    }
    ```

### 4. Product Catalog
- **GET `/api/v1/products`** — Retrieve product list.
- **POST `/api/v1/products`** — Add a new product.

### 5. Payment Processing
- **GET `/api/v1/payments`** — List all payment transactions.
- **POST `/api/v1/payments`** — Process a payment.
  - *Payload example:*
    ```json
    {
      "tenant_id": 1,
      "order_id": 1,
      "amount": 150.00,
      "payment_method": "credit_card"
    }
    ```

### 6. Shipment & Logistics
- **GET `/api/v1/shipments`** — List shipments.
- **POST `/api/v1/shipments`** — Create a shipment / tracking record.

### 7. Inventory & Warehouse
- **GET `/api/v1/inventory`** — Check inventory stock levels.
- **GET `/api/v1/warehouses`** — List active warehouses.

### 8. Audit Logging
- **GET `/api/v1/audit`** — Retrieve system audit trails.
- **POST `/api/v1/audit`** — Record an audit event.
  - *Payload example:*
    ```json
    {
      "tenant_id": 1,
      "user_id": 1,
      "action": "DELETE_USER",
      "details": "Admin removed inactive user account."
    }
    ```

---

## 🔍 Observability Stack Ports

- **Application API**: `http://localhost:8080`
- **Jaeger UI (Traces)**: `http://localhost:16686`
- **Prometheus (Metrics)**: `http://localhost:9090`
- **Grafana (Dashboards)**: `http://localhost:3000`
- **Kibana (Logs)**: `http://localhost:5601`
- **Elasticsearch**: `http://localhost:9200`
>>>>>>> ad52c4f (feat: complete modular enterprise architecture and README)
