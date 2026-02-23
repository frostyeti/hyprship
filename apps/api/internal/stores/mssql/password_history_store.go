package mssql

import (
	"context"
	"database/sql"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/google/uuid"
)

type UserPasswordHistoryStore struct {
	db *sql.DB
}

func NewUserPasswordHistoryStore(db *sql.DB) *UserPasswordHistoryStore {
	return &UserPasswordHistoryStore{db: db}
}

func (s *UserPasswordHistoryStore) ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]string, error) {
	query := `SELECT TOP (?) password_digest FROM user_password_history WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, core.Rebind("mssql", query), limit, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []string
	for rows.Next() {
		var digest string
		if err := rows.Scan(&digest); err != nil {
			return nil, err
		}
		history = append(history, digest)
	}
	return history, nil
}

func (s *UserPasswordHistoryStore) Create(ctx context.Context, userID uuid.UUID, passwordDigest string) error {
	query := `INSERT INTO user_password_history (id, user_id, password_digest, created_at) VALUES (?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query), uuid.New().String(), userID.String(), passwordDigest, time.Now())
	return err
}

func (s *UserPasswordHistoryStore) DeleteOld(ctx context.Context, userID uuid.UUID, keepCount int) error {
	query := `DELETE FROM user_password_history WHERE id IN (
		SELECT id FROM user_password_history WHERE user_id = ? ORDER BY created_at DESC OFFSET ? ROWS
	)`
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query), userID.String(), keepCount)
	return err
}
