package sqlite

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
	var passwordExpiresAt, lastAttemptedAt, otpExpiresAt, updatedAt sql.NullInt64
	query := `SELECT user_id, password_digest, password_expires_at, last_attempted_at, attempt_count, otp_digest, otp_expires_at, otp_link_token, is_locked, created_at, updated_at FROM user_password_auth WHERE user_id = ?`
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
		t := time.Unix(passwordExpiresAt.Int64, 0)
		a.PasswordExpiresAt = &t
	}
	if lastAttemptedAt.Valid {
		t := time.Unix(lastAttemptedAt.Int64, 0)
		a.LastAttemptedAt = &t
	}
	if otpExpiresAt.Valid {
		t := time.Unix(otpExpiresAt.Int64, 0)
		a.OtpExpiresAt = &t
	}
	if updatedAt.Valid {
		t := time.Unix(updatedAt.Int64, 0)
		a.UpdatedAt = &t
	}
	return &a, nil
}

func (s *UserPasswordAuthStore) Create(ctx context.Context, a *models.UserPasswordAuth) error {
	query := `INSERT INTO user_password_auth (user_id, password_digest, password_expires_at, last_attempted_at, attempt_count, otp_digest, otp_expires_at, otp_link_token, is_locked, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	var passwordExpiresAt, lastAttemptedAt, otpExpiresAt, updatedAt *int64
	if a.PasswordExpiresAt != nil {
		t := a.PasswordExpiresAt.Unix()
		passwordExpiresAt = &t
	}
	if a.LastAttemptedAt != nil {
		t := a.LastAttemptedAt.Unix()
		lastAttemptedAt = &t
	}
	if a.OtpExpiresAt != nil {
		t := a.OtpExpiresAt.Unix()
		otpExpiresAt = &t
	}
	if a.UpdatedAt != nil {
		t := a.UpdatedAt.Unix()
		updatedAt = &t
	}

	_, err := s.db.ExecContext(ctx, query, a.UserID, a.PasswordDigest, passwordExpiresAt, lastAttemptedAt, a.AttemptCount, a.OtpDigest, otpExpiresAt, a.OtpLinkToken, a.IsLocked, a.CreatedAt.Unix(), updatedAt)
	return err
}

func (s *UserPasswordAuthStore) Update(ctx context.Context, a *models.UserPasswordAuth) error {
	query := `UPDATE user_password_auth SET password_digest = ?, password_expires_at = ?, last_attempted_at = ?, attempt_count = ?, otp_digest = ?, otp_expires_at = ?, otp_link_token = ?, is_locked = ?, updated_at = ? WHERE user_id = ?`

	var passwordExpiresAt, lastAttemptedAt, otpExpiresAt, updatedAt *int64
	if a.PasswordExpiresAt != nil {
		t := a.PasswordExpiresAt.Unix()
		passwordExpiresAt = &t
	}
	if a.LastAttemptedAt != nil {
		t := a.LastAttemptedAt.Unix()
		lastAttemptedAt = &t
	}
	if a.OtpExpiresAt != nil {
		t := a.OtpExpiresAt.Unix()
		otpExpiresAt = &t
	}
	now := time.Now().Unix()
	updatedAt = &now

	_, err := s.db.ExecContext(ctx, query, a.PasswordDigest, passwordExpiresAt, lastAttemptedAt, a.AttemptCount, a.OtpDigest, otpExpiresAt, a.OtpLinkToken, a.IsLocked, updatedAt, a.UserID)
	return err
}

func (s *UserPasswordAuthStore) Delete(ctx context.Context, userID uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM user_password_auth WHERE user_id = ?`, userID)
	return err
}
