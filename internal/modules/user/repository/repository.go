package repository

import (
	"context"
	"database/sql"
	"time"
	"otel-kibana-grafana-test/internal/domain"
	"otel-kibana-grafana-test/internal/infrastructure"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, "SELECT id, tenant_id, name, email, status, created_at FROM users LIMIT 100")
	if err != nil {
		infrastructure.PushLog("repository", "UserRepo.GetAll", "FAILED", err.Error(), nil)
		return nil, err
	}
	defer rows.Close()

	var list []domain.User
	for rows.Next() {
		var u domain.User
		rows.Scan(&u.ID, &u.TenantID, &u.Name, &u.Email, &u.Status, &u.CreatedAt)
		list = append(list, u)
	}
	infrastructure.PushLog("repository", "UserRepo.GetAll", "SUCCESS", "Fetched users", map[string]interface{}{"count": len(list), "duration_ms": time.Since(start).Milliseconds()})
	return list, nil
}

func (r *UserRepository) Create(ctx context.Context, tenantID int, name, email string) (*domain.User, error) {
	start := time.Now()
	var u domain.User
	err := r.db.QueryRowContext(ctx, "INSERT INTO users (tenant_id, name, email, status) VALUES ($1, $2, $3, 'ACTIVE') RETURNING id, tenant_id, name, email, status, created_at", tenantID, name, email).
		Scan(&u.ID, &u.TenantID, &u.Name, &u.Email, &u.Status, &u.CreatedAt)
	if err != nil {
		infrastructure.PushLog("repository", "UserRepo.Create", "FAILED", err.Error(), nil)
		return nil, err
	}
	infrastructure.PushLog("repository", "UserRepo.Create", "SUCCESS", "User created", map[string]interface{}{"id": u.ID, "duration_ms": time.Since(start).Milliseconds()})
	return &u, nil
}
