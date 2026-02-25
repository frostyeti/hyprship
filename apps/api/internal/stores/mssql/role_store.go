package mssql

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
	fieldMap := map[string]string{
		"id":   "id",
		"name": "name",
	}

	baseQuery := `SELECT id, name, description FROM roles`
	query, args := core.BuildListQuery("mssql", baseQuery, opts, fieldMap)

	rows, err := s.db.QueryContext(ctx, core.Rebind("mssql", query), args...)
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
	countQuery, countArgs := core.BuildCountQuery(countQuery, core.ListOptions{Filter: opts.Filter}, fieldMap)
	err = s.db.QueryRowContext(ctx, core.Rebind("mssql", countQuery), countArgs...).Scan(&total)
	if err != nil {
		return core.ListResult[models.Role]{}, err
	}

	return core.ListResult[models.Role]{Items: items, TotalCount: int(total)}, nil
}

func (s *RoleStore) Get(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var r models.Role
	query := `SELECT CAST(id AS CHAR(36)), name, description FROM roles WHERE id = ?`
	err := s.db.QueryRowContext(ctx, core.Rebind("mssql", query), id.String()).Scan(&r.ID, &r.Name, &r.Description)
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
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query), role.ID.String(), role.Name, role.NameUpcase, role.Description)
	return err
}

func (s *RoleStore) Update(ctx context.Context, role *models.Role) error {
	query := `UPDATE roles SET name = ?, name_upcase = ?, description = ? WHERE id = ?`
	role.NameUpcase = strings.ToUpper(role.Name)
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", query), role.Name, role.NameUpcase, role.Description, role.ID.String())
	return err
}

func (s *RoleStore) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", `DELETE FROM roles WHERE id = ?`), id.String())
	return err
}

func (s *RoleStore) ListClaims(ctx context.Context, roleID uuid.UUID) ([]models.RoleClaim, error) {
	query := `SELECT id, CAST(role_id AS CHAR(36)), type, claim_value FROM role_claims WHERE role_id = ?`
	rows, err := s.db.QueryContext(ctx, core.Rebind("mssql", query), roleID.String())
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
	query := `INSERT INTO role_claims (role_id, type, claim_value) OUTPUT inserted.id VALUES (?, ?, ?)`
	err := s.db.QueryRowContext(ctx, core.Rebind("mssql", query), claim.RoleID.String(), claim.Type, claim.Value).Scan(&claim.ID)
	return err
}

func (s *RoleStore) RemoveClaim(ctx context.Context, claimID int32) error {
	_, err := s.db.ExecContext(ctx, core.Rebind("mssql", `DELETE FROM role_claims WHERE id = ?`), claimID)
	return err
}
