# Enterprise Go Observability & Modular Monolith API

A high-performance enterprise-grade Go application built with a **Modular Monolith / Clean Layered Architecture**, integrated with a **Full Observability Stack** (OpenTelemetry, Elasticsearch, Kibana, Prometheus, Grafana, and Jaeger).

---

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

## 🔒 Security & Middleware

All `/api/v1/*` endpoints are protected by an enterprise-grade **Bearer Token Authorization Middleware** (`internal/middleware/auth.go`).
- **Header required**: `Authorization: Bearer <token>`
- **Context Injection**: Injects `tenant_id`, `user_id`, and `role` into request contexts.

---

## 📡 Comprehensive API Endpoints List

### 1. Tenant Management
- **GET `/api/v1/tenants`** — Retrieve all tenants.
- **POST `/api/v1/tenants`** — Create a new tenant.

### 2. User Management
- **GET `/api/v1/users`** — List all registered users.
- **POST `/api/v1/users`** — Register a new user.

### 3. Order Management
- **GET `/api/v1/orders`** — List all orders.
- **POST `/api/v1/orders`** — Place a new order.

### 4. Product Catalog
- **GET `/api/v1/products`** — Retrieve product list.
- **POST `/api/v1/products`** — Add a new product.

### 5. Payment Processing
- **GET `/api/v1/payments`** — List all payment transactions.
- **POST `/api/v1/payments`** — Process a payment.

### 6. Shipment & Logistics
- **GET `/api/v1/shipments`** — List shipments.
- **POST `/api/v1/shipments`** — Create a shipment / tracking record.

### 7. Inventory & Warehouse
- **GET `/api/v1/inventory`** — Check inventory stock levels.
- **GET `/api/v1/warehouses`** — List active warehouses.

### 8. Audit Logging
- **GET `/api/v1/audit`** — Retrieve system audit trails.
- **POST `/api/v1/audit`** — Record an audit event.

---

## 🔍 Observability Stack Ports

- **Application API**: `http://localhost:8080`
- **Jaeger UI (Traces)**: `http://localhost:16686`
- **Prometheus (Metrics)**: `http://localhost:9090`
- **Grafana (Dashboards)**: `http://localhost:3000`
- **Kibana (Logs)**: `http://localhost:5601`
- **Elasticsearch**: `http://localhost:9200`
