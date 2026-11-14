package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	trace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var tracer trace.Tracer

func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("otel-collector:4317"),
		otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String("massive-enterprise-app"),
		)),
	)
	otel.SetTracerProvider(tp)
	tracer = otel.Tracer("enterprise-tracer")
	return tp, nil
}

// PushAdvancedLog sends structured per-layer logs (Controller, Repository, Service) to Elasticsearch for Kibana Discover
func PushAdvancedLog(layer, component, status, message string, payload map[string]interface{}) {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://elasticsearch:9200"
	}

	logEntry := map[string]interface{}{
		"@timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		"service":    "massive-enterprise-app",
		"layer":      layer,     // e.g., "controller", "repository", "system"
		"component":  component, // e.g., "GetOrdersHandler", "OrdersRepo.FindAll"
		"status":     status,    // "SUCCESS" or "FAILED"
		"message":    message,
		"details":    payload,   // Request params, response summary, SQL errors, duration
	}

	body, err := json.Marshal(logEntry)
	if err != nil {
		return
	}

	go func() {
		client := &http.Client{Timeout: 3 * time.Second}
		indexName := fmt.Sprintf("app-logs-%s", time.Now().Format("2006.01.02"))
		url := fmt.Sprintf("%s/%s/_doc", esURL, indexName)
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
		if err == nil {
			resp.Body.Close()
		}
	}()
}

func main() {
	tp, err := initTracer()
	if err != nil {
		log.Printf("Warning: Failed to init tracer: %v", err)
	} else {
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}()
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:***@host.docker.internal:5432/testing-otel?sslmode=disable"
	}

	var db *sql.DB
	for i := 1; i <= 10; i++ {
		db, err = sql.Open("postgres", dbURL)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Attempt %d: Waiting for PostgreSQL connection...", i)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migration for 10+ relational enterprise tables
	schemaQueries := []string{
		`CREATE TABLE IF NOT EXISTS tenants (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			subdomain VARCHAR(100) UNIQUE NOT NULL,
			status VARCHAR(50) DEFAULT 'ACTIVE',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			tenant_id INT REFERENCES tenants(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			status VARCHAR(50) DEFAULT 'ACTIVE',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			tenant_id INT REFERENCES tenants(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			tenant_id INT REFERENCES tenants(id) ON DELETE CASCADE,
			category_id INT REFERENCES categories(id) ON DELETE SET NULL,
			sku VARCHAR(100) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			price NUMERIC(12, 2) NOT NULL,
			cost_price NUMERIC(12, 2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS warehouses (
			id SERIAL PRIMARY KEY,
			tenant_id INT REFERENCES tenants(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			location VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS inventory_stocks (
			id SERIAL PRIMARY KEY,
			warehouse_id INT REFERENCES warehouses(id) ON DELETE CASCADE,
			product_id INT REFERENCES products(id) ON DELETE CASCADE,
			quantity INT NOT NULL DEFAULT 0,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			tenant_id INT REFERENCES tenants(id) ON DELETE CASCADE,
			user_id INT REFERENCES users(id) ON DELETE CASCADE,
			total_amount NUMERIC(12, 2) NOT NULL,
			status VARCHAR(50) DEFAULT 'PENDING',
			shipping_address TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id INT REFERENCES orders(id) ON DELETE CASCADE,
			product_id INT REFERENCES products(id) ON DELETE CASCADE,
			quantity INT NOT NULL,
			unit_price NUMERIC(12, 2) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS payments (
			id SERIAL PRIMARY KEY,
			order_id INT REFERENCES orders(id) ON DELETE CASCADE,
			amount NUMERIC(12, 2) NOT NULL,
			payment_gateway VARCHAR(100) DEFAULT 'MIDTRANS',
			status VARCHAR(50) DEFAULT 'SUCCESS',
			transaction_ref VARCHAR(255) UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS shipments (
			id SERIAL PRIMARY KEY,
			order_id INT REFERENCES orders(id) ON DELETE CASCADE,
			courier VARCHAR(100) NOT NULL,
			tracking_number VARCHAR(255) UNIQUE NOT NULL,
			status VARCHAR(50) DEFAULT 'PROCESSING',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id SERIAL PRIMARY KEY,
			action VARCHAR(255) NOT NULL,
			details TEXT,
			ip_address VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, q := range schemaQueries {
		if _, err := db.Exec(q); err != nil {
			PushAdvancedLog("system", "MigrationRepo", "FAILED", "Failed to execute schema migration", map[string]interface{}{"error": err.Error()})
			log.Fatalf("Failed to execute schema migration: %v", err)
		}
	}
	PushAdvancedLog("system", "MigrationRepo", "SUCCESS", "Database schema migration completed successfully", nil)

	// Automated Massive Seeder if empty
	var userCount int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if userCount == 0 {
		log.Println("Seeding massive enterprise data...")
		PushAdvancedLog("system", "Seeder", "SUCCESS", "Starting database seeding", nil)

		db.Exec(`INSERT INTO tenants (id, name, subdomain, status) VALUES (1, 'Enterprise Corp Alpha', 'alpha', 'ACTIVE')`)
		db.Exec(`INSERT INTO tenants (id, name, subdomain, status) VALUES (2, 'Enterprise Corp Beta', 'beta', 'ACTIVE')`)
		db.Exec(`INSERT INTO warehouses (tenant_id, name, location) VALUES (1, 'Central Jakarta Warehouse', 'Jakarta')`)
		db.Exec(`INSERT INTO warehouses (tenant_id, name, location) VALUES (2, 'Surabaya Distribution Hub', 'Surabaya')`)

		for i := 1; i <= 500; i++ {
			db.Exec(`INSERT INTO users (tenant_id, name, email, status) VALUES ($1, $2, $3, $4)`,
				(i%2)+1, fmt.Sprintf("Enterprise User %d", i), fmt.Sprintf("user%d@enterprise.com", i), "ACTIVE")
		}
		for i := 1; i <= 10; i++ {
			db.Exec(`INSERT INTO categories (tenant_id, name, slug) VALUES ($1, $2, $3)`,
				(i%2)+1, fmt.Sprintf("Category %d", i), fmt.Sprintf("cat-%d", i))
		}
		for i := 1; i <= 300; i++ {
			db.Exec(`INSERT INTO products (tenant_id, category_id, sku, name, price, cost_price) VALUES ($1, $2, $3, $4, $5, $6)`,
				(i%2)+1, (i%10)+1, fmt.Sprintf("SKU-ENT-%04d", i), fmt.Sprintf("Enterprise Product %d", i), float64(i)*150.5, float64(i)*100.0)
			db.Exec(`INSERT INTO inventory_stocks (warehouse_id, product_id, quantity) VALUES ($1, $2, $3)`,
				(i%2)+1, i, 1000)
		}
		for i := 1; i <= 1000; i++ {
			db.Exec(`INSERT INTO orders (tenant_id, user_id, total_amount, status, shipping_address) VALUES ($1, $2, $3, $4, $5)`,
				(i%2)+1, (i%500)+1, float64(i)*250.0, "COMPLETED", fmt.Sprintf("Jl. Sudirman No. %d, Jakarta", i))
			db.Exec(`INSERT INTO payments (order_id, amount, status, transaction_ref) VALUES ($1, $2, $3, $4)`,
				i, float64(i)*250.0, "SUCCESS", fmt.Sprintf("TRX-ENT-REF-%06d", i))
			db.Exec(`INSERT INTO shipments (order_id, courier, tracking_number, status) VALUES ($1, $2, $3, $4)`,
				i, "JNE-EXPRESS", fmt.Sprintf("JD-AWB-%06d", i), "DELIVERED")
		}
		for i := 1; i <= 100; i++ {
			db.Exec(`INSERT INTO audit_logs (action, details, ip_address) VALUES ($1, $2, $3)`,
				"SYSTEM_BOOT_SEED", fmt.Sprintf("Audit log entry batch %d", i), "192.168.1.100")
		}

		PushAdvancedLog("system", "Seeder", "SUCCESS", "Massive enterprise database seeding completed successfully", map[string]interface{}{"total_users": 500, "total_orders": 1000})
		log.Println("Massive enterprise database seeding completed successfully!")
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("massive-enterprise-app"))

	// Health check controller
	r.GET("/health", func(c *gin.Context) {
		start := time.Now()
		PushAdvancedLog("controller", "HealthHandler", "SUCCESS", "Health check accessed", map[string]interface{}{
			"client_ip": c.ClientIP(),
			"method":    c.Request.Method,
		})
		c.JSON(http.StatusOK, gin.H{
			"service":      "massive-enterprise-observability-app",
			"status":       "healthy",
			"duration_ms":  time.Since(start).Milliseconds(),
			"timestamp":    time.Now().UTC(),
		})
	})

	// API v1 Group
	api := r.Group("/api/v1")
	{
		// 1. Tenants
		api.GET("/tenants", func(c *gin.Context) {
			start := time.Now()
			PushAdvancedLog("controller", "GetTenantsHandler", "SUCCESS", "Request received for tenants", map[string]interface{}{"client_ip": c.ClientIP()})

			// Repository Layer execution
			repoStart := time.Now()
			rows, err := db.Query("SELECT id, name, subdomain, status, created_at FROM tenants")
			if err != nil {
				PushAdvancedLog("repository", "TenantsRepo.FindAll", "FAILED", "Database query failed", map[string]interface{}{"error": err.Error(), "duration_ms": time.Since(repoStart).Milliseconds()})
				c.JSON(500, gin.H{"success": false, "error": err.Error()})
				return
			}
			defer rows.Close()

			var list []map[string]interface{}
			for rows.Next() {
				var id int
				var name, subdomain, status string
				var createdAt time.Time
				rows.Scan(&id, &name, &subdomain, &status, &createdAt)
				list = append(list, map[string]interface{}{"id": id, "name": name, "subdomain": subdomain, "status": status, "created_at": createdAt})
			}
			PushAdvancedLog("repository", "TenantsRepo.FindAll", "SUCCESS", "Query executed successfully", map[string]interface{}{"rows_count": len(list), "duration_ms": time.Since(repoStart).Milliseconds()})

			PushAdvancedLog("controller", "GetTenantsHandler", "SUCCESS", "Request completed", map[string]interface{}{"total_duration_ms": time.Since(start).Milliseconds()})
			c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
		})

		// 2. Users
		api.GET("/users", func(c *gin.Context) {
			start := time.Now()
			PushAdvancedLog("controller", "GetUsersHandler", "SUCCESS", "Request received for users", nil)

			repoStart := time.Now()
			rows, err := db.Query("SELECT id, tenant_id, name, email, status, created_at FROM users LIMIT 50")
			if err != nil {
				PushAdvancedLog("repository", "UsersRepo.FindLimit", "FAILED", "Database query failed", map[string]interface{}{"error": err.Error()})
				c.JSON(500, gin.H{"success": false, "error": err.Error()})
				return
			}
			defer rows.Close()

			var list []map[string]interface{}
			for rows.Next() {
				var id, tenantID int
				var name, email, status string
				var createdAt time.Time
				rows.Scan(&id, &tenantID, &name, &email, &status, &createdAt)
				list = append(list, map[string]interface{}{"id": id, "tenant_id": tenantID, "name": name, "email": email, "status": status, "created_at": createdAt})
			}
			PushAdvancedLog("repository", "UsersRepo.FindLimit", "SUCCESS", "Query executed", map[string]interface{}{"rows_count": len(list), "duration_ms": time.Since(repoStart).Milliseconds()})

			PushAdvancedLog("controller", "GetUsersHandler", "SUCCESS", "Request completed", map[string]interface{}{"total_duration_ms": time.Since(start).Milliseconds()})
			c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
		})

		// 3. Products
		api.GET("/products", func(c *gin.Context) {
			start := time.Now()
			PushAdvancedLog("controller", "GetProductsHandler", "SUCCESS", "Request received for products", nil)

			repoStart := time.Now()
			rows, err := db.Query("SELECT id, sku, name, price, cost_price FROM products LIMIT 50")
			if err != nil {
				PushAdvancedLog("repository", "ProductsRepo.FindLimit", "FAILED", "Database query failed", map[string]interface{}{"error": err.Error()})
				c.JSON(500, gin.H{"success": false, "error": err.Error()})
				return
			}
			defer rows.Close()

			var list []map[string]interface{}
			for rows.Next() {
				var id int
				var sku, name string
				var price, costPrice float64
				rows.Scan(&id, &sku, &name, &price, &costPrice)
				list = append(list, map[string]interface{}{"id": id, "sku": sku, "name": name, "price": price, "cost_price": costPrice})
			}
			PushAdvancedLog("repository", "ProductsRepo.FindLimit", "SUCCESS", "Query executed", map[string]interface{}{"rows_count": len(list), "duration_ms": time.Since(repoStart).Milliseconds()})

			PushAdvancedLog("controller", "GetProductsHandler", "SUCCESS", "Request completed", map[string]interface{}{"total_duration_ms": time.Since(start).Milliseconds()})
			c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
		})

		// 4. Orders
		api.GET("/orders", func(c *gin.Context) {
			start := time.Now()
			PushAdvancedLog("controller", "GetOrdersHandler", "SUCCESS", "Request received for orders", nil)

			repoStart := time.Now()
			rows, err := db.Query("SELECT id, tenant_id, user_id, total_amount, status, shipping_address, created_at FROM orders LIMIT 50")
			if err != nil {
				PushAdvancedLog("repository", "OrdersRepo.FindLimit", "FAILED", "Database query failed", map[string]interface{}{"error": err.Error()})
				c.JSON(500, gin.H{"success": false, "error": err.Error()})
				return
			}
			defer rows.Close()

			var list []map[string]interface{}
			for rows.Next() {
				var id, tenantID, userID int
				var totalAmount float64
				var status, shippingAddress string
				var createdAt time.Time
				rows.Scan(&id, &tenantID, &userID, &totalAmount, &status, &shippingAddress, &createdAt)
				list = append(list, map[string]interface{}{"id": id, "tenant_id": tenantID, "user_id": userID, "total_amount": totalAmount, "status": status, "shipping_address": shippingAddress, "created_at": createdAt})
			}
			PushAdvancedLog("repository", "OrdersRepo.FindLimit", "SUCCESS", "Query executed", map[string]interface{}{"rows_count": len(list), "duration_ms": time.Since(repoStart).Milliseconds()})

			PushAdvancedLog("controller", "GetOrdersHandler", "SUCCESS", "Request completed", map[string]interface{}{"total_duration_ms": time.Since(start).Milliseconds()})
			c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
		})

		// 5. Payments
		api.GET("/payments", func(c *gin.Context) {
			start := time.Now()
			PushAdvancedLog("controller", "GetPaymentsHandler", "SUCCESS", "Request received for payments", nil)

			repoStart := time.Now()
			rows, err := db.Query("SELECT id, order_id, amount, payment_gateway, status, transaction_ref, created_at FROM payments LIMIT 50")
			if err != nil {
				PushAdvancedLog("repository", "PaymentsRepo.FindLimit", "FAILED", "Database query failed", map[string]interface{}{"error": err.Error()})
				c.JSON(500, gin.H{"success": false, "error": err.Error()})
				return
			}
			defer rows.Close()

			var list []map[string]interface{}
			for rows.Next() {
				var id, orderID int
				var amount float64
				var gateway, status, ref string
				var createdAt time.Time
				rows.Scan(&id, &orderID, &amount, &gateway, &status, &ref, &createdAt)
				list = append(list, map[string]interface{}{"id": id, "order_id": orderID, "amount": amount, "gateway": gateway, "status": status, "transaction_ref": ref, "created_at": createdAt})
			}
			PushAdvancedLog("repository", "PaymentsRepo.FindLimit", "SUCCESS", "Query executed", map[string]interface{}{"rows_count": len(list), "duration_ms": time.Since(repoStart).Milliseconds()})

			PushAdvancedLog("controller", "GetPaymentsHandler", "SUCCESS", "Request completed", map[string]interface{}{"total_duration_ms": time.Since(start).Milliseconds()})
			c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
		})

		// 6. Analytics Summary
		api.GET("/analytics/summary", func(c *gin.Context) {
			start := time.Now()
			PushAdvancedLog("controller", "GetAnalyticsSummaryHandler", "SUCCESS", "Request received for analytics summary", nil)

			repoStart := time.Now()
			var totalUsers, totalProducts, totalOrders int
			var totalRevenue float64
			errU := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
			errP := db.QueryRow("SELECT COUNT(*) FROM products").Scan(&totalProducts)
			errO := db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&totalOrders)
			errR := db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM payments WHERE status='SUCCESS'").Scan(&totalRevenue)

			if errU != nil || errP != nil || errO != nil || errR != nil {
				PushAdvancedLog("repository", "AnalyticsRepo.Summary", "FAILED", "Database aggregated query failed", map[string]interface{}{"error": "query execution error"})
				c.JSON(500, gin.H{"success": false, "error": "internal database error"})
				return
			}
			PushAdvancedLog("repository", "AnalyticsRepo.Summary", "SUCCESS", "Aggregated query executed", map[string]interface{}{"duration_ms": time.Since(repoStart).Milliseconds()})

			PushAdvancedLog("controller", "GetAnalyticsSummaryHandler", "SUCCESS", "Request completed", map[string]interface{}{"total_duration_ms": time.Since(start).Milliseconds()})
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"analytics": gin.H{
					"total_users":    totalUsers,
					"total_products": totalProducts,
					"total_orders":   totalOrders,
					"total_revenue":  totalRevenue,
				},
			})
		})
	}

	log.Println("Starting Massive Enterprise Observability Application on :8080 with Layer-Level Logging...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
