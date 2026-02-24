package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type userPasskeyStore struct {
	db *sql.DB
}

func NewUserPasskeyStore(db *sql.DB) *userPasskeyStore {
	return &userPasskeyStore{db: db}
}

func (s *userPasskeyStore) Get(ctx context.Context, id uuid.UUID) (*models.UserPasskey, error) {
	query := `
		SELECT id, user_id, credential_id, data, created_at
		FROM user_passkeys
		WHERE id = ?
	`

	var p models.UserPasskey
	var idStr, userIDStr string

	err := s.db.QueryRowContext(ctx, core.Rebind("sqlite", query), id.String()).Scan(
		&idStr, &userIDStr, &p.CredentialID, &p.Data, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	uid, _ := uuid.Parse(userIDStr)
	pid, _ := uuid.Parse(idStr)
	p.ID = pid
	p.UserID = uid

	return &p, nil
}

func (s *userPasskeyStore) GetByCredentialID(ctx context.Context, credentialID []byte) (*models.UserPasskey, error) {
	query := `
		SELECT id, user_id, credential_id, data, created_at
		FROM user_passkeys
		WHERE credential_id = ?
	`

	var p models.UserPasskey
	var idStr, userIDStr string

	err := s.db.QueryRowContext(ctx, core.Rebind("sqlite", query), credentialID).Scan(
		&idStr, &userIDStr, &p.CredentialID, &p.Data, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	uid, _ := uuid.Parse(userIDStr)
	pid, _ := uuid.Parse(idStr)
	p.ID = pid
	p.UserID = uid

	return &p, nil
}

func (s *userPasskeyStore) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.UserPasskey, error) {
	query := `
		SELECT id, user_id, credential_id, data, created_at
		FROM user_passkeys
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, core.Rebind("sqlite", query), userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passkeys []models.UserPasskey
	for rows.Next() {
		var p models.UserPasskey
		var idStr, uidStr string
		err := rows.Scan(
			&idStr, &uidStr, &p.CredentialID, &p.Data, &p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		uid, _ := uuid.Parse(uidStr)
		pid, _ := uuid.Parse(idStr)
		p.ID = pid
		p.UserID = uid
		passkeys = append(passkeys, p)
	}

	return passkeys, rows.Err()
}

func (s *userPasskeyStore) Create(ctx context.Context, passkey *models.UserPasskey) error {
	query := `
		INSERT INTO user_passkeys (id, user_id, credential_id, data, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	now := time.Now()
	if passkey.CreatedAt.IsZero() {
		passkey.CreatedAt = now
	}
	if passkey.ID == uuid.Nil {
		passkey.ID = uuid.New()
	}

	_, err := s.db.ExecContext(ctx, core.Rebind("sqlite", query),
		passkey.ID.String(), passkey.UserID.String(), passkey.CredentialID, passkey.Data, passkey.CreatedAt,
	)

	return err
}

func (s *userPasskeyStore) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM user_passkeys WHERE id = ?`
	_, err := s.db.ExecContext(ctx, core.Rebind("sqlite", query), id.String())
	return err
}
