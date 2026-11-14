# Enterprise Modular Monolith with Full Observability

This is a production-grade enterprise application built with a Modular Monolith architecture in Go, featuring a full observability stack (OpenTelemetry, Elasticsearch, Kibana, Prometheus, Grafana, and Jaeger) backed by PostgreSQL.

---

## Architecture & Modules

The application is organized into isolated domain modules under `internal/modules/`:
- **Tenant**: Multi-tenancy support and organization management.
- **User**: User accounts, roles, and authentication context.
- **Product**: Product catalog and SKU tracking.
- **Warehouse**: Inventory location and warehouse management.
- **Inventory**: Stock levels per warehouse and product.
- **Order**: Customer order placement and processing.
- **Payment**: Transaction processing and payment gateways.
- **Shipment**: Courier tracking and delivery logistics.

---

## API Endpoints & CRUD Reference

### 1. Tenants (`/api/v1/tenants`)
- `GET /api/v1/tenants` - List all tenants
- `POST /api/v1/tenants` - Create a new tenant
- `GET /api/v1/tenants/:id` - Get tenant by ID
- `PUT /api/v1/tenants/:id` - Update tenant
- `DELETE /api/v1/tenants/:id` - Delete tenant

### 2. Users (`/api/v1/users`)
- `GET /api/v1/users` - List all users
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### 3. Products (`/api/v1/products`)
- `GET /api/v1/products` - List all products
- `POST /api/v1/products` - Create a new product
- `GET /api/v1/products/:id` - Get product by ID
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product

### 4. Warehouses (`/api/v1/warehouses`)
- `GET /api/v1/warehouses` - List all warehouses
- `POST /api/v1/warehouses` - Create a new warehouse
- `GET /api/v1/warehouses/:id` - Get warehouse by ID
- `PUT /api/v1/warehouses/:id` - Update warehouse
- `DELETE /api/v1/warehouses/:id` - Delete warehouse

### 5. Inventory (`/api/v1/inventory`)
- `GET /api/v1/inventory` - List inventory levels
- `POST /api/v1/inventory` - Adjust stock levels
- `GET /api/v1/inventory/:id` - Get inventory item by ID

### 6. Orders (`/api/v1/orders`)
- `GET /api/v1/orders` - List all orders
- `POST /api/v1/orders` - Create a new order
- `GET /api/v1/orders/:id` - Get order by ID
- `PUT /api/v1/orders/:id` - Update order
- `DELETE /api/v1/orders/:id` - Delete order

### 7. Payments (`/api/v1/payments`)
- `GET /api/v1/payments` - List all payments
- `POST /api/v1/payments` - Process a payment
- `GET /api/v1/payments/:id` - Get payment by ID
- `PUT /api/v1/payments/:id` - Update payment
- `DELETE /api/v1/payments/:id` - Delete payment

### 8. Shipments (`/api/v1/shipments`)
- `GET /api/v1/shipments` - List all shipments
- `POST /api/v1/shipments` - Create a shipment
- `GET /api/v1/shipments/:id` - Get shipment by ID
- `PUT /api/v1/shipments/:id` - Update shipment
- `DELETE /api/v1/shipments/:id` - Delete shipment

---

## Ready-to-Use cURL Commands

### Tenants
```bash
# List Tenants
curl -X GET http://localhost:8080/api/v1/tenants

# Create Tenant
curl -X POST http://localhost:8080/api/v1/tenants \
  -H "Content-Type: application/json" \
  -d '{"name": "Enterprise Corp", "subdomain": "enterprise", "domain": "enterprise.com", "status": "active"}'

# Get Tenant by ID
curl -X GET http://localhost:8080/api/v1/tenants/1

# Update Tenant
curl -X PUT http://localhost:8080/api/v1/tenants/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Enterprise Corp Updated", "subdomain": "enterprise", "domain": "enterprise.com", "status": "active"}'

# Delete Tenant
curl -X DELETE http://localhost:8080/api/v1/tenants/1
```

### Users
```bash
# List Users
curl -X GET http://localhost:8080/api/v1/users

# Create User
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": 1, "name": "John Doe", "email": "john@enterprise.com", "role": "admin"}'

# Get User by ID
curl -X GET http://localhost:8080/api/v1/users/1

# Update User
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": 1, "name": "John Updated", "email": "john@enterprise.com", "role": "admin"}'

# Delete User
curl -X DELETE http://localhost:8080/api/v1/users/1
```

### Products
```bash
# List Products
curl -X GET http://localhost:8080/api/v1/products

# Create Product
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": 1, "name": "MacBook Pro", "description": "M3 Max 16-inch", "price": 3499.00, "stock": 50, "sku": "MBP-M3-16", "status": "active"}'

# Get Product by ID
curl -X GET http://localhost:8080/api/v1/products/1

# Update Product
curl -X PUT http://localhost:8080/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": 1, "name": "MacBook Pro 2026", "description": "Updated specs", "price": 3599.00, "stock": 45, "sku": "MBP-M3-16", "status": "active"}'

# Delete Product
curl -X DELETE http://localhost:8080/api/v1/products/1
```

### Warehouses
```bash
# List Warehouses
curl -X GET http://localhost:8080/api/v1/warehouses

# Create Warehouse
curl -X POST http://localhost:8080/api/v1/warehouses \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": 1, "name": "Warehouse Jakarta", "location": "Jakarta Selatan", "status": "active"}'

# Get Warehouse by ID
curl -X GET http://localhost:8080/api/v1/warehouses/1

# Update Warehouse
curl -X PUT http://localhost:8080/api/v1/warehouses/1 \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": 1, "name": "Warehouse Jakarta Pusat", "location": "Jakarta Pusat", "status": "active"}'

# Delete Warehouse
curl -X DELETE http://localhost:8080/api/v1/warehouses/1
```

### Inventory
```bash
# List Inventory
curl -X GET http://localhost:8080/api/v1/inventory

# Adjust/Upsert Stock
curl -X POST http://localhost:8080/api/v1/inventory \
  -H "Content-Type: application/json" \
  -d '{"product_id": 1, "warehouse_id": 1, "quantity": 100}'

# Get Inventory by ID
curl -X GET http://localhost:8080/api/v1/inventory/1
```

### Orders
```bash
# List Orders
curl -X GET http://localhost:8080/api/v1/orders

# Create Order
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": 1, "user_id": 1, "total_amount": 3499.00, "shipping_addr": "Jl. Sudirman No. 1", "status": "pending"}'

# Get Order by ID
curl -X GET http://localhost:8080/api/v1/orders/1

# Update Order
curl -X PUT http://localhost:8080/api/v1/orders/1 \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": 1, "user_id": 1, "total_amount": 3499.00, "shipping_addr": "Jl. Sudirman No. 1", "status": "paid"}'

# Delete Order
curl -X DELETE http://localhost:8080/api/v1/orders/1
```

### Payments
```bash
# List Payments
curl -X GET http://localhost:8080/api/v1/payments

# Create Payment
curl -X POST http://localhost:8080/api/v1/payments \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": 1, "order_id": 1, "amount": 3499.00, "payment_method": "credit_card", "status": "success"}'

# Get Payment by ID
curl -X GET http://localhost:8080/api/v1/payments/1

# Update Payment
curl -X PUT http://localhost:8080/api/v1/payments/1 \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": 1, "order_id": 1, "amount": 3499.00, "payment_method": "credit_card", "status": "settled"}'

# Delete Payment
curl -X DELETE http://localhost:8080/api/v1/payments/1
```

### Shipments
```bash
# List Shipments
curl -X GET http://localhost:8080/api/v1/shipments

# Create Shipment
curl -X POST http://localhost:8080/api/v1/shipments \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": 1, "order_id": 1, "courier": "JNE Express", "tracking_number": "JNE-123456789", "status": "dispatched"}'

# Get Shipment by ID
curl -X GET http://localhost:8080/api/v1/shipments/1

# Update Shipment
curl -X PUT http://localhost:8080/api/v1/shipments/1 \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": 1, "order_id": 1, "courier": "JNE Express", "tracking_number": "JNE-123456789", "status": "delivered"}'

# Delete Shipment
curl -X DELETE http://localhost:8080/api/v1/shipments/1
```
