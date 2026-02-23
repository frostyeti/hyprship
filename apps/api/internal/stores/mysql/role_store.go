package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type RoleStore struct {
	db *sql.DB
}

func NewRoleStore(db *sql.DB) *RoleStore {
	return &RoleStore{db: db}
}

func (s *RoleStore) List(ctx context.Context, opts core.ListOptions) (core.ListResult[models.Role], error) {
	query := `SELECT id, name, description FROM roles`
	args := []any{}
	if opts.Filter != "" {
		query += ` WHERE name_upcase LIKE ?`
		filter := "%" + strings.ToUpper(opts.Filter) + "%"
		args = append(args, filter)
	}

	if opts.Sort != "" {
		query += " ORDER BY " + opts.Sort
		if opts.SortDesc {
			query += " DESC"
		}
	} else {
		query += " ORDER BY name_upcase ASC"
	}

	if opts.PageSize > 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, opts.PageSize, (opts.Page-1)*opts.PageSize)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return core.ListResult[models.Role]{}, err
	}
	defer rows.Close()

	var items []models.Role
	for rows.Next() {
		var r models.Role
		if err := rows.Scan(&r.ID, &r.Name, &r.Description); err != nil {
			return core.ListResult[models.Role]{}, err
		}
		items = append(items, r)
	}

	var total int64
	countQuery := `SELECT COUNT(*) FROM roles`
	countArgs := []any{}
	if opts.Filter != "" {
		countQuery += ` WHERE name_upcase LIKE ?`
		countArgs = args[:1]
	}
	err = s.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return core.ListResult[models.Role]{}, err
	}

	return core.ListResult[models.Role]{Items: items, Total: total}, nil
}

func (s *RoleStore) Get(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var r models.Role
	query := `SELECT id, name, description FROM roles WHERE id = ?`
	err := s.db.QueryRowContext(ctx, query, id).Scan(&r.ID, &r.Name, &r.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *RoleStore) Create(ctx context.Context, role *models.Role) error {
	query := `INSERT INTO roles (id, name, name_upcase, description) VALUES (?, ?, ?, ?)`
	role.NameUpcase = strings.ToUpper(role.Name)
	_, err := s.db.ExecContext(ctx, query, role.ID, role.Name, role.NameUpcase, role.Description)
	return err
}

func (s *RoleStore) Update(ctx context.Context, role *models.Role) error {
	query := `UPDATE roles SET name = ?, name_upcase = ?, description = ? WHERE id = ?`
	role.NameUpcase = strings.ToUpper(role.Name)
	_, err := s.db.ExecContext(ctx, query, role.Name, role.NameUpcase, role.Description, role.ID)
	return err
}

func (s *RoleStore) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM roles WHERE id = ?`, id)
	return err
}

func (s *RoleStore) ListClaims(ctx context.Context, roleID uuid.UUID) ([]models.RoleClaim, error) {
	query := `SELECT id, role_id, type, claim_value FROM role_claims WHERE role_id = ?`
	rows, err := s.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var claims []models.RoleClaim
	for rows.Next() {
		var c models.RoleClaim
		if err := rows.Scan(&c.ID, &c.RoleID, &c.Type, &c.Value); err != nil {
			return nil, err
		}
		claims = append(claims, c)
	}
	return claims, nil
}

func (s *RoleStore) AddClaim(ctx context.Context, claim *models.RoleClaim) error {
	query := `INSERT INTO role_claims (role_id, type, claim_value) VALUES (?, ?, ?)`
	res, err := s.db.ExecContext(ctx, query, claim.RoleID, claim.Type, claim.Value)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		claim.ID = int32(id)
	}
	return err
}

func (s *RoleStore) RemoveClaim(ctx context.Context, claimID int32) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM role_claims WHERE id = ?`, claimID)
	return err
}

func (s *RoleStore) Export(opts core.ListOptions) ([]models.Role, error) {
	res, err := s.List(context.Background(), opts)
	return res.Items, err
}

func (s *RoleStore) Import(items []models.Role) error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query := `INSERT INTO roles (id, name, name_upcase, description) VALUES (?, ?, ?, ?)`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, role := range items {
		role.NameUpcase = strings.ToUpper(role.Name)
		_, err := stmt.Exec(role.ID, role.Name, role.NameUpcase, role.Description)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}
