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
			semconv.ServiceNameKey.String("clean-enterprise-app"),
		)),
	)
	otel.SetTracerProvider(tp)
	tracer = otel.Tracer("clean-arch-tracer")
	return tp, nil
}

// ==========================================
// INFRASTRUCTURE LAYER: Elasticsearch Logger
// ==========================================

func PushLog(layer, component, status, message string, payload map[string]interface{}) {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://elasticsearch:9200"
	}

	logEntry := map[string]interface{}{
		"@timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		"service":    "clean-enterprise-app",
		"layer":      layer,
		"component":  component,
		"status":     status,
		"message":    message,
		"details":    payload,
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

// ==========================================
// DOMAIN / ENTITY LAYER
// ==========================================

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	TotalAmount  float64   `json:"total_amount"`
	Status       string    `json:"status"`
	ShippingAddr string    `json:"shipping_address"`
	CreatedAt    time.Time `json:"created_at"`
}

// ==========================================
// REPOSITORY LAYER (Data Access)
// ==========================================

type UserRepository interface {
	FindAll(ctx context.Context) ([]User, error)
	Create(ctx context.Context, name, email string) (*User, error)
}

type OrderRepository interface {
	FindAll(ctx context.Context) ([]Order, error)
	Create(ctx context.Context, userID int, totalAmount float64, addr string) (*Order, error)
}

type postgresUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) FindAll(ctx context.Context) ([]User, error) {
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, email, status, created_at FROM users LIMIT 50")
	if err != nil {
		PushLog("repository", "UserRepository.FindAll", "FAILED", err.Error(), map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Status, &u.CreatedAt); err != nil {
			continue
		}
		users = append(users, u)
	}
	PushLog("repository", "UserRepository.FindAll", "SUCCESS", "Fetched users successfully", map[string]interface{}{"count": len(users), "duration_ms": time.Since(start).Milliseconds()})
	return users, nil
}

func (r *postgresUserRepository) Create(ctx context.Context, name, email string) (*User, error) {
	start := time.Now()
	var u User
	err := r.db.QueryRowContext(ctx, "INSERT INTO users (tenant_id, name, email, status) VALUES (1, $1, $2, 'ACTIVE') RETURNING id, name, email, status, created_at", name, email).
		Scan(&u.ID, &u.Name, &u.Email, &u.Status, &u.CreatedAt)
	if err != nil {
		PushLog("repository", "UserRepository.Create", "FAILED", err.Error(), map[string]interface{}{"email": email, "duration_ms": time.Since(start).Milliseconds()})
		return nil, err
	}
	PushLog("repository", "UserRepository.Create", "SUCCESS", "User created in database", map[string]interface{}{"user_id": u.ID, "duration_ms": time.Since(start).Milliseconds()})
	return &u, nil
}

type postgresOrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &postgresOrderRepository{db: db}
}

func (r *postgresOrderRepository) FindAll(ctx context.Context) ([]Order, error) {
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, total_amount, status, shipping_address, created_at FROM orders LIMIT 50")
	if err != nil {
		PushLog("repository", "OrderRepository.FindAll", "FAILED", err.Error(), map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.ShippingAddr, &o.CreatedAt); err != nil {
			continue
		}
		orders = append(orders, o)
	}
	PushLog("repository", "OrderRepository.FindAll", "SUCCESS", "Fetched orders successfully", map[string]interface{}{"count": len(orders), "duration_ms": time.Since(start).Milliseconds()})
	return orders, nil
}

func (r *postgresOrderRepository) Create(ctx context.Context, userID int, totalAmount float64, addr string) (*Order, error) {
	start := time.Now()
	var o Order
	err := r.db.QueryRowContext(ctx, "INSERT INTO orders (tenant_id, user_id, total_amount, status, shipping_address) VALUES (1, $1, $2, 'COMPLETED', $3) RETURNING id, user_id, total_amount, status, shipping_address, created_at", userID, totalAmount, addr).
		Scan(&o.ID, &o.UserID, &o.TotalAmount, &o.Status, &o.ShippingAddr, &o.CreatedAt)
	if err != nil {
		PushLog("repository", "OrderRepository.Create", "FAILED", err.Error(), map[string]interface{}{"user_id": userID, "duration_ms": time.Since(start).Milliseconds()})
		return nil, err
	}
	PushLog("repository", "OrderRepository.Create", "SUCCESS", "Order created in database", map[string]interface{}{"order_id": o.ID, "duration_ms": time.Since(start).Milliseconds()})
	return &o, nil
}

// ==========================================
// SERVICE LAYER (Business Logic)
// ==========================================

type UserService interface {
	GetAllUsers(ctx context.Context) ([]User, error)
	RegisterUser(ctx context.Context, name, email string) (*User, error)
}

type OrderService interface {
	GetAllOrders(ctx context.Context) ([]Order, error)
	PlaceOrder(ctx context.Context, userID int, amount float64, addr string) (*Order, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{userRepo: repo}
}

func (s *userService) GetAllUsers(ctx context.Context) ([]User, error) {
	start := time.Now()
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		PushLog("service", "UserService.GetAllUsers", "FAILED", err.Error(), nil)
		return nil, err
	}
	PushLog("service", "UserService.GetAllUsers", "SUCCESS", "Business logic executed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return users, nil
}

func (s *userService) RegisterUser(ctx context.Context, name, email string) (*User, error) {
	start := time.Now()
	u, err := s.userRepo.Create(ctx, name, email)
	if err != nil {
		PushLog("service", "UserService.RegisterUser", "FAILED", err.Error(), nil)
		return nil, err
	}
	PushLog("service", "UserService.RegisterUser", "SUCCESS", "User registered successfully", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return u, nil
}

type orderService struct {
	orderRepo OrderRepository
}

func NewOrderService(repo OrderRepository) OrderService {
	return &orderService{orderRepo: repo}
}

func (s *orderService) GetAllOrders(ctx context.Context) ([]Order, error) {
	start := time.Now()
	orders, err := s.orderRepo.FindAll(ctx)
	if err != nil {
		PushLog("service", "OrderService.GetAllOrders", "FAILED", err.Error(), nil)
		return nil, err
	}
	PushLog("service", "OrderService.GetAllOrders", "SUCCESS", "Business logic executed", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return orders, nil
}

func (s *orderService) PlaceOrder(ctx context.Context, userID int, amount float64, addr string) (*Order, error) {
	start := time.Now()
	o, err := s.orderRepo.Create(ctx, userID, amount, addr)
	if err != nil {
		PushLog("service", "OrderService.PlaceOrder", "FAILED", err.Error(), nil)
		return nil, err
	}
	PushLog("service", "OrderService.PlaceOrder", "SUCCESS", "Order placed successfully", map[string]interface{}{"duration_ms": time.Since(start).Milliseconds()})
	return o, nil
}

// ==========================================
// CONTROLLER / TRANSPORT LAYER (HTTP Handlers)
// ==========================================

type UserController struct {
	userService UserService
}

func NewUserController(r *gin.RouterGroup, service UserService) {
	uc := &UserController{userService: service}
	r.GET("/users", uc.GetUsers)
	r.POST("/users", uc.CreateUser)
}

func (ctrl *UserController) GetUsers(c *gin.Context) {
	start := time.Now()
	PushLog("controller", "UserController.GetUsers", "SUCCESS", "HTTP GET /users received", nil)

	users, err := ctrl.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	PushLog("controller", "UserController.GetUsers", "SUCCESS", "HTTP request completed", map[string]interface{}{"total_duration_ms": time.Since(start).Milliseconds()})
	c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	start := time.Now()
	var input struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		PushLog("controller", "UserController.CreateUser", "FAILED", "Invalid request body", map[string]interface{}{"error": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.userService.RegisterUser(c.Request.Context(), input.Name, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	PushLog("controller", "UserController.CreateUser", "SUCCESS", "User created successfully via HTTP", map[string]interface{}{"total_duration_ms": time.Since(start).Milliseconds()})
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": user})
}

type OrderController struct {
	orderService OrderService
}

func NewOrderController(r *gin.RouterGroup, service OrderService) {
	oc := &OrderController{orderService: service}
	r.GET("/orders", oc.GetOrders)
	r.POST("/orders", oc.CreateOrder)
}

func (ctrl *OrderController) GetOrders(c *gin.Context) {
	start := time.Now()
	PushLog("controller", "OrderController.GetOrders", "SUCCESS", "HTTP GET /orders received", nil)

	orders, err := ctrl.orderService.GetAllOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	PushLog("controller", "OrderController.GetOrders", "SUCCESS", "HTTP request completed", map[string]interface{}{"total_duration_ms": time.Since(start).Milliseconds()})
	c.JSON(http.StatusOK, gin.H{"success": true, "data": orders})
}

func (ctrl *OrderController) CreateOrder(c *gin.Context) {
	start := time.Now()
	var input struct {
		UserID       int     `json:"user_id" binding:"required"`
		TotalAmount  float64 `json:"total_amount" binding:"required"`
		ShippingAddr string  `json:"shipping_address" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		PushLog("controller", "OrderController.CreateOrder", "FAILED", "Invalid request body", map[string]interface{}{"error": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := ctrl.orderService.PlaceOrder(c.Request.Context(), input.UserID, input.TotalAmount, input.ShippingAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	PushLog("controller", "OrderController.CreateOrder", "SUCCESS", "Order placed successfully via HTTP", map[string]interface{}{"total_duration_ms": time.Since(start).Milliseconds()})
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": order})
}

// ==========================================
// MAIN BOOTSTRAPPER
// ==========================================

func main() {
	tp, err := initTracer()
	if err != nil {
		log.Printf("Warning: Failed to init tracer: %v", err)
	} else {
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer: %v", err)
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
		log.Printf("Attempt %d: Waiting for database...", i)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Auto-migration tables
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, tenant_id INT DEFAULT 1, name VARCHAR(255), email VARCHAR(255) UNIQUE, status VARCHAR(50), created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`)
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS orders (id SERIAL PRIMARY KEY, tenant_id INT DEFAULT 1, user_id INT REFERENCES users(id), total_amount NUMERIC(12,2), status VARCHAR(50), shipping_address TEXT, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`)
	PushLog("system", "Migration", "SUCCESS", "Clean architecture tables migrated", nil)

	// Dependency Injection Wiring
	userRepo := NewUserRepository(db)
	orderRepo := NewOrderRepository(db)

	userService := NewUserService(userRepo)
	orderService := NewOrderService(orderRepo)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("clean-enterprise-app"))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "architecture": "clean-layered"})
	})

	api := r.Group("/api/v1")
	NewUserController(api, userService)
	NewOrderController(api, orderService)

	log.Println("Starting Clean Architecture Enterprise App on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server crashed: %v", err)
	}
}
