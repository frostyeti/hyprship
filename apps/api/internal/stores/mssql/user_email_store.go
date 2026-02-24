package mssql

import (
	"context"
	"database/sql"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type UserEmailStore struct {
	db *sql.DB
}

func NewUserEmailStore(db *sql.DB) *UserEmailStore {
	return &UserEmailStore{db: db}
}

func (s *UserEmailStore) Get(ctx context.Context, id uuid.UUID) (*models.UserEmail, error) {
	query := `
		SELECT id, user_id, email, email_upcase, email_upcase_digest, is_active, is_verified, is_primary, created_at, updated_at
		FROM users_emails
		WHERE id = ?
	`

	var k models.UserEmail
	var idStr, userIDStr string

	err := s.db.QueryRowContext(ctx, core.Rebind("mssql", query), id.String()).Scan(
		&idStr, &userIDStr, &k.Email, &k.EmailUpcase, &k.EmailUpcaseDigest, &k.IsActive, &k.IsVerified, &k.IsPrimary, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	k.ID, _ = uuid.Parse(idStr)
	k.UserID, _ = uuid.Parse(userIDStr)

	return &k, nil
}

func (s *UserEmailStore) GetByEmail(ctx context.Context, email string) (*models.UserEmail, error) {
	query := `
		SELECT id, user_id, email, email_upcase, email_upcase_digest, is_active, is_verified, is_primary, created_at, updated_at
		FROM users_emails
		WHERE email_upcase = ?
	`

	var k models.UserEmail
	var idStr, userIDStr string

	err := s.db.QueryRowContext(ctx, core.Rebind("mssql", query), email).Scan(
		&idStr, &userIDStr, &k.Email, &k.EmailUpcase, &k.EmailUpcaseDigest, &k.IsActive, &k.IsVerified, &k.IsPrimary, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	k.ID, _ = uuid.Parse(idStr)
	k.UserID, _ = uuid.Parse(userIDStr)

	return &k, nil
}

func (s *UserEmailStore) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.UserEmail, error) {
	query := `
		SELECT id, user_id, email, email_upcase, email_upcase_digest, is_active, is_verified, is_primary, created_at, updated_at
		FROM users_emails
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, core.Rebind("mssql", query), userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []models.UserEmail
	for rows.Next() {
		var k models.UserEmail
		var idStr, uidStr string
		err := rows.Scan(
			&idStr, &uidStr, &k.Email, &k.EmailUpcase, &k.EmailUpcaseDigest, &k.IsActive, &k.IsVerified, &k.IsPrimary, &k.CreatedAt, &k.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		k.ID, _ = uuid.Parse(idStr)
		k.UserID, _ = uuid.Parse(uidStr)
		keys = append(keys, k)
	}

	return keys, rows.Err()
}

func (s *UserEmailStore) Create(ctx context.Context, email *models.UserEmail) error {
	now := time.Now()
	if email.CreatedAt.IsZero() {
		email.CreatedAt = now
	}
	if email.ID == uuid.Nil {
		email.ID = uuid.New()
	}

	query := `
		INSERT INTO users_emails (id, user_id, email, email_upcase, email_upcase_digest, is_active, is_verified, is_primary, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query),
		email.ID.String(), email.UserID.String(), email.Email, email.EmailUpcase, email.EmailUpcaseDigest, email.IsActive, email.IsVerified, email.IsPrimary, email.CreatedAt, email.UpdatedAt,
	)
	return err
}

func (s *UserEmailStore) Update(ctx context.Context, email *models.UserEmail) error {
	now := time.Now()
	email.UpdatedAt = &now

	query := `
		UPDATE users_emails
        SET email = ?, email_upcase = ?, email_upcase_digest = ?, is_active = ?, is_verified = ?, is_primary = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query),
		email.Email, email.EmailUpcase, email.EmailUpcaseDigest, email.IsActive, email.IsVerified, email.IsPrimary, email.UpdatedAt, email.ID.String(),
	)
	return err
}

func (s *UserEmailStore) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users_emails WHERE id = ?`
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query), id.String())
	return err
}
