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

type EnvironmentStore struct {
	db *sql.DB
}

func NewEnvironmentStore(db *sql.DB) *EnvironmentStore {
	return &EnvironmentStore{db: db}
}

func (s *EnvironmentStore) List(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.Environment], error) {
	fieldMap := map[string]string{
		"id":          "id",
		"name":        "name_upcase",
		"description": "description",
		"createdAt":   "created_at",
		"updatedAt":   "updated_at",
	}

	baseQuery := "SELECT id, project_id, name, description, created_at, updated_at FROM environments WHERE project_id = ?"
	query, listArgs := core.BuildListQuery("sqlite", baseQuery, opts, fieldMap)
	finalArgs := append([]interface{}{projectID}, listArgs...)

	rows, err := s.db.QueryContext(ctx, core.Rebind("sqlite", query), finalArgs...)
	if err != nil {
		return core.ListResult[models.Environment]{}, err
	}
	defer rows.Close()

	var items []models.Environment
	for rows.Next() {
		var e models.Environment
		var idStr, pIdStr string

		var createdAt, updatedAt sql.NullInt64
		if err := rows.Scan(&idStr, &pIdStr, &e.Name, &e.Description, &createdAt, &updatedAt); err != nil {
			return core.ListResult[models.Environment]{}, err
		}
		if createdAt.Valid {
			e.CreatedAt = time.Unix(createdAt.Int64, 0)
		}
		if updatedAt.Valid {
			t := time.Unix(updatedAt.Int64, 0)
			e.UpdatedAt = &t
		}

		e.ID, _ = uuid.Parse(idStr)
		e.ProjectID, _ = uuid.Parse(pIdStr)
		items = append(items, e)
	}

	var total int64
	countQuery := "SELECT COUNT(*) FROM environments WHERE project_id = ?"
	countQuery, countArgs := core.BuildCountQuery(countQuery, core.ListOptions{Filter: opts.Filter}, fieldMap)
	finalCountArgs := append([]interface{}{projectID}, countArgs...)
	err = s.db.QueryRowContext(ctx, core.Rebind("sqlite", countQuery), finalCountArgs...).Scan(&total)
	if err != nil {
		return core.ListResult[models.Environment]{}, err
	}

	return core.ListResult[models.Environment]{Items: items, TotalCount: int(total)}, nil
}

func (s *EnvironmentStore) Get(ctx context.Context, id uuid.UUID) (*models.Environment, error) {
	query := "SELECT id, project_id, name, description, created_at, updated_at FROM environments WHERE id = ?"
	var e models.Environment
	var idStr, pIdStr string

	var createdAt, updatedAt sql.NullInt64
	err := s.db.QueryRowContext(ctx, core.Rebind("sqlite", query), id).Scan(&idStr, &pIdStr, &e.Name, &e.Description, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if createdAt.Valid {
		e.CreatedAt = time.Unix(createdAt.Int64, 0)
	}
	if updatedAt.Valid {
		t := time.Unix(updatedAt.Int64, 0)
		e.UpdatedAt = &t
	}

	e.ID, _ = uuid.Parse(idStr)
	e.ProjectID, _ = uuid.Parse(pIdStr)
	return &e, nil
}

func (s *EnvironmentStore) GetByName(ctx context.Context, projectID uuid.UUID, name string) (*models.Environment, error) {
	nameUpcase := strings.ToUpper(name)
	query := "SELECT id, project_id, name, description, created_at, updated_at FROM environments WHERE project_id = ? AND name_upcase = ?"
	var e models.Environment
	var idStr, pIdStr string

	var createdAt, updatedAt sql.NullInt64
	err := s.db.QueryRowContext(ctx, core.Rebind("sqlite", query), projectID, nameUpcase).Scan(&idStr, &pIdStr, &e.Name, &e.Description, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if createdAt.Valid {
		e.CreatedAt = time.Unix(createdAt.Int64, 0)
	}
	if updatedAt.Valid {
		t := time.Unix(updatedAt.Int64, 0)
		e.UpdatedAt = &t
	}

	e.ID, _ = uuid.Parse(idStr)
	e.ProjectID, _ = uuid.Parse(pIdStr)
	return &e, nil
}

func (s *EnvironmentStore) Create(ctx context.Context, env *models.Environment) error {
	nameUpcase := strings.ToUpper(env.Name)
	env.NameUpcase = &nameUpcase
	query := "INSERT INTO environments (id, project_id, name, name_upcase, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"

	var updatedAt *int64
	if env.UpdatedAt != nil {
		t := env.UpdatedAt.Unix()
		updatedAt = &t
	}
	_, err := s.db.ExecContext(ctx, core.Rebind("sqlite", query), env.ID, env.ProjectID, env.Name, env.NameUpcase, env.Description, env.CreatedAt.Unix(), updatedAt)

	return err
}

func (s *EnvironmentStore) Update(ctx context.Context, env *models.Environment) error {
	nameUpcase := strings.ToUpper(env.Name)
	env.NameUpcase = &nameUpcase

	query := "UPDATE environments SET name = ?, name_upcase = ?, description = ?, updated_at = ? WHERE id = ?"

	var updatedAt int64
	if env.UpdatedAt != nil {
		updatedAt = env.UpdatedAt.Unix()
	} else {
		updatedAt = time.Now().Unix()
	}
	_, err := s.db.ExecContext(ctx, core.Rebind("sqlite", query), env.Name, env.NameUpcase, env.Description, updatedAt, env.ID)

	return err
}

func (s *EnvironmentStore) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, core.Rebind("sqlite", "DELETE FROM environments WHERE id = ?"), id)
	return err
}
