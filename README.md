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

## 📡 Full CRUD API Endpoints

Every enterprise domain now supports complete **CRUD** operations (`GET`, `POST`, `PUT`, `DELETE`, and ID lookups):

### 1. Tenants (`/api/v1/tenants`)
- `GET /api/v1/tenants` — List all tenants
- `POST /api/v1/tenants` — Create a new tenant
- `GET /api/v1/tenants/:id` — Get tenant by ID
- `PUT /api/v1/tenants/:id` — Update tenant details
- `DELETE /api/v1/tenants/:id` — Delete a tenant

### 2. Users (`/api/v1/users`)
- `GET /api/v1/users` — List all users
- `POST /api/v1/users` — Register a new user
- `GET /api/v1/users/:id` — Get user by ID
- `PUT /api/v1/users/:id` — Update user details
- `DELETE /api/v1/users/:id` — Remove user

### 3. Orders (`/api/v1/orders`)
- `GET /api/v1/orders` — List orders
- `POST /api/v1/orders` — Place order
- `GET /api/v1/orders/:id` — Get order details
- `PUT /api/v1/orders/:id` — Update order status
- `DELETE /api/v1/orders/:id` — Cancel/Delete order

### 4. Products (`/api/v1/products`)
- `GET /api/v1/products` — List catalog products
- `POST /api/v1/products` — Add product
- `GET /api/v1/products/:id` — Get product details
- `PUT /api/v1/products/:id` — Update product info
- `DELETE /api/v1/products/:id` — Remove product

### 5. Payments (`/api/v1/payments`)
- `GET /api/v1/payments` — List payments
- `POST /api/v1/payments` — Process payment
- `GET /api/v1/payments/:id` — Get payment receipt/status
- `PUT /api/v1/payments/:id` — Update payment status
- `DELETE /api/v1/payments/:id` — Void/Delete payment

### 6. Shipments (`/api/v1/shipments`)
- `GET /api/v1/shipments` — List shipments
- `POST /api/v1/shipments` — Create shipment
- `GET /api/v1/shipments/:id` — Track shipment by ID
- `PUT /api/v1/shipments/:id` — Update tracking status
- `DELETE /api/v1/shipments/:id` — Cancel shipment

### 7. Inventory & Warehouses (`/api/v1/inventory`, `/api/v1/warehouses`)
- `GET /api/v1/inventory` — Check stock levels
- `POST /api/v1/inventory` — Add/Initialize stock
- `PUT /api/v1/inventory/:id` — Adjust stock quantity
- `GET /api/v1/warehouses` — List warehouses
- `POST /api/v1/warehouses` — Add warehouse
- `GET /api/v1/warehouses/:id` — Get warehouse details
- `PUT /api/v1/warehouses/:id` — Update warehouse info
- `DELETE /api/v1/warehouses/:id` — Remove warehouse

### 8. Audit Logs (`/api/v1/audit`)
- `GET /api/v1/audit` — List audit logs
- `POST /api/v1/audit` — Record audit event

---

## 🔍 Observability Stack Ports

- **Application API**: `http://localhost:8080`
- **Jaeger UI (Traces)**: `http://localhost:16686`
- **Prometheus (Metrics)**: `http://localhost:9090`
- **Grafana (Dashboards)**: `http://localhost:3000`
- **Kibana (Logs)**: `http://localhost:5601`
- **Elasticsearch**: `http://localhost:9200`
