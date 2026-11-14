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
	rows, err := r.db.QueryContext(ctx, "SELECT id, tenant_id, name, email, role, created_at FROM users")
	if err != nil {
		infrastructure.PushLog("repository", "UserRepository.GetAll", "FAILED", err.Error(), nil)
		return nil, err
	}
	defer rows.Close()

	var list []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.TenantID, &u.Name, &u.Email, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	infrastructure.PushLog("repository", "UserRepository.GetAll", "SUCCESS", "Fetched users", map[string]interface{}{"count": len(list), "duration_ms": time.Since(start).Milliseconds()})
	return list, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	start := time.Now()
	var u domain.User
	err := r.db.QueryRowContext(ctx, "SELECT id, tenant_id, name, email, role, created_at FROM users WHERE id = $1", id).
		Scan(&u.ID, &u.TenantID, &u.Name, &u.Email, &u.Role, &u.CreatedAt)
	if err != nil {
		infrastructure.PushLog("repository", "UserRepository.GetByID", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return nil, err
	}
	infrastructure.PushLog("repository", "UserRepository.GetByID", "SUCCESS", "Fetched user by ID", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	start := time.Now()
	query := "INSERT INTO users (tenant_id, name, email, role, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, u.TenantID, u.Name, u.Email, u.Role, time.Now()).Scan(&u.ID)
	if err != nil {
		infrastructure.PushLog("repository", "UserRepository.Create", "FAILED", err.Error(), nil)
		return err
	}
	infrastructure.PushLog("repository", "UserRepository.Create", "SUCCESS", "Created user", map[string]interface{}{"id": u.ID, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) error {
	start := time.Now()
	query := "UPDATE users SET tenant_id = $1, name = $2, email = $3, role = $4 WHERE id = $5"
	_, err := r.db.ExecContext(ctx, query, u.TenantID, u.Name, u.Email, u.Role, u.ID)
	if err != nil {
		infrastructure.PushLog("repository", "UserRepository.Update", "FAILED", err.Error(), map[string]interface{}{"id": u.ID})
		return err
	}
	infrastructure.PushLog("repository", "UserRepository.Update", "SUCCESS", "Updated user", map[string]interface{}{"id": u.ID, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	start := time.Now()
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		infrastructure.PushLog("repository", "UserRepository.Delete", "FAILED", err.Error(), map[string]interface{}{"id": id})
		return err
	}
	infrastructure.PushLog("repository", "UserRepository.Delete", "SUCCESS", "Deleted user", map[string]interface{}{"id": id, "duration_ms": time.Since(start).Milliseconds()})
	return nil
}
