package pg

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

type ProjectStore struct {
	db *sql.DB
}

func NewProjectStore(db *sql.DB) *ProjectStore {
	return &ProjectStore{db: db}
}

func (s *ProjectStore) List(ctx context.Context, opts core.ListOptions) (core.ListResult[models.Project], error) {
	fieldMap := map[string]string{
		"id":        "id",
		"name":      "name_upcase",
		"slug":      "slug",
		"isActive":  "is_active",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}

	baseQuery := `SELECT id, name, slug, description, is_active, created_at, updated_at FROM projects`
	query, args := core.BuildListQuery("postgres", baseQuery, opts, fieldMap)

	rows, err := s.db.QueryContext(ctx, core.Rebind("postgres", query), args...)
	if err != nil {
		return core.ListResult[models.Project]{}, err
	}
	defer rows.Close()

	var items []models.Project
	for rows.Next() {
		var p models.Project
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.IsActive, &createdAt, &updatedAt); err != nil {
			return core.ListResult[models.Project]{}, err
		}
		if createdAt.Valid {
			p.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			p.UpdatedAt = &t
		}
		items = append(items, p)
	}

	var total int64
	countQuery := `SELECT COUNT(*) FROM projects`
	countQuery, countArgs := core.BuildCountQuery(countQuery, core.ListOptions{Filter: opts.Filter}, fieldMap)
	err = s.db.QueryRowContext(ctx, core.Rebind("postgres", countQuery), countArgs...).Scan(&total)
	if err != nil {
		return core.ListResult[models.Project]{}, err
	}

	return core.ListResult[models.Project]{Items: items, TotalCount: int(total)}, nil
}

func (s *ProjectStore) Get(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	var p models.Project
	var createdAt, updatedAt sql.NullTime
	query := `SELECT id, name, slug, description, is_active, created_at, updated_at FROM projects WHERE id = ?`
	err := s.db.QueryRowContext(ctx, core.Rebind("postgres", query), id).Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.IsActive, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if createdAt.Valid {
		p.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		p.UpdatedAt = &t
	}
	return &p, nil
}

func (s *ProjectStore) GetBySlug(ctx context.Context, slug string) (*models.Project, error) {
	var p models.Project
	var createdAt, updatedAt sql.NullTime
	query := `SELECT id, name, slug, description, is_active, created_at, updated_at FROM projects WHERE slug = ?`
	err := s.db.QueryRowContext(ctx, core.Rebind("postgres", query), slug).Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.IsActive, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	if createdAt.Valid {
		p.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		p.UpdatedAt = &t
	}
	return &p, nil
}

func (s *ProjectStore) Create(ctx context.Context, project *models.Project) error {
	query := `INSERT INTO projects (id, name, name_upcase, slug, description, is_active, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	nameUpcase := strings.ToUpper(project.Name)
	project.NameUpcase = &nameUpcase

	_, err := s.db.ExecContext(ctx, core.Rebind("postgres", query), project.ID, project.Name, project.NameUpcase, project.Slug, project.Description, project.IsActive, project.CreatedAt, project.UpdatedAt)
	return err
}

func (s *ProjectStore) Update(ctx context.Context, project *models.Project) error {
	query := `UPDATE projects SET name = ?, name_upcase = ?, slug = ?, description = ?, is_active = ?, updated_at = ? WHERE id = ?`

	nameUpcase := strings.ToUpper(project.Name)
	project.NameUpcase = &nameUpcase

	var updatedAt *time.Time
	if project.UpdatedAt != nil {
		updatedAt = project.UpdatedAt
	} else {
		now := time.Now()
		updatedAt = &now
	}

	_, err := s.db.ExecContext(ctx, core.Rebind("postgres", query), project.Name, project.NameUpcase, project.Slug, project.Description, project.IsActive, updatedAt, project.ID)
	return err
}

func (s *ProjectStore) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, core.Rebind("postgres", `DELETE FROM projects WHERE id = ?`), id)
	return err
}

func (s *ProjectStore) AddGroup(ctx context.Context, projectID, groupID uuid.UUID, permissions int64) error {
	query := `INSERT INTO projects_groups (project_id, group_id, permissions) VALUES (?, ?, ?)`
	_, err := s.db.ExecContext(ctx, core.Rebind("postgres", query), projectID, groupID, permissions)
	return err
}

func (s *ProjectStore) RemoveGroup(ctx context.Context, projectID, groupID uuid.UUID) error {
	query := `DELETE FROM projects_groups WHERE project_id = ? AND group_id = ?`
	_, err := s.db.ExecContext(ctx, core.Rebind("postgres", query), projectID, groupID)
	return err
}

func (s *ProjectStore) UpdateGroupPermissions(ctx context.Context, projectID, groupID uuid.UUID, permissions int64) error {
	query := `UPDATE projects_groups SET permissions = ? WHERE project_id = ? AND group_id = ?`
	_, err := s.db.ExecContext(ctx, core.Rebind("postgres", query), permissions, projectID, groupID)
	return err
}

func (s *ProjectStore) ListGroups(ctx context.Context, projectID uuid.UUID) ([]models.ProjectGroup, error) {
	query := `SELECT project_id, group_id, permissions FROM projects_groups WHERE project_id = ?`
	rows, err := s.db.QueryContext(ctx, core.Rebind("postgres", query), projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.ProjectGroup
	for rows.Next() {
		var pg models.ProjectGroup
		if err := rows.Scan(&pg.ProjectID, &pg.GroupID, &pg.Permissions); err != nil {
			return nil, err
		}
		groups = append(groups, pg)
	}
	return groups, rows.Err()
}
