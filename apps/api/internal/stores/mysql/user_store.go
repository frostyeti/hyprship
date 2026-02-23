package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) List(ctx context.Context, opts core.ListOptions) (core.ListResult[models.User], error) {
	query := `SELECT id, primary_email, primary_phone, name, image_uri, is_banned, created_at, updated_at FROM users`
	args := []any{}
	if opts.Filter != "" {
		query += ` WHERE name_upcase LIKE ? OR primary_email_upcase LIKE ?`
		filter := "%" + strings.ToUpper(opts.Filter) + "%"
		args = append(args, filter, filter)
	}

	if opts.Sort != "" {
		query += " ORDER BY " + opts.Sort
		if opts.SortDesc {
			query += " DESC"
		}
	} else {
		query += " ORDER BY created_at DESC"
	}

	if opts.PageSize > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, opts.PageSize, (opts.Page-1)*opts.PageSize)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return core.ListResult[models.User]{}, err
	}
	defer rows.Close()

	var items []models.User
	for rows.Next() {
		var u models.User
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&u.ID, &u.PrimaryEmail, &u.PrimaryPhone, &u.Name, &u.ImageURI, &u.IsBanned, &createdAt, &updatedAt); err != nil {
			return core.ListResult[models.User]{}, err
		}
		if createdAt.Valid {
			u.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			u.UpdatedAt = &t
		}
		items = append(items, u)
	}

	var total int64
	countQuery := `SELECT COUNT(*) FROM users`
	countArgs := []any{}
	if opts.Filter != "" {
		countQuery += ` WHERE name_upcase LIKE ? OR primary_email_upcase LIKE ?`
		countArgs = args[:2]
	}
	err = s.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return core.ListResult[models.User]{}, err
	}

	return core.ListResult[models.User]{Items: items, Total: total}, nil
}

func (s *UserStore) Get(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var u models.User
	var createdAt, updatedAt sql.NullTime
	query := `SELECT id, primary_email, primary_phone, name, image_uri, is_banned, created_at, updated_at FROM users WHERE id = ?`
	err := s.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.PrimaryEmail, &u.PrimaryPhone, &u.Name, &u.ImageURI, &u.IsBanned, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if createdAt.Valid {
		u.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		u.UpdatedAt = &t
	}
	return &u, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	emailUpcase := strings.ToUpper(email)
	var u models.User
	var createdAt, updatedAt sql.NullTime
	query := `SELECT id, primary_email, primary_phone, name, image_uri, is_banned, created_at, updated_at FROM users WHERE primary_email_upcase = ?`
	err := s.db.QueryRowContext(ctx, query, emailUpcase).Scan(&u.ID, &u.PrimaryEmail, &u.PrimaryPhone, &u.Name, &u.ImageURI, &u.IsBanned, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if createdAt.Valid {
		u.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		u.UpdatedAt = &t
	}
	return &u, nil
}

func (s *UserStore) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, primary_email, primary_email_upcase, primary_phone, name, name_upcase, image_uri, is_banned, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	emailUpcase, nameUpcase := "", ""
	if user.PrimaryEmail != nil {
		emailUpcase = strings.ToUpper(*user.PrimaryEmail)
		user.PrimaryEmailUpcase = &emailUpcase
	}
	if user.Name != nil {
		nameUpcase = strings.ToUpper(*user.Name)
		user.NameUpcase = &nameUpcase
	}

	_, err := s.db.ExecContext(ctx, query, user.ID, user.PrimaryEmail, user.PrimaryEmailUpcase, user.PrimaryPhone, user.Name, user.NameUpcase, user.ImageURI, user.IsBanned, user.CreatedAt, user.UpdatedAt)
	return err
}

func (s *UserStore) Update(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET primary_email = ?, primary_email_upcase = ?, primary_phone = ?, name = ?, name_upcase = ?, image_uri = ?, is_banned = ?, updated_at = ? WHERE id = ?`

	emailUpcase, nameUpcase := "", ""
	if user.PrimaryEmail != nil {
		emailUpcase = strings.ToUpper(*user.PrimaryEmail)
		user.PrimaryEmailUpcase = &emailUpcase
	}
	if user.Name != nil {
		nameUpcase = strings.ToUpper(*user.Name)
		user.NameUpcase = &nameUpcase
	}

	user.UpdatedAt = &[]time.Time{time.Now()}[0]

	_, err := s.db.ExecContext(ctx, query, user.PrimaryEmail, user.PrimaryEmailUpcase, user.PrimaryPhone, user.Name, user.NameUpcase, user.ImageURI, user.IsBanned, user.UpdatedAt, user.ID)
	return err
}

func (s *UserStore) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	return err
}

func (s *UserStore) ListClaims(ctx context.Context, userID uuid.UUID) ([]models.UserClaim, error) {
	query := `SELECT id, user_id, type, value FROM user_claims WHERE user_id = ?`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var claims []models.UserClaim
	for rows.Next() {
		var c models.UserClaim
		if err := rows.Scan(&c.ID, &c.UserID, &c.Type, &c.Value); err != nil {
			return nil, err
		}
		claims = append(claims, c)
	}
	return claims, nil
}

func (s *UserStore) AddClaim(ctx context.Context, claim *models.UserClaim) error {
	query := `INSERT INTO user_claims (user_id, type, value) VALUES (?, ?, ?)`
	res, err := s.db.ExecContext(ctx, query, claim.UserID, claim.Type, claim.Value)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		claim.ID = int32(id)
	}
	return err
}

func (s *UserStore) RemoveClaim(ctx context.Context, claimID int32) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM user_claims WHERE id = ?`, claimID)
	return err
}

func (s *UserStore) Export(opts core.ListOptions) ([]models.User, error) {
	res, err := s.List(context.Background(), opts)
	return res.Items, err
}

func (s *UserStore) Import(items []models.User) error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query := `INSERT INTO users (id, primary_email, primary_email_upcase, primary_phone, name, name_upcase, image_uri, is_banned, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, user := range items {
		emailUpcase, nameUpcase := "", ""
		if user.PrimaryEmail != nil {
			emailUpcase = strings.ToUpper(*user.PrimaryEmail)
			user.PrimaryEmailUpcase = &emailUpcase
		}
		if user.Name != nil {
			nameUpcase = strings.ToUpper(*user.Name)
			user.NameUpcase = &nameUpcase
		}

		_, err := stmt.Exec(user.ID, user.PrimaryEmail, user.PrimaryEmailUpcase, user.PrimaryPhone, user.Name, user.NameUpcase, user.ImageURI, user.IsBanned, user.CreatedAt, user.UpdatedAt)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
