package pg

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type UserPasswordAuthStore struct {
	db *sql.DB
}

func NewUserPasswordAuthStore(db *sql.DB) *UserPasswordAuthStore {
	return &UserPasswordAuthStore{db: db}
}

func (s *UserPasswordAuthStore) Get(ctx context.Context, userID uuid.UUID) (*models.UserPasswordAuth, error) {
	var a models.UserPasswordAuth
	var passwordExpiresAt, lastAttemptedAt, otpExpiresAt, updatedAt sql.NullTime
	query := `SELECT user_id, password_digest, password_expires_at, last_attempted_at, attempt_count, otp_digest, otp_expires_at, otp_link_token, is_locked, created_at, updated_at FROM user_password_auth WHERE user_id = $1`
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&a.UserID, &a.PasswordDigest, &passwordExpiresAt, &lastAttemptedAt, &a.AttemptCount,
		&a.OtpDigest, &otpExpiresAt, &a.OtpLinkToken, &a.IsLocked, &a.CreatedAt, &updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("not found")
	} else if err != nil {
		return nil, err
	}

	if passwordExpiresAt.Valid {
		a.PasswordExpiresAt = &passwordExpiresAt.Time
	}
	if lastAttemptedAt.Valid {
		a.LastAttemptedAt = &lastAttemptedAt.Time
	}
	if otpExpiresAt.Valid {
		a.OtpExpiresAt = &otpExpiresAt.Time
	}
	if updatedAt.Valid {
		a.UpdatedAt = &updatedAt.Time
	}
	return &a, nil
}

func (s *UserPasswordAuthStore) Create(ctx context.Context, a *models.UserPasswordAuth) error {
	query := `INSERT INTO user_password_auth (user_id, password_digest, password_expires_at, last_attempted_at, attempt_count, otp_digest, otp_expires_at, otp_link_token, is_locked, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := s.db.ExecContext(ctx, query, a.UserID, a.PasswordDigest, a.PasswordExpiresAt, a.LastAttemptedAt, a.AttemptCount, a.OtpDigest, a.OtpExpiresAt, a.OtpLinkToken, a.IsLocked, a.CreatedAt, a.UpdatedAt)
	return err
}

func (s *UserPasswordAuthStore) Update(ctx context.Context, a *models.UserPasswordAuth) error {
	query := `UPDATE user_password_auth SET password_digest = $1, password_expires_at = $2, last_attempted_at = $3, attempt_count = $4, otp_digest = $5, otp_expires_at = $6, otp_link_token = $7, is_locked = $8, updated_at = $9 WHERE user_id = $10`

	now := time.Now()
	_, err := s.db.ExecContext(ctx, query, a.PasswordDigest, a.PasswordExpiresAt, a.LastAttemptedAt, a.AttemptCount, a.OtpDigest, a.OtpExpiresAt, a.OtpLinkToken, a.IsLocked, now, a.UserID)
	return err
}

func (s *UserPasswordAuthStore) Delete(ctx context.Context, userID uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM user_password_auth WHERE user_id = $1`, userID)
	return err
}
