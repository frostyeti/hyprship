package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserPasswordHistoryStore struct {
	db *sql.DB
}

func NewUserPasswordHistoryStore(db *sql.DB) *UserPasswordHistoryStore {
	return &UserPasswordHistoryStore{db: db}
}

func (s *UserPasswordHistoryStore) ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]string, error) {
	query := `SELECT password_digest FROM user_password_history WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`
	rows, err := s.db.QueryContext(ctx, query, userID, limit)
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
	query := `INSERT INTO user_password_history (id, user_id, password_digest, created_at) VALUES ($1, $2, $3, $4)`
	_, err := s.db.ExecContext(ctx, query, uuid.New(), userID, passwordDigest, time.Now())
	return err
}

func (s *UserPasswordHistoryStore) DeleteOld(ctx context.Context, userID uuid.UUID, keepCount int) error {
	query := `DELETE FROM user_password_history WHERE id IN (
		SELECT id FROM user_password_history WHERE user_id = $1 ORDER BY created_at DESC OFFSET $2
	)`
	_, err := s.db.ExecContext(ctx, query, userID, keepCount)
	return err
}
