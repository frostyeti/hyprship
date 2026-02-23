package pg

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
	query := `SELECT id, user_id, expires_at, token, ip_address, user_agent, created_at, updated_at FROM user_sessions WHERE id = $1`
	
	err := s.db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID, &session.UserID, &session.ExpiresAt, &session.Token,
		&ipAddress, &userAgent, &session.CreatedAt, &updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("not found")
	} else if err != nil {
		return nil, err
	}
	
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
	query := `SELECT id, user_id, expires_at, token, ip_address, user_agent, created_at, updated_at FROM user_sessions WHERE token = $1`
	
	err := s.db.QueryRowContext(ctx, query, token).Scan(
		&session.ID, &session.UserID, &session.ExpiresAt, &session.Token,
		&ipAddress, &userAgent, &session.CreatedAt, &updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("not found")
	} else if err != nil {
		return nil, err
	}
	
	session.IPAddress = ipAddress
	session.UserAgent = userAgent
	if updatedAt.Valid {
		session.UpdatedAt = &updatedAt.Time
	}
	return &session, nil
}

func (s *UserSessionStore) ListByUser(ctx context.Context, userID uuid.UUID) ([]models.UserSession, error) {
	query := `SELECT id, user_id, expires_at, token, ip_address, user_agent, created_at, updated_at FROM user_sessions WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.UserSession
	for rows.Next() {
		var session models.UserSession
		var ipAddress, userAgent *string
		var updatedAt sql.NullTime
		if err := rows.Scan(&session.ID, &session.UserID, &session.ExpiresAt, &session.Token, &ipAddress, &userAgent, &session.CreatedAt, &updatedAt); err != nil {
			return nil, err
		}
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
	query := `INSERT INTO user_sessions (id, user_id, expires_at, token, ip_address, user_agent, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := s.db.ExecContext(ctx, query, session.ID, session.UserID, session.ExpiresAt, session.Token, session.IPAddress, session.UserAgent, session.CreatedAt, session.UpdatedAt)
	return err
}

func (s *UserSessionStore) Update(ctx context.Context, session *models.UserSession) error {
	query := `UPDATE user_sessions SET expires_at = $1, token = $2, ip_address = $3, user_agent = $4, updated_at = $5 WHERE id = $6`

	now := time.Now()
	_, err := s.db.ExecContext(ctx, query, session.ExpiresAt, session.Token, session.IPAddress, session.UserAgent, now, session.ID)
	return err
}

func (s *UserSessionStore) Delete(ctx context.Context, sessionID uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM user_sessions WHERE id = $1`, sessionID)
	return err
}

func (s *UserSessionStore) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM user_sessions WHERE user_id = $1`, userID)
	return err
}
