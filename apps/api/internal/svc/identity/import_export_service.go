package identity

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/google/uuid"
)

type ImportExportService interface {
	ExportUsers(ctx context.Context, format string) ([]byte, error)
	ImportUsers(ctx context.Context, format string, data []byte) error

	ExportRoles(ctx context.Context, format string) ([]byte, error)
	ImportRoles(ctx context.Context, format string, data []byte) error
}

type importExportService struct {
	userStore stores.UserStore
	roleStore stores.RoleStore
}

func NewImportExportService(userStore stores.UserStore, roleStore stores.RoleStore) ImportExportService {
	return &importExportService{
		userStore: userStore,
		roleStore: roleStore,
	}
}

func (s *importExportService) userCol() core.CollectionDef[models.User] {
	return core.CollectionDef[models.User]{
		NameStr:     "users",
		HeaderNames: []string{"ID", "PrimaryEmail", "PrimaryPhone", "Name", "IsBanned", "CreatedAt"},
		ToRow: func(u models.User) []string {
			email := ""
			if u.PrimaryEmail != nil {
				email = *u.PrimaryEmail
			}
			phone := ""
			if u.PrimaryPhone != nil {
				phone = *u.PrimaryPhone
			}
			name := ""
			if u.Name != nil {
				name = *u.Name
			}
			return []string{
				u.ID.String(), email, phone, name, strconv.FormatBool(u.IsBanned), u.CreatedAt.Format(time.RFC3339),
			}
		},
		FromRow: func(row []string) (models.User, error) {
			if len(row) < 6 {
				return models.User{}, fmt.Errorf("row invalid format")
			}
			uid, err := uuid.Parse(row[0])
			if err != nil {
				uid = uuid.New()
			}
			var email *string
			if row[1] != "" {
				email = &row[1]
			}
			var phone *string
			if row[2] != "" {
				phone = &row[2]
			}
			var name *string
			if row[3] != "" {
				name = &row[3]
			}
			isBanned, _ := strconv.ParseBool(row[4])
			createdAt, err := time.Parse(time.RFC3339, row[5])
			if err != nil {
				createdAt = time.Now()
			}

			return models.User{
				ID:           uid,
				PrimaryEmail: email,
				PrimaryPhone: phone,
				Name:         name,
				IsBanned:     isBanned,
				CreatedAt:    createdAt,
			}, nil
		},
		GetItems: func(ctx context.Context, opts core.ListOptions) ([]models.User, error) {
			res, err := s.userStore.List(ctx, opts)
			return res.Items, err
		},
		SaveItems: func(ctx context.Context, items []models.User) ([]core.ImportError, error) {
			var errs []core.ImportError
			for i, u := range items {
				if u.ID == uuid.Nil {
					u.ID = uuid.New()
				}
				existing, _ := s.userStore.Get(ctx, u.ID)
				if existing != nil {
					if err := s.userStore.Update(ctx, &u); err != nil {
						errs = append(errs, core.ImportError{RowNumber: i + 2, Error: err})
					}
				} else {
					if err := s.userStore.Create(ctx, &u); err != nil {
						errs = append(errs, core.ImportError{RowNumber: i + 2, Error: err})
					}
				}
			}
			return errs, nil
		},
	}
}

func (s *importExportService) roleCol() core.CollectionDef[models.Role] {
	return core.CollectionDef[models.Role]{
		NameStr:     "roles",
		HeaderNames: []string{"ID", "Name", "Description"},
		ToRow: func(r models.Role) []string {
			return []string{
				r.ID.String(), r.Name, r.Description,
			}
		},
		FromRow: func(row []string) (models.Role, error) {
			if len(row) < 3 {
				return models.Role{}, fmt.Errorf("row invalid format")
			}
			uid, err := uuid.Parse(row[0])
			if err != nil {
				uid = uuid.New()
			}
			return models.Role{
				ID:          uid,
				Name:        row[1],
				Description: row[2],
			}, nil
		},
		GetItems: func(ctx context.Context, opts core.ListOptions) ([]models.Role, error) {
			res, err := s.roleStore.List(ctx, opts)
			return res.Items, err
		},
		SaveItems: func(ctx context.Context, items []models.Role) ([]core.ImportError, error) {
			var errs []core.ImportError
			for i, r := range items {
				if r.ID == uuid.Nil {
					r.ID = uuid.New()
				}
				existing, _ := s.roleStore.Get(ctx, r.ID)
				if existing != nil {
					if err := s.roleStore.Update(ctx, &r); err != nil {
						errs = append(errs, core.ImportError{RowNumber: i + 2, Error: err})
					}
				} else {
					if err := s.roleStore.Create(ctx, &r); err != nil {
						errs = append(errs, core.ImportError{RowNumber: i + 2, Error: err})
					}
				}
			}
			return errs, nil
		},
	}
}

func (s *importExportService) ExportUsers(ctx context.Context, format string) ([]byte, error) {
	opts := core.ListOptions{Limit: 10000}
	var buf bytes.Buffer
	err := core.BatchExport(ctx, &buf, core.ExportFormat(format), opts, s.userCol())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *importExportService) ImportUsers(ctx context.Context, format string, data []byte) error {
	errs, err := core.BatchImport(ctx, bytes.NewReader(data), core.ExportFormat(format), s.userCol())
	if err != nil {
		return err
	}
	if len(errs) > 0 {
		return fmt.Errorf("import finished with %d errors. first error: %v", len(errs), errs[0].Error)
	}
	return nil
}

func (s *importExportService) ExportRoles(ctx context.Context, format string) ([]byte, error) {
	opts := core.ListOptions{Limit: 10000}
	var buf bytes.Buffer
	err := core.BatchExport(ctx, &buf, core.ExportFormat(format), opts, s.roleCol())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *importExportService) ImportRoles(ctx context.Context, format string, data []byte) error {
	errs, err := core.BatchImport(ctx, bytes.NewReader(data), core.ExportFormat(format), s.roleCol())
	if err != nil {
		return err
	}
	if len(errs) > 0 {
		return fmt.Errorf("import finished with %d errors. first error: %v", len(errs), errs[0].Error)
	}
	return nil
}
