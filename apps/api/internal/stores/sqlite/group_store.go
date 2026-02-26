package sqlite

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

type GroupStore struct {
	db *sql.DB
}

func NewGroupStore(db *sql.DB) *GroupStore {
	return &GroupStore{db: db}
}

func (s *GroupStore) List(ctx context.Context, opts core.ListOptions) (core.ListResult[models.Group], error) {
	fieldMap := map[string]string{
		"id":        "id",
		"name":      "name_upcase",
		"email":     "email_upcase",
		"isActive":  "is_active",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}

	baseQuery := `SELECT id, name, description, image_uri, email, is_active, created_at, updated_at FROM "groups"`
	query, args := core.BuildListQuery("sqlite3", baseQuery, opts, fieldMap)

	rows, err := s.db.QueryContext(ctx, core.Rebind("sqlite3", query), args...)
	if err != nil {
		return core.ListResult[models.Group]{}, err
	}
	defer rows.Close()

	var items []models.Group
	for rows.Next() {
		var g models.Group
		var createdAt, updatedAt sql.NullInt64
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.ImageURI, &g.Email, &g.IsActive, &createdAt, &updatedAt); err != nil {
			return core.ListResult[models.Group]{}, err
		}
		if createdAt.Valid {
			g.CreatedAt = time.Unix(createdAt.Int64, 0).UTC()
		}
		if updatedAt.Valid {
			t := time.Unix(updatedAt.Int64, 0).UTC()
			g.UpdatedAt = &t
		}
		items = append(items, g)
	}

	var total int64
	countQuery := `SELECT COUNT(*) FROM "groups"`
	countQuery, countArgs := core.BuildCountQuery(countQuery, core.ListOptions{Filter: opts.Filter}, fieldMap)
	err = s.db.QueryRowContext(ctx, core.Rebind("sqlite3", countQuery), countArgs...).Scan(&total)
	if err != nil {
		return core.ListResult[models.Group]{}, err
	}

	return core.ListResult[models.Group]{Items: items, TotalCount: int(total)}, nil
}

func (s *GroupStore) Get(ctx context.Context, id uuid.UUID) (*models.Group, error) {
	var g models.Group
	var createdAt, updatedAt sql.NullInt64
	query := `SELECT id, name, description, image_uri, email, is_active, created_at, updated_at FROM "groups" WHERE id = ?`
	err := s.db.QueryRowContext(ctx, core.Rebind("sqlite", query), id).Scan(&g.ID, &g.Name, &g.Description, &g.ImageURI, &g.Email, &g.IsActive, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if createdAt.Valid {
		g.CreatedAt = time.Unix(createdAt.Int64, 0).UTC()
	}
	if updatedAt.Valid {
		t := time.Unix(updatedAt.Int64, 0).UTC()
		g.UpdatedAt = &t
	}
	return &g, nil
}

func (s *GroupStore) GetByName(ctx context.Context, name string) (*models.Group, error) {
	nameUpcase := strings.ToUpper(name)
	var g models.Group
	var createdAt, updatedAt sql.NullInt64
	query := `SELECT id, name, description, image_uri, email, is_active, created_at, updated_at FROM "groups" WHERE name_upcase = ?`
	err := s.db.QueryRowContext(ctx, core.Rebind("sqlite", query), nameUpcase).Scan(&g.ID, &g.Name, &g.Description, &g.ImageURI, &g.Email, &g.IsActive, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if createdAt.Valid {
		g.CreatedAt = time.Unix(createdAt.Int64, 0).UTC()
	}
	if updatedAt.Valid {
		t := time.Unix(updatedAt.Int64, 0).UTC()
		g.UpdatedAt = &t
	}
	return &g, nil
}

func (s *GroupStore) Create(ctx context.Context, group *models.Group) error {
	query := `INSERT INTO "groups" (id, name, name_upcase, description, image_uri, email, email_upcase, is_active, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	nameUpcase := strings.ToUpper(group.Name)
	group.NameUpcase = &nameUpcase
	var emailUpcase *string
	if group.Email != nil {
		up := strings.ToUpper(*group.Email)
		emailUpcase = &up
		group.EmailUpcase = emailUpcase
	}

	createdAtUnix := group.CreatedAt.Unix()
	var updatedAtUnix *int64
	if group.UpdatedAt != nil {
		u := group.UpdatedAt.Unix()
		updatedAtUnix = &u
	}

	_, err := s.db.ExecContext(ctx, core.Rebind("sqlite3", query), group.ID, group.Name, group.NameUpcase, group.Description, group.ImageURI, group.Email, group.EmailUpcase, group.IsActive, createdAtUnix, updatedAtUnix)
	return err
}

func (s *GroupStore) Update(ctx context.Context, group *models.Group) error {
	query := `UPDATE "groups" SET name = ?, name_upcase = ?, description = ?, image_uri = ?, email = ?, email_upcase = ?, is_active = ?, updated_at = ? WHERE id = ?`

	nameUpcase := strings.ToUpper(group.Name)
	group.NameUpcase = &nameUpcase
	var emailUpcase *string
	if group.Email != nil {
		up := strings.ToUpper(*group.Email)
		emailUpcase = &up
		group.EmailUpcase = emailUpcase
	} else {
		group.EmailUpcase = nil
	}

	var updatedAt *time.Time
	var updatedAtUnix *int64
	if group.UpdatedAt != nil {
		updatedAt = group.UpdatedAt
		u := updatedAt.Unix()
		updatedAtUnix = &u
	} else {
		now := time.Now()
		updatedAt = &now
		u := now.Unix()
		updatedAtUnix = &u
	}

	_, err := s.db.ExecContext(ctx, core.Rebind("sqlite3", query), group.Name, group.NameUpcase, group.Description, group.ImageURI, group.Email, group.EmailUpcase, group.IsActive, updatedAtUnix, group.ID)
	return err
}

func (s *GroupStore) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, core.Rebind("sqlite", `DELETE FROM "groups" WHERE id = ?`), id)
	return err
}

func (s *GroupStore) AddUser(ctx context.Context, groupID, userID uuid.UUID) error {
	query := `INSERT INTO groups_users (group_id, user_id) VALUES (?, ?)`
	_, err := s.db.ExecContext(ctx, core.Rebind("sqlite", query), groupID.String(), userID.String())
	return err
}

func (s *GroupStore) RemoveUser(ctx context.Context, groupID, userID uuid.UUID) error {
	query := `DELETE FROM groups_users WHERE group_id = ? AND user_id = ?`
	_, err := s.db.ExecContext(ctx, core.Rebind("sqlite", query), groupID.String(), userID.String())
	return err
}

func (s *GroupStore) ListUsers(ctx context.Context, groupID uuid.UUID) ([]models.GroupUser, error) {
	query := `SELECT group_id, user_id FROM groups_users WHERE group_id = ?`
	rows, err := s.db.QueryContext(ctx, core.Rebind("sqlite", query), groupID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.GroupUser
	for rows.Next() {
		var gu models.GroupUser
		if err := rows.Scan(&gu.GroupID, &gu.UserID); err != nil {
			return nil, err
		}
		items = append(items, gu)
	}
	return items, nil
}

func (s *GroupStore) ListGroupsByUserID(ctx context.Context, userID uuid.UUID) ([]models.GroupUser, error) {
	query := `SELECT group_id, user_id FROM groups_users WHERE user_id = ?`
	rows, err := s.db.QueryContext(ctx, core.Rebind("sqlite", query), userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.GroupUser
	for rows.Next() {
		var gu models.GroupUser
		if err := rows.Scan(&gu.GroupID, &gu.UserID); err != nil {
			return nil, err
		}
		items = append(items, gu)
	}
	return items, nil
}

func (s *GroupStore) AddAdmin(ctx context.Context, groupID, userID uuid.UUID) error {
	query := `INSERT INTO groups_admins (group_id, user_id) VALUES (?, ?)`
	_, err := s.db.ExecContext(ctx, core.Rebind("sqlite", query), groupID.String(), userID.String())
	return err
}

func (s *GroupStore) RemoveAdmin(ctx context.Context, groupID, userID uuid.UUID) error {
	query := `DELETE FROM groups_admins WHERE group_id = ? AND user_id = ?`
	_, err := s.db.ExecContext(ctx, core.Rebind("sqlite", query), groupID.String(), userID.String())
	return err
}

func (s *GroupStore) ListAdmins(ctx context.Context, groupID uuid.UUID) ([]models.GroupAdmin, error) {
	query := `SELECT group_id, user_id FROM groups_admins WHERE group_id = ?`
	rows, err := s.db.QueryContext(ctx, core.Rebind("sqlite", query), groupID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.GroupAdmin
	for rows.Next() {
		var ga models.GroupAdmin
		if err := rows.Scan(&ga.GroupID, &ga.UserID); err != nil {
			return nil, err
		}
		items = append(items, ga)
	}
	return items, nil
}

func (s *GroupStore) AddRole(ctx context.Context, groupID, roleID uuid.UUID) error {
	query := `INSERT INTO groups_roles (group_id, role_id) VALUES (?, ?)`
	_, err := s.db.ExecContext(ctx, core.Rebind("sqlite", query), groupID.String(), roleID.String())
	return err
}

func (s *GroupStore) RemoveRole(ctx context.Context, groupID, roleID uuid.UUID) error {
	query := `DELETE FROM groups_roles WHERE group_id = ? AND role_id = ?`
	_, err := s.db.ExecContext(ctx, core.Rebind("sqlite", query), groupID.String(), roleID.String())
	return err
}

func (s *GroupStore) ListRoles(ctx context.Context, groupID uuid.UUID) ([]models.GroupRole, error) {
	query := `SELECT group_id, role_id FROM groups_roles WHERE group_id = ?`
	rows, err := s.db.QueryContext(ctx, core.Rebind("sqlite", query), groupID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.GroupRole
	for rows.Next() {
		var gr models.GroupRole
		if err := rows.Scan(&gr.GroupID, &gr.RoleID); err != nil {
			return nil, err
		}
		items = append(items, gr)
	}
	return items, nil
}
