package mssql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type UserSessionStore struct {
	db *sql.DB
}

func NewUserSessionStore(db *sql.DB) *UserSessionStore {
	return &UserSessionStore{db: db}
}

func (s *UserSessionStore) Get(ctx context.Context, sessionID uuid.UUID) (*models.UserSession, error) {
	var session models.UserSession
	var ipAddress, userAgent *string
	var updatedAt sql.NullTime
	var idStr, userIDStr string
	query := `SELECT CAST(id AS CHAR(36)), CAST(user_id AS CHAR(36)), expires_at, token, ip_address, user_agent, created_at, updated_at FROM user_sessions WHERE id = ?`
	
	err := s.db.QueryRowContext(ctx, core.Rebind("mssql", query), sessionID.String()).Scan(
		&idStr, &userIDStr, &session.ExpiresAt, &session.Token,
		&ipAddress, &userAgent, &session.CreatedAt, &updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("not found")
	} else if err != nil {
		return nil, err
	}
	
	session.ID, _ = uuid.Parse(idStr)
	session.UserID, _ = uuid.Parse(userIDStr)
	session.IPAddress = ipAddress
	session.UserAgent = userAgent
	if updatedAt.Valid {
		session.UpdatedAt = &updatedAt.Time
	}
	return &session, nil
}

func (s *UserSessionStore) GetByToken(ctx context.Context, token string) (*models.UserSession, error) {
	var session models.UserSession
	var ipAddress, userAgent *string
	var updatedAt sql.NullTime
	var idStr, userIDStr string
	query := `SELECT CAST(id AS CHAR(36)), CAST(user_id AS CHAR(36)), expires_at, token, ip_address, user_agent, created_at, updated_at FROM user_sessions WHERE token = ?`
	
	err := s.db.QueryRowContext(ctx, core.Rebind("mssql", query), token).Scan(
		&idStr, &userIDStr, &session.ExpiresAt, &session.Token,
		&ipAddress, &userAgent, &session.CreatedAt, &updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("not found")
	} else if err != nil {
		return nil, err
	}
	
	session.ID, _ = uuid.Parse(idStr)
	session.UserID, _ = uuid.Parse(userIDStr)
	session.IPAddress = ipAddress
	session.UserAgent = userAgent
	if updatedAt.Valid {
		session.UpdatedAt = &updatedAt.Time
	}
	return &session, nil
}

func (s *UserSessionStore) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.UserSession, error) {
	query := `SELECT CAST(id AS CHAR(36)), CAST(user_id AS CHAR(36)), expires_at, token, ip_address, user_agent, created_at, updated_at FROM user_sessions WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, core.Rebind("mssql", query), userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.UserSession
	for rows.Next() {
		var session models.UserSession
		var ipAddress, userAgent *string
		var updatedAt sql.NullTime
		var idStr, userIDStr string
		if err := rows.Scan(&idStr, &userIDStr, &session.ExpiresAt, &session.Token, &ipAddress, &userAgent, &session.CreatedAt, &updatedAt); err != nil {
			return nil, err
		}
		session.ID, _ = uuid.Parse(idStr)
		session.UserID, _ = uuid.Parse(userIDStr)
		session.IPAddress = ipAddress
		session.UserAgent = userAgent
		if updatedAt.Valid {
			session.UpdatedAt = &updatedAt.Time
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (s *UserSessionStore) Create(ctx context.Context, session *models.UserSession) error {
	query := `INSERT INTO user_sessions (id, user_id, expires_at, token, ip_address, user_agent, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query), session.ID.String(), session.UserID.String(), session.ExpiresAt, session.Token, session.IPAddress, session.UserAgent, session.CreatedAt, session.UpdatedAt)
	return err
}

func (s *UserSessionStore) Update(ctx context.Context, session *models.UserSession) error {
	query := `UPDATE user_sessions SET expires_at = ?, token = ?, ip_address = ?, user_agent = ?, updated_at = ? WHERE id = ?`

	now := time.Now()
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query), session.ExpiresAt, session.Token, session.IPAddress, session.UserAgent, now, session.ID.String())
	return err
}

func (s *UserSessionStore) Delete(ctx context.Context, sessionID uuid.UUID) error {
	query := `DELETE FROM user_sessions WHERE id = ?`
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query), sessionID.String())
	return err
}

func (s *UserSessionStore) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_sessions WHERE user_id = ?`
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query), userID.String())
	return err
}
