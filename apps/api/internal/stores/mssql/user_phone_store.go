package mssql

import (
	"context"
	"database/sql"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type UserPhoneStore struct {
	db *sql.DB
}

func NewUserPhoneStore(db *sql.DB) *UserPhoneStore {
	return &UserPhoneStore{db: db}
}

func (s *UserPhoneStore) Get(ctx context.Context, id uuid.UUID) (*models.UserPhone, error) {
	query := `
		SELECT id, user_id, phone, phone_digest, is_active, is_verified, is_primary, created_at, updated_at
		FROM users_phones
		WHERE id = ?
	`

	var k models.UserPhone
	var idStr, userIDStr string

	err := s.db.QueryRowContext(ctx, core.Rebind("mssql", query), id.String()).Scan(
		&idStr, &userIDStr, &k.Phone, &k.PhoneDigest, &k.IsActive, &k.IsVerified, &k.IsPrimary, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	k.ID, _ = uuid.Parse(idStr)
	k.UserID, _ = uuid.Parse(userIDStr)

	return &k, nil
}

func (s *UserPhoneStore) GetByPhone(ctx context.Context, phone string) (*models.UserPhone, error) {
	query := `
		SELECT id, user_id, phone, phone_digest, is_active, is_verified, is_primary, created_at, updated_at
		FROM users_phones
		WHERE phone = ?
	`

	var k models.UserPhone
	var idStr, userIDStr string

	err := s.db.QueryRowContext(ctx, core.Rebind("mssql", query), phone).Scan(
		&idStr, &userIDStr, &k.Phone, &k.PhoneDigest, &k.IsActive, &k.IsVerified, &k.IsPrimary, &k.CreatedAt, &k.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	k.ID, _ = uuid.Parse(idStr)
	k.UserID, _ = uuid.Parse(userIDStr)

	return &k, nil
}

func (s *UserPhoneStore) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.UserPhone, error) {
	query := `
		SELECT id, user_id, phone, phone_digest, is_active, is_verified, is_primary, created_at, updated_at
		FROM users_phones
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, core.Rebind("mssql", query), userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []models.UserPhone
	for rows.Next() {
		var k models.UserPhone
		var idStr, uidStr string
		err := rows.Scan(
			&idStr, &uidStr, &k.Phone, &k.PhoneDigest, &k.IsActive, &k.IsVerified, &k.IsPrimary, &k.CreatedAt, &k.UpdatedAt,
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

func (s *UserPhoneStore) Create(ctx context.Context, phone *models.UserPhone) error {
	now := time.Now()
	if phone.CreatedAt.IsZero() {
		phone.CreatedAt = now
	}
    if phone.ID == uuid.Nil {
        phone.ID = uuid.New()
    }

	query := `
		INSERT INTO users_phones (id, user_id, phone, phone_digest, is_active, is_verified, is_primary, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query),
		phone.ID.String(), phone.UserID.String(), phone.Phone, phone.PhoneDigest, phone.IsActive, phone.IsVerified, phone.IsPrimary, phone.CreatedAt, phone.UpdatedAt,
	)
	return err
}

func (s *UserPhoneStore) Update(ctx context.Context, phone *models.UserPhone) error {
	now := time.Now()
	phone.UpdatedAt = &now

	query := `
		UPDATE users_phones
        SET phone = ?, phone_digest = ?, is_active = ?, is_verified = ?, is_primary = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query),
		phone.Phone, phone.PhoneDigest, phone.IsActive, phone.IsVerified, phone.IsPrimary, phone.UpdatedAt, phone.ID.String(),
	)
	return err
}

func (s *UserPhoneStore) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users_phones WHERE id = ?`
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query), id.String())
	return err
}
