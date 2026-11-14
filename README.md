# Enterprise Modular Monolith with Full Observability

A robust, enterprise-grade Go application implementing clean architecture, modular monolith patterns, and full observability.

## Tech Stack
- **Language**: Go (Golang)
- **Framework**: Gin Gonic
- **Database**: PostgreSQL (External)
- **Observability**: OpenTelemetry (Traces), Prometheus (Metrics), Elasticsearch/Kibana (Logs), Jaeger (Tracing)
- **Infrastructure**: Docker, Docker Compose

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
