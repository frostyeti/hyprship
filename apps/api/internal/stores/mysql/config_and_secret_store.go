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

// ======== CONFIG STORE ========
type ConfigStore struct {
	db *sql.DB
}

func NewConfigStore(db *sql.DB) *ConfigStore {
	return &ConfigStore{db: db}
}

// Config Files

func (s *ConfigStore) ListConfigFileNames(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.ConfigFileName], error) {
	fieldMap := map[string]string{
		"id":        "id",
		"name":      "name_upcase",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}
	baseQuery := "SELECT CAST(id AS CHAR(36)), CAST(project_id AS CHAR(36)), name, created_at, updated_at FROM config_file_names WHERE project_id = ?"
	query, listArgs := core.BuildListQuery("mysql", baseQuery, opts, fieldMap)
	finalArgs := append([]interface{}{projectID.String()}, listArgs...)

	rows, err := s.db.QueryContext(ctx, core.Rebind("mysql", query), finalArgs...)
	if err != nil {
		return core.ListResult[models.ConfigFileName]{}, err
	}
	defer rows.Close()

	var items []models.ConfigFileName
	for rows.Next() {
		var n models.ConfigFileName
		var idStr, pIdStr string
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&idStr, &pIdStr, &n.Name, &createdAt, &updatedAt); err != nil {
			return core.ListResult[models.ConfigFileName]{}, err
		}
		n.ID, _ = uuid.Parse(idStr)
		n.ProjectID, _ = uuid.Parse(pIdStr)
		if createdAt.Valid {
			n.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			n.UpdatedAt = &t
		}
		items = append(items, n)
	}

	var total int64
	countQuery, countArgs := core.BuildCountQuery("SELECT COUNT(*) FROM config_file_names WHERE project_id = ?", core.ListOptions{Filter: opts.Filter}, fieldMap)
	finalCountArgs := append([]interface{}{projectID.String()}, countArgs...)
	err = s.db.QueryRowContext(ctx, core.Rebind("mysql", countQuery), finalCountArgs...).Scan(&total)
	return core.ListResult[models.ConfigFileName]{Items: items, TotalCount: int(total)}, err
}

func (s *ConfigStore) GetConfigFileName(ctx context.Context, id uuid.UUID) (*models.ConfigFileName, error) {
	query := "SELECT CAST(id AS CHAR(36)), CAST(project_id AS CHAR(36)), name, created_at, updated_at FROM config_file_names WHERE id = ?"
	var n models.ConfigFileName
	var idStr, pIdStr string
	var createdAt, updatedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, core.Rebind("mysql", query), id.String()).Scan(&idStr, &pIdStr, &n.Name, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	n.ID, _ = uuid.Parse(idStr)
	n.ProjectID, _ = uuid.Parse(pIdStr)
	if createdAt.Valid {
		n.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		n.UpdatedAt = &t
	}
	return &n, nil
}

func (s *ConfigStore) GetConfigFileNameByName(ctx context.Context, projectID uuid.UUID, name string) (*models.ConfigFileName, error) {
	nameUpcase := strings.ToUpper(name)
	query := "SELECT CAST(id AS CHAR(36)), CAST(project_id AS CHAR(36)), name, created_at, updated_at FROM config_file_names WHERE project_id = ? AND name_upcase = ?"
	var n models.ConfigFileName
	var idStr, pIdStr string
	var createdAt, updatedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, core.Rebind("mysql", query), projectID.String(), nameUpcase).Scan(&idStr, &pIdStr, &n.Name, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	n.ID, _ = uuid.Parse(idStr)
	n.ProjectID, _ = uuid.Parse(pIdStr)
	if createdAt.Valid {
		n.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		n.UpdatedAt = &t
	}
	return &n, nil
}

func (s *ConfigStore) CreateConfigFileName(ctx context.Context, name *models.ConfigFileName) error {
	nameUpcase := strings.ToUpper(name.Name)
	name.NameUpcase = &nameUpcase
	query := "INSERT INTO config_file_names (id, project_id, name, name_upcase, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := s.db.ExecContext(ctx, core.Rebind("mysql", query), name.ID.String(), name.ProjectID.String(), name.Name, name.NameUpcase, name.CreatedAt, name.UpdatedAt)
	return err
}

func (s *ConfigStore) DeleteConfigFileName(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, core.Rebind("mysql", "DELETE FROM config_file_names WHERE id = ?"), id.String())
	return err
}

func (s *ConfigStore) GetConfigFile(ctx context.Context, nameID, envID uuid.UUID) (*models.ConfigFile, error) {
	query := "SELECT CAST(id AS CHAR(36)), CAST(config_file_name_id AS CHAR(36)), CAST(environment_id AS CHAR(36)), content, created_at, updated_at FROM config_files WHERE config_file_name_id = ? AND environment_id = ?"
	var f models.ConfigFile
	var idStr, nIdStr, eIdStr string
	var createdAt, updatedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, core.Rebind("mysql", query), nameID.String(), envID.String()).Scan(&idStr, &nIdStr, &eIdStr, &f.Content, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	f.ID, _ = uuid.Parse(idStr)
	f.ConfigFileNameID, _ = uuid.Parse(nIdStr)
	f.EnvironmentID, _ = uuid.Parse(eIdStr)
	if createdAt.Valid {
		f.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		f.UpdatedAt = &t
	}
	return &f, nil
}

func (s *ConfigStore) UpsertConfigFile(ctx context.Context, file *models.ConfigFile) error {
	existing, err := s.GetConfigFile(ctx, file.ConfigFileNameID, file.EnvironmentID)
	if err != nil {
		return err
	}
	if existing == nil {
		query := "INSERT INTO config_files (id, config_file_name_id, environment_id, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
		_, err = s.db.ExecContext(ctx, core.Rebind("mysql", query), file.ID.String(), file.ConfigFileNameID.String(), file.EnvironmentID.String(), file.Content, file.CreatedAt, file.UpdatedAt)
		return err
	}

	now := time.Now()
	query := "UPDATE config_files SET content = ?, updated_at = ? WHERE id = ?"
	_, err = s.db.ExecContext(ctx, core.Rebind("mysql", query), file.Content, now, existing.ID.String())
	return err
}

// Env Variables

func (s *ConfigStore) ListEnvVariableNames(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.EnvVariableName], error) {
	fieldMap := map[string]string{
		"id":        "id",
		"name":      "name_upcase",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}
	baseQuery := "SELECT CAST(id AS CHAR(36)), CAST(project_id AS CHAR(36)), name, created_at, updated_at FROM env_variable_names WHERE project_id = ?"
	query, listArgs := core.BuildListQuery("mysql", baseQuery, opts, fieldMap)
	finalArgs := append([]interface{}{projectID.String()}, listArgs...)

	rows, err := s.db.QueryContext(ctx, core.Rebind("mysql", query), finalArgs...)
	if err != nil {
		return core.ListResult[models.EnvVariableName]{}, err
	}
	defer rows.Close()

	var items []models.EnvVariableName
	for rows.Next() {
		var n models.EnvVariableName
		var idStr, pIdStr string
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&idStr, &pIdStr, &n.Name, &createdAt, &updatedAt); err != nil {
			return core.ListResult[models.EnvVariableName]{}, err
		}
		n.ID, _ = uuid.Parse(idStr)
		n.ProjectID, _ = uuid.Parse(pIdStr)
		if createdAt.Valid {
			n.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			n.UpdatedAt = &t
		}
		items = append(items, n)
	}

	var total int64
	countQuery, countArgs := core.BuildCountQuery("SELECT COUNT(*) FROM env_variable_names WHERE project_id = ?", core.ListOptions{Filter: opts.Filter}, fieldMap)
	finalCountArgs := append([]interface{}{projectID.String()}, countArgs...)
	err = s.db.QueryRowContext(ctx, core.Rebind("mysql", countQuery), finalCountArgs...).Scan(&total)
	return core.ListResult[models.EnvVariableName]{Items: items, TotalCount: int(total)}, err
}

func (s *ConfigStore) GetEnvVariableName(ctx context.Context, id uuid.UUID) (*models.EnvVariableName, error) {
	query := "SELECT CAST(id AS CHAR(36)), CAST(project_id AS CHAR(36)), name, created_at, updated_at FROM env_variable_names WHERE id = ?"
	var n models.EnvVariableName
	var idStr, pIdStr string
	var createdAt, updatedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, core.Rebind("mysql", query), id.String()).Scan(&idStr, &pIdStr, &n.Name, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	n.ID, _ = uuid.Parse(idStr)
	n.ProjectID, _ = uuid.Parse(pIdStr)
	if createdAt.Valid {
		n.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		n.UpdatedAt = &t
	}
	return &n, nil
}

func (s *ConfigStore) GetEnvVariableNameByName(ctx context.Context, projectID uuid.UUID, name string) (*models.EnvVariableName, error) {
	nameUpcase := strings.ToUpper(name)
	query := "SELECT CAST(id AS CHAR(36)), CAST(project_id AS CHAR(36)), name, created_at, updated_at FROM env_variable_names WHERE project_id = ? AND name_upcase = ?"
	var n models.EnvVariableName
	var idStr, pIdStr string
	var createdAt, updatedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, core.Rebind("mysql", query), projectID.String(), nameUpcase).Scan(&idStr, &pIdStr, &n.Name, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	n.ID, _ = uuid.Parse(idStr)
	n.ProjectID, _ = uuid.Parse(pIdStr)
	if createdAt.Valid {
		n.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		n.UpdatedAt = &t
	}
	return &n, nil
}

func (s *ConfigStore) CreateEnvVariableName(ctx context.Context, name *models.EnvVariableName) error {
	nameUpcase := strings.ToUpper(name.Name)
	name.NameUpcase = &nameUpcase
	query := "INSERT INTO env_variable_names (id, project_id, name, name_upcase, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := s.db.ExecContext(ctx, core.Rebind("mysql", query), name.ID.String(), name.ProjectID.String(), name.Name, name.NameUpcase, name.CreatedAt, name.UpdatedAt)
	return err
}

func (s *ConfigStore) DeleteEnvVariableName(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, core.Rebind("mysql", "DELETE FROM env_variable_names WHERE id = ?"), id.String())
	return err
}

func (s *ConfigStore) GetEnvVariable(ctx context.Context, nameID, envID uuid.UUID) (*models.EnvVariable, error) {
	query := "SELECT CAST(id AS CHAR(36)), CAST(env_variable_name_id AS CHAR(36)), CAST(environment_id AS CHAR(36)), value, created_at, updated_at FROM env_variables WHERE env_variable_name_id = ? AND environment_id = ?"
	var v models.EnvVariable
	var idStr, nIdStr, eIdStr string
	var createdAt, updatedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, core.Rebind("mysql", query), nameID.String(), envID.String()).Scan(&idStr, &nIdStr, &eIdStr, &v.Value, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	v.ID, _ = uuid.Parse(idStr)
	v.EnvVariableNameID, _ = uuid.Parse(nIdStr)
	v.EnvironmentID, _ = uuid.Parse(eIdStr)
	if createdAt.Valid {
		v.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		v.UpdatedAt = &t
	}
	return &v, nil
}

func (s *ConfigStore) UpsertEnvVariable(ctx context.Context, v *models.EnvVariable) error {
	existing, err := s.GetEnvVariable(ctx, v.EnvVariableNameID, v.EnvironmentID)
	if err != nil {
		return err
	}
	if existing == nil {
		query := "INSERT INTO env_variables (id, env_variable_name_id, environment_id, value, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
		_, err = s.db.ExecContext(ctx, core.Rebind("mysql", query), v.ID.String(), v.EnvVariableNameID.String(), v.EnvironmentID.String(), v.Value, v.CreatedAt, v.UpdatedAt)
		return err
	}

	now := time.Now()
	query := "UPDATE env_variables SET value = ?, updated_at = ? WHERE id = ?"
	_, err = s.db.ExecContext(ctx, core.Rebind("mysql", query), v.Value, now, existing.ID.String())
	return err
}

// ======== SECRET STORE ========
type SecretStore struct {
	db *sql.DB
}

func NewSecretStore(db *sql.DB) *SecretStore {
	return &SecretStore{db: db}
}

// Secrets

func (s *SecretStore) ListSecretNames(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.SecretName], error) {
	fieldMap := map[string]string{
		"id":        "id",
		"name":      "name_upcase",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}
	baseQuery := "SELECT CAST(id AS CHAR(36)), CAST(project_id AS CHAR(36)), name, created_at, updated_at FROM secret_names WHERE project_id = ?"
	query, listArgs := core.BuildListQuery("mysql", baseQuery, opts, fieldMap)
	finalArgs := append([]interface{}{projectID.String()}, listArgs...)

	rows, err := s.db.QueryContext(ctx, core.Rebind("mysql", query), finalArgs...)
	if err != nil {
		return core.ListResult[models.SecretName]{}, err
	}
	defer rows.Close()

	var items []models.SecretName
	for rows.Next() {
		var n models.SecretName
		var idStr, pIdStr string
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&idStr, &pIdStr, &n.Name, &createdAt, &updatedAt); err != nil {
			return core.ListResult[models.SecretName]{}, err
		}
		n.ID, _ = uuid.Parse(idStr)
		n.ProjectID, _ = uuid.Parse(pIdStr)
		if createdAt.Valid {
			n.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			n.UpdatedAt = &t
		}
		items = append(items, n)
	}

	var total int64
	countQuery, countArgs := core.BuildCountQuery("SELECT COUNT(*) FROM secret_names WHERE project_id = ?", core.ListOptions{Filter: opts.Filter}, fieldMap)
	finalCountArgs := append([]interface{}{projectID.String()}, countArgs...)
	err = s.db.QueryRowContext(ctx, core.Rebind("mysql", countQuery), finalCountArgs...).Scan(&total)
	return core.ListResult[models.SecretName]{Items: items, TotalCount: int(total)}, err
}

func (s *SecretStore) GetSecretName(ctx context.Context, id uuid.UUID) (*models.SecretName, error) {
	query := "SELECT CAST(id AS CHAR(36)), CAST(project_id AS CHAR(36)), name, created_at, updated_at FROM secret_names WHERE id = ?"
	var n models.SecretName
	var idStr, pIdStr string
	var createdAt, updatedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, core.Rebind("mysql", query), id.String()).Scan(&idStr, &pIdStr, &n.Name, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	n.ID, _ = uuid.Parse(idStr)
	n.ProjectID, _ = uuid.Parse(pIdStr)
	if createdAt.Valid {
		n.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		n.UpdatedAt = &t
	}
	return &n, nil
}

func (s *SecretStore) GetSecretNameByName(ctx context.Context, projectID uuid.UUID, name string) (*models.SecretName, error) {
	nameUpcase := strings.ToUpper(name)
	query := "SELECT CAST(id AS CHAR(36)), CAST(project_id AS CHAR(36)), name, created_at, updated_at FROM secret_names WHERE project_id = ? AND name_upcase = ?"
	var n models.SecretName
	var idStr, pIdStr string
	var createdAt, updatedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, core.Rebind("mysql", query), projectID.String(), nameUpcase).Scan(&idStr, &pIdStr, &n.Name, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	n.ID, _ = uuid.Parse(idStr)
	n.ProjectID, _ = uuid.Parse(pIdStr)
	if createdAt.Valid {
		n.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		n.UpdatedAt = &t
	}
	return &n, nil
}

func (s *SecretStore) CreateSecretName(ctx context.Context, name *models.SecretName) error {
	nameUpcase := strings.ToUpper(name.Name)
	name.NameUpcase = &nameUpcase
	query := "INSERT INTO secret_names (id, project_id, name, name_upcase, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := s.db.ExecContext(ctx, core.Rebind("mysql", query), name.ID.String(), name.ProjectID.String(), name.Name, name.NameUpcase, name.CreatedAt, name.UpdatedAt)
	return err
}

func (s *SecretStore) DeleteSecretName(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, core.Rebind("mysql", "DELETE FROM secret_names WHERE id = ?"), id.String())
	return err
}

func (s *SecretStore) GetSecret(ctx context.Context, nameID, envID uuid.UUID) (*models.Secret, error) {
	query := "SELECT CAST(id AS CHAR(36)), CAST(secret_name_id AS CHAR(36)), CAST(environment_id AS CHAR(36)), value, created_at, updated_at FROM secrets WHERE secret_name_id = ? AND environment_id = ?"
	var sec models.Secret
	var idStr, nIdStr, eIdStr string
	var createdAt, updatedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, core.Rebind("mysql", query), nameID.String(), envID.String()).Scan(&idStr, &nIdStr, &eIdStr, &sec.Value, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	sec.ID, _ = uuid.Parse(idStr)
	sec.SecretNameID, _ = uuid.Parse(nIdStr)
	sec.EnvironmentID, _ = uuid.Parse(eIdStr)
	if createdAt.Valid {
		sec.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		sec.UpdatedAt = &t
	}
	return &sec, nil
}

func (s *SecretStore) UpsertSecret(ctx context.Context, sec *models.Secret) error {
	existing, err := s.GetSecret(ctx, sec.SecretNameID, sec.EnvironmentID)
	if err != nil {
		return err
	}
	if existing == nil {
		query := "INSERT INTO secrets (id, secret_name_id, environment_id, value, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
		_, err = s.db.ExecContext(ctx, core.Rebind("mysql", query), sec.ID.String(), sec.SecretNameID.String(), sec.EnvironmentID.String(), sec.Value, sec.CreatedAt, sec.UpdatedAt)
		return err
	}

	now := time.Now()
	query := "UPDATE secrets SET value = ?, updated_at = ? WHERE id = ?"
	_, err = s.db.ExecContext(ctx, core.Rebind("mysql", query), sec.Value, now, existing.ID.String())
	return err
}

// Certificates

func (s *SecretStore) ListCertificateNames(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.CertificateName], error) {
	fieldMap := map[string]string{
		"id":        "id",
		"name":      "name_upcase",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}
	baseQuery := "SELECT CAST(id AS CHAR(36)), CAST(project_id AS CHAR(36)), name, created_at, updated_at FROM certificate_names WHERE project_id = ?"
	query, listArgs := core.BuildListQuery("mysql", baseQuery, opts, fieldMap)
	finalArgs := append([]interface{}{projectID.String()}, listArgs...)

	rows, err := s.db.QueryContext(ctx, core.Rebind("mysql", query), finalArgs...)
	if err != nil {
		return core.ListResult[models.CertificateName]{}, err
	}
	defer rows.Close()

	var items []models.CertificateName
	for rows.Next() {
		var n models.CertificateName
		var idStr, pIdStr string
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&idStr, &pIdStr, &n.Name, &createdAt, &updatedAt); err != nil {
			return core.ListResult[models.CertificateName]{}, err
		}
		n.ID, _ = uuid.Parse(idStr)
		n.ProjectID, _ = uuid.Parse(pIdStr)
		if createdAt.Valid {
			n.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			n.UpdatedAt = &t
		}
		items = append(items, n)
	}

	var total int64
	countQuery, countArgs := core.BuildCountQuery("SELECT COUNT(*) FROM certificate_names WHERE project_id = ?", core.ListOptions{Filter: opts.Filter}, fieldMap)
	finalCountArgs := append([]interface{}{projectID.String()}, countArgs...)
	err = s.db.QueryRowContext(ctx, core.Rebind("mysql", countQuery), finalCountArgs...).Scan(&total)
	return core.ListResult[models.CertificateName]{Items: items, TotalCount: int(total)}, err
}

func (s *SecretStore) GetCertificateName(ctx context.Context, id uuid.UUID) (*models.CertificateName, error) {
	query := "SELECT CAST(id AS CHAR(36)), CAST(project_id AS CHAR(36)), name, created_at, updated_at FROM certificate_names WHERE id = ?"
	var n models.CertificateName
	var idStr, pIdStr string
	var createdAt, updatedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, core.Rebind("mysql", query), id.String()).Scan(&idStr, &pIdStr, &n.Name, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	n.ID, _ = uuid.Parse(idStr)
	n.ProjectID, _ = uuid.Parse(pIdStr)
	if createdAt.Valid {
		n.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		n.UpdatedAt = &t
	}
	return &n, nil
}

func (s *SecretStore) GetCertificateNameByName(ctx context.Context, projectID uuid.UUID, name string) (*models.CertificateName, error) {
	nameUpcase := strings.ToUpper(name)
	query := "SELECT CAST(id AS CHAR(36)), CAST(project_id AS CHAR(36)), name, created_at, updated_at FROM certificate_names WHERE project_id = ? AND name_upcase = ?"
	var n models.CertificateName
	var idStr, pIdStr string
	var createdAt, updatedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, core.Rebind("mysql", query), projectID.String(), nameUpcase).Scan(&idStr, &pIdStr, &n.Name, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	n.ID, _ = uuid.Parse(idStr)
	n.ProjectID, _ = uuid.Parse(pIdStr)
	if createdAt.Valid {
		n.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		n.UpdatedAt = &t
	}
	return &n, nil
}

func (s *SecretStore) CreateCertificateName(ctx context.Context, name *models.CertificateName) error {
	nameUpcase := strings.ToUpper(name.Name)
	name.NameUpcase = &nameUpcase
	query := "INSERT INTO certificate_names (id, project_id, name, name_upcase, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := s.db.ExecContext(ctx, core.Rebind("mysql", query), name.ID.String(), name.ProjectID.String(), name.Name, name.NameUpcase, name.CreatedAt, name.UpdatedAt)
	return err
}

func (s *SecretStore) DeleteCertificateName(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, core.Rebind("mysql", "DELETE FROM certificate_names WHERE id = ?"), id.String())
	return err
}

func (s *SecretStore) GetCertificate(ctx context.Context, nameID, envID uuid.UUID) (*models.Certificate, error) {
	query := "SELECT CAST(id AS CHAR(36)), CAST(certificate_name_id AS CHAR(36)), CAST(environment_id AS CHAR(36)), certificate_data, private_key_secret, created_at, updated_at FROM certificates WHERE certificate_name_id = ? AND environment_id = ?"
	var cert models.Certificate
	var idStr, nIdStr, eIdStr string
	var privateKeySecret sql.NullString
	var createdAt, updatedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, core.Rebind("mysql", query), nameID.String(), envID.String()).Scan(&idStr, &nIdStr, &eIdStr, &cert.CertificateData, &privateKeySecret, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	cert.ID, _ = uuid.Parse(idStr)
	cert.CertificateNameID, _ = uuid.Parse(nIdStr)
	cert.EnvironmentID, _ = uuid.Parse(eIdStr)
	if privateKeySecret.Valid {
		cert.PrivateKeySecret = &privateKeySecret.String
	}
	if createdAt.Valid {
		cert.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		cert.UpdatedAt = &t
	}
	return &cert, nil
}

func (s *SecretStore) UpsertCertificate(ctx context.Context, cert *models.Certificate) error {
	existing, err := s.GetCertificate(ctx, cert.CertificateNameID, cert.EnvironmentID)
	if err != nil {
		return err
	}
	if existing == nil {
		query := "INSERT INTO certificates (id, certificate_name_id, environment_id, certificate_data, private_key_secret, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
		_, err = s.db.ExecContext(ctx, core.Rebind("mysql", query), cert.ID.String(), cert.CertificateNameID.String(), cert.EnvironmentID.String(), cert.CertificateData, cert.PrivateKeySecret, cert.CreatedAt, cert.UpdatedAt)
		return err
	}

	now := time.Now()
	query := "UPDATE certificates SET certificate_data = ?, private_key_secret = ?, updated_at = ? WHERE id = ?"
	_, err = s.db.ExecContext(ctx, core.Rebind("mysql", query), cert.CertificateData, cert.PrivateKeySecret, now, existing.ID.String())
	return err
}

// SSH Keys

func (s *SecretStore) ListSSHKeyNames(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.SSHKeyName], error) {
	fieldMap := map[string]string{
		"id":        "id",
		"name":      "name_upcase",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	}
	baseQuery := "SELECT CAST(id AS CHAR(36)), CAST(project_id AS CHAR(36)), name, created_at, updated_at FROM ssh_key_names WHERE project_id = ?"
	query, listArgs := core.BuildListQuery("mysql", baseQuery, opts, fieldMap)
	finalArgs := append([]interface{}{projectID.String()}, listArgs...)

	rows, err := s.db.QueryContext(ctx, core.Rebind("mysql", query), finalArgs...)
	if err != nil {
		return core.ListResult[models.SSHKeyName]{}, err
	}
	defer rows.Close()

	var items []models.SSHKeyName
	for rows.Next() {
		var n models.SSHKeyName
		var idStr, pIdStr string
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&idStr, &pIdStr, &n.Name, &createdAt, &updatedAt); err != nil {
			return core.ListResult[models.SSHKeyName]{}, err
		}
		n.ID, _ = uuid.Parse(idStr)
		n.ProjectID, _ = uuid.Parse(pIdStr)
		if createdAt.Valid {
			n.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			t := updatedAt.Time
			n.UpdatedAt = &t
		}
		items = append(items, n)
	}

	var total int64
	countQuery, countArgs := core.BuildCountQuery("SELECT COUNT(*) FROM ssh_key_names WHERE project_id = ?", core.ListOptions{Filter: opts.Filter}, fieldMap)
	finalCountArgs := append([]interface{}{projectID.String()}, countArgs...)
	err = s.db.QueryRowContext(ctx, core.Rebind("mysql", countQuery), finalCountArgs...).Scan(&total)
	return core.ListResult[models.SSHKeyName]{Items: items, TotalCount: int(total)}, err
}

func (s *SecretStore) GetSSHKeyName(ctx context.Context, id uuid.UUID) (*models.SSHKeyName, error) {
	query := "SELECT CAST(id AS CHAR(36)), CAST(project_id AS CHAR(36)), name, created_at, updated_at FROM ssh_key_names WHERE id = ?"
	var n models.SSHKeyName
	var idStr, pIdStr string
	var createdAt, updatedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, core.Rebind("mysql", query), id.String()).Scan(&idStr, &pIdStr, &n.Name, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	n.ID, _ = uuid.Parse(idStr)
	n.ProjectID, _ = uuid.Parse(pIdStr)
	if createdAt.Valid {
		n.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		n.UpdatedAt = &t
	}
	return &n, nil
}

func (s *SecretStore) GetSSHKeyNameByName(ctx context.Context, projectID uuid.UUID, name string) (*models.SSHKeyName, error) {
	nameUpcase := strings.ToUpper(name)
	query := "SELECT CAST(id AS CHAR(36)), CAST(project_id AS CHAR(36)), name, created_at, updated_at FROM ssh_key_names WHERE project_id = ? AND name_upcase = ?"
	var n models.SSHKeyName
	var idStr, pIdStr string
	var createdAt, updatedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, core.Rebind("mysql", query), projectID.String(), nameUpcase).Scan(&idStr, &pIdStr, &n.Name, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	n.ID, _ = uuid.Parse(idStr)
	n.ProjectID, _ = uuid.Parse(pIdStr)
	if createdAt.Valid {
		n.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		n.UpdatedAt = &t
	}
	return &n, nil
}

func (s *SecretStore) CreateSSHKeyName(ctx context.Context, name *models.SSHKeyName) error {
	nameUpcase := strings.ToUpper(name.Name)
	name.NameUpcase = &nameUpcase
	query := "INSERT INTO ssh_key_names (id, project_id, name, name_upcase, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := s.db.ExecContext(ctx, core.Rebind("mysql", query), name.ID.String(), name.ProjectID.String(), name.Name, name.NameUpcase, name.CreatedAt, name.UpdatedAt)
	return err
}

func (s *SecretStore) DeleteSSHKeyName(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.ExecContext(ctx, core.Rebind("mysql", "DELETE FROM ssh_key_names WHERE id = ?"), id.String())
	return err
}

func (s *SecretStore) GetSSHKey(ctx context.Context, nameID, envID uuid.UUID) (*models.SSHKey, error) {
	query := "SELECT CAST(id AS CHAR(36)), CAST(ssh_key_name_id AS CHAR(36)), CAST(environment_id AS CHAR(36)), public_key, private_key_secret, created_at, updated_at FROM ssh_keys WHERE ssh_key_name_id = ? AND environment_id = ?"
	var key models.SSHKey
	var idStr, nIdStr, eIdStr string
	var createdAt, updatedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, core.Rebind("mysql", query), nameID.String(), envID.String()).Scan(&idStr, &nIdStr, &eIdStr, &key.PublicKey, &key.PrivateKeySecret, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	key.ID, _ = uuid.Parse(idStr)
	key.SSHKeyNameID, _ = uuid.Parse(nIdStr)
	key.EnvironmentID, _ = uuid.Parse(eIdStr)
	if createdAt.Valid {
		key.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		t := updatedAt.Time
		key.UpdatedAt = &t
	}
	return &key, nil
}

func (s *SecretStore) UpsertSSHKey(ctx context.Context, key *models.SSHKey) error {
	existing, err := s.GetSSHKey(ctx, key.SSHKeyNameID, key.EnvironmentID)
	if err != nil {
		return err
	}
	if existing == nil {
		query := "INSERT INTO ssh_keys (id, ssh_key_name_id, environment_id, public_key, private_key_secret, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
		_, err = s.db.ExecContext(ctx, core.Rebind("mysql", query), key.ID.String(), key.SSHKeyNameID.String(), key.EnvironmentID.String(), key.PublicKey, key.PrivateKeySecret, key.CreatedAt, key.UpdatedAt)
		return err
	}

	now := time.Now()
	query := "UPDATE ssh_keys SET public_key = ?, private_key_secret = ?, updated_at = ? WHERE id = ?"
	_, err = s.db.ExecContext(ctx, core.Rebind("mysql", query), key.PublicKey, key.PrivateKeySecret, now, existing.ID.String())
	return err
}
