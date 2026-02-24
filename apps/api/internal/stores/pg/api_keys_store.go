package pg

import (
	"context"
	"database/sql"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type UserAPIKeyStore struct {
	db *sql.DB
}

func NewUserAPIKeyStore(db *sql.DB) *UserAPIKeyStore {
	return &UserAPIKeyStore{db: db}
}

func (s *UserAPIKeyStore) Get(ctx context.Context, id int32) (*models.UserAPIKey, error) {
	query := `
		SELECT id, user_id, name, name_upcase, key_hint, key_digest, expires_at, is_locked, is_revoked, comment, created_at, updated_at
		FROM user_api_keys
		WHERE id = ?
	`

	var k models.UserAPIKey
	var userID string

	err := s.db.QueryRowContext(ctx, core.Rebind("pg", query), id).Scan(
		&k.ID, &userID, &k.Name, &k.NameUpcase, &k.KeyHint, &k.KeyDigest, &k.ExpiresAt, &k.IsLocked, &k.IsRevoked, &k.Comment, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	k.UserID = uid

	return &k, nil
}

func (s *UserAPIKeyStore) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.UserAPIKey, error) {
	query := `
		SELECT id, user_id, name, name_upcase, key_hint, key_digest, expires_at, is_locked, is_revoked, comment, created_at, updated_at
		FROM user_api_keys
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, core.Rebind("pg", query), userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []models.UserAPIKey
	for rows.Next() {
		var k models.UserAPIKey
		var uidStr string
		err := rows.Scan(
			&k.ID, &uidStr, &k.Name, &k.NameUpcase, &k.KeyHint, &k.KeyDigest, &k.ExpiresAt, &k.IsLocked, &k.IsRevoked, &k.Comment, &k.CreatedAt, &k.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		uid, _ := uuid.Parse(uidStr)
		k.UserID = uid
		keys = append(keys, k)
	}

	return keys, rows.Err()
}

func (s *UserAPIKeyStore) Create(ctx context.Context, apiKey *models.UserAPIKey) error {
	now := time.Now()
	if apiKey.CreatedAt.IsZero() {
		apiKey.CreatedAt = now
	}

	query := `
		INSERT INTO user_api_keys (user_id, name, name_upcase, key_hint, key_digest, expires_at, is_locked, is_revoked, comment, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	err := s.db.QueryRowContext(ctx, core.Rebind("pg", query),
		apiKey.UserID.String(), apiKey.Name, apiKey.NameUpcase, apiKey.KeyHint, apiKey.KeyDigest, apiKey.ExpiresAt, apiKey.IsLocked, apiKey.IsRevoked, apiKey.Comment, apiKey.CreatedAt, apiKey.UpdatedAt,
	).Scan(&apiKey.ID)
	return err
}

func (s *UserAPIKeyStore) Delete(ctx context.Context, id int32) error {
	query := `DELETE FROM user_api_keys WHERE id = ?`
	_, err := s.db.ExecContext(ctx, core.Rebind("pg", query), id)
	return err
}
