package identity

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/frostyeti/hyprship/apps/api/internal/stores"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
)

var (
	ErrUnsupportedFormat = errors.New("unsupported import/export format")
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

// User Export
func (s *importExportService) ExportUsers(ctx context.Context, format string) ([]byte, error) {
	opts := core.ListOptions{Limit: 10000}
	res, err := s.userStore.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return s.serializeUsers(res.Items, format)
}

func (s *importExportService) serializeUsers(users []models.User, format string) ([]byte, error) {
	switch format {
	case "json":
		return json.MarshalIndent(map[string]any{"users": users}, "", "  ")
	case "yaml", "yml":
		return yaml.Marshal(map[string]any{"users": users})
	case "csv":
		var buf bytes.Buffer
		w := csv.NewWriter(&buf)
		w.Write([]string{"ID", "PrimaryEmail", "PrimaryPhone", "Name", "IsBanned", "CreatedAt"})
		for _, u := range users {
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
			w.Write([]string{
				u.ID.String(), email, phone, name, strconv.FormatBool(u.IsBanned), u.CreatedAt.Format(time.RFC3339),
			})
		}
		w.Flush()
		return buf.Bytes(), w.Error()
	case "excel", "xlsx":
		f := excelize.NewFile()
		sheet := "Users"
		f.SetSheetName("Sheet1", sheet)
		f.SetSheetRow(sheet, "A1", &[]interface{}{"ID", "PrimaryEmail", "PrimaryPhone", "Name", "IsBanned", "CreatedAt"})
		for i, u := range users {
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
			row := i + 2
			f.SetSheetRow(sheet, fmt.Sprintf("A%d", row), &[]interface{}{
				u.ID.String(), email, phone, name, u.IsBanned, u.CreatedAt.Format(time.RFC3339),
			})
		}
		var buf bytes.Buffer
		err := f.Write(&buf)
		return buf.Bytes(), err
	default:
		return nil, ErrUnsupportedFormat
	}
}

// User Import
func (s *importExportService) ImportUsers(ctx context.Context, format string, data []byte) error {
	var users []models.User
	var err error

	switch format {
	case "json":
		var payload struct {
			Users []models.User `json:"users"`
		}
		if err = json.Unmarshal(data, &payload); err != nil {
			return err
		}
		users = payload.Users
	case "yaml", "yml":
		var payload struct {
			Users []models.User `yaml:"users"`
		}
		if err = yaml.Unmarshal(data, &payload); err != nil {
			return err
		}
		users = payload.Users
	case "csv":
		r := csv.NewReader(bytes.NewReader(data))
		records, err := r.ReadAll()
		if err != nil {
			return err
		}
		users, err = s.parseUsersCSV(records)
		if err != nil {
			return err
		}
	case "excel", "xlsx":
		f, err := excelize.OpenReader(bytes.NewReader(data))
		if err != nil {
			return err
		}
		defer f.Close()
		sheet := f.GetSheetName(0)
		rows, err := f.GetRows(sheet)
		if err != nil {
			return err
		}
		users, err = s.parseUsersCSV(rows)
		if err != nil {
			return err
		}
	default:
		return ErrUnsupportedFormat
	}

	for _, u := range users {
		if u.ID == uuid.Nil {
			u.ID = uuid.New()
		}
		existing, _ := s.userStore.Get(ctx, u.ID)
		if existing != nil {
			s.userStore.Update(ctx, &u)
		} else {
			s.userStore.Create(ctx, &u)
		}
	}

	return nil
}

func (s *importExportService) parseUsersCSV(records [][]string) ([]models.User, error) {
	if len(records) < 2 {
		return nil, errors.New("empty records")
	}
	var users []models.User
	for i, row := range records[1:] {
		if len(row) < 6 {
			return nil, fmt.Errorf("row %d invalid format", i+2)
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

		users = append(users, models.User{
			ID:           uid,
			PrimaryEmail: email,
			PrimaryPhone: phone,
			Name:         name,
			IsBanned:     isBanned,
			CreatedAt:    createdAt,
		})
	}
	return users, nil
}

// Role Export
func (s *importExportService) ExportRoles(ctx context.Context, format string) ([]byte, error) {
	opts := core.ListOptions{Limit: 10000}
	res, err := s.roleStore.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return s.serializeRoles(res.Items, format)
}

func (s *importExportService) serializeRoles(roles []models.Role, format string) ([]byte, error) {
	switch format {
	case "json":
		return json.MarshalIndent(map[string]any{"roles": roles}, "", "  ")
	case "yaml", "yml":
		return yaml.Marshal(map[string]any{"roles": roles})
	case "csv":
		var buf bytes.Buffer
		w := csv.NewWriter(&buf)
		w.Write([]string{"ID", "Name", "Description"})
		for _, r := range roles {
			w.Write([]string{
				r.ID.String(), r.Name, r.Description,
			})
		}
		w.Flush()
		return buf.Bytes(), w.Error()
	case "excel", "xlsx":
		f := excelize.NewFile()
		sheet := "Roles"
		f.SetSheetName("Sheet1", sheet)
		f.SetSheetRow(sheet, "A1", &[]interface{}{"ID", "Name", "Description"})
		for i, r := range roles {
			row := i + 2
			f.SetSheetRow(sheet, fmt.Sprintf("A%d", row), &[]interface{}{
				r.ID.String(), r.Name, r.Description,
			})
		}
		var buf bytes.Buffer
		err := f.Write(&buf)
		return buf.Bytes(), err
	default:
		return nil, ErrUnsupportedFormat
	}
}

// Role Import
func (s *importExportService) ImportRoles(ctx context.Context, format string, data []byte) error {
	var roles []models.Role
	var err error

	switch format {
	case "json":
		var payload struct {
			Roles []models.Role `json:"roles"`
		}
		if err = json.Unmarshal(data, &payload); err != nil {
			return err
		}
		roles = payload.Roles
	case "yaml", "yml":
		var payload struct {
			Roles []models.Role `yaml:"roles"`
		}
		if err = yaml.Unmarshal(data, &payload); err != nil {
			return err
		}
		roles = payload.Roles
	case "csv":
		r := csv.NewReader(bytes.NewReader(data))
		records, err := r.ReadAll()
		if err != nil {
			return err
		}
		roles, err = s.parseRolesCSV(records)
		if err != nil {
			return err
		}
	case "excel", "xlsx":
		f, err := excelize.OpenReader(bytes.NewReader(data))
		if err != nil {
			return err
		}
		defer f.Close()
		sheet := f.GetSheetName(0)
		rows, err := f.GetRows(sheet)
		if err != nil {
			return err
		}
		roles, err = s.parseRolesCSV(rows)
		if err != nil {
			return err
		}
	default:
		return ErrUnsupportedFormat
	}

	for _, r := range roles {
		if r.ID == uuid.Nil {
			r.ID = uuid.New()
		}
		existing, _ := s.roleStore.Get(ctx, r.ID)
		if existing != nil {
			s.roleStore.Update(ctx, &r)
		} else {
			s.roleStore.Create(ctx, &r)
		}
	}

	return nil
}

func (s *importExportService) parseRolesCSV(records [][]string) ([]models.Role, error) {
	if len(records) < 2 {
		return nil, errors.New("empty records")
	}
	var roles []models.Role
	for i, row := range records[1:] {
		if len(row) < 3 {
			return nil, fmt.Errorf("row %d invalid format", i+2)
		}
		uid, err := uuid.Parse(row[0])
		if err != nil {
			uid = uuid.New()
		}

		roles = append(roles, models.Role{
			ID:          uid,
			Name:        row[1],
			Description: row[2],
		})
	}
	return roles, nil
}
