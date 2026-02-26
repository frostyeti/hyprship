package stores

import (
	"context"

	"github.com/frostyeti/hyprship/apps/api/internal/core"
	"github.com/frostyeti/hyprship/apps/api/internal/models"
	"github.com/google/uuid"
)

type EnvironmentStore interface {
	List(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.Environment], error)
	Get(ctx context.Context, id uuid.UUID) (*models.Environment, error)
	GetByName(ctx context.Context, projectID uuid.UUID, name string) (*models.Environment, error)
	Create(ctx context.Context, env *models.Environment) error
	Update(ctx context.Context, env *models.Environment) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ConfigStore interface {
	// Config Files
	ListConfigFileNames(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.ConfigFileName], error)
	GetConfigFileName(ctx context.Context, id uuid.UUID) (*models.ConfigFileName, error)
	GetConfigFileNameByName(ctx context.Context, projectID uuid.UUID, name string) (*models.ConfigFileName, error)
	CreateConfigFileName(ctx context.Context, name *models.ConfigFileName) error
	DeleteConfigFileName(ctx context.Context, id uuid.UUID) error

	GetConfigFile(ctx context.Context, nameID, envID uuid.UUID) (*models.ConfigFile, error)
	UpsertConfigFile(ctx context.Context, file *models.ConfigFile) error

	// Env Variables
	ListEnvVariableNames(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.EnvVariableName], error)
	GetEnvVariableName(ctx context.Context, id uuid.UUID) (*models.EnvVariableName, error)
	GetEnvVariableNameByName(ctx context.Context, projectID uuid.UUID, name string) (*models.EnvVariableName, error)
	CreateEnvVariableName(ctx context.Context, name *models.EnvVariableName) error
	DeleteEnvVariableName(ctx context.Context, id uuid.UUID) error

	GetEnvVariable(ctx context.Context, nameID, envID uuid.UUID) (*models.EnvVariable, error)
	UpsertEnvVariable(ctx context.Context, v *models.EnvVariable) error
}

type SecretStore interface {
	// Secrets
	ListSecretNames(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.SecretName], error)
	GetSecretName(ctx context.Context, id uuid.UUID) (*models.SecretName, error)
	GetSecretNameByName(ctx context.Context, projectID uuid.UUID, name string) (*models.SecretName, error)
	CreateSecretName(ctx context.Context, name *models.SecretName) error
	DeleteSecretName(ctx context.Context, id uuid.UUID) error

	GetSecret(ctx context.Context, nameID, envID uuid.UUID) (*models.Secret, error)
	UpsertSecret(ctx context.Context, secret *models.Secret) error

	// Certificates
	ListCertificateNames(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.CertificateName], error)
	GetCertificateName(ctx context.Context, id uuid.UUID) (*models.CertificateName, error)
	GetCertificateNameByName(ctx context.Context, projectID uuid.UUID, name string) (*models.CertificateName, error)
	CreateCertificateName(ctx context.Context, name *models.CertificateName) error
	DeleteCertificateName(ctx context.Context, id uuid.UUID) error

	GetCertificate(ctx context.Context, nameID, envID uuid.UUID) (*models.Certificate, error)
	UpsertCertificate(ctx context.Context, cert *models.Certificate) error

	// SSH Keys
	ListSSHKeyNames(ctx context.Context, projectID uuid.UUID, opts core.ListOptions) (core.ListResult[models.SSHKeyName], error)
	GetSSHKeyName(ctx context.Context, id uuid.UUID) (*models.SSHKeyName, error)
	GetSSHKeyNameByName(ctx context.Context, projectID uuid.UUID, name string) (*models.SSHKeyName, error)
	CreateSSHKeyName(ctx context.Context, name *models.SSHKeyName) error
	DeleteSSHKeyName(ctx context.Context, id uuid.UUID) error

	GetSSHKey(ctx context.Context, nameID, envID uuid.UUID) (*models.SSHKey, error)
	UpsertSSHKey(ctx context.Context, key *models.SSHKey) error
}
