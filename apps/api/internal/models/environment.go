package models

import (
	"time"

	"github.com/google/uuid"
)

type Environment struct {
	ID          uuid.UUID  `json:"id"`
	ProjectID   uuid.UUID  `json:"projectId"`
	Name        string     `json:"name"`
	NameUpcase  *string    `json:"-"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type ConfigFileName struct {
	ID         uuid.UUID  `json:"id"`
	ProjectID  uuid.UUID  `json:"projectId"`
	Name       string     `json:"name"`
	NameUpcase *string    `json:"-"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

type ConfigFile struct {
	ID               uuid.UUID  `json:"id"`
	ConfigFileNameID uuid.UUID  `json:"configFileNameId"`
	EnvironmentID    uuid.UUID  `json:"environmentId"`
	Content          string     `json:"content"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
}

type EnvVariableName struct {
	ID         uuid.UUID  `json:"id"`
	ProjectID  uuid.UUID  `json:"projectId"`
	Name       string     `json:"name"`
	NameUpcase *string    `json:"-"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

type EnvVariable struct {
	ID                uuid.UUID  `json:"id"`
	EnvVariableNameID uuid.UUID  `json:"envVariableNameId"`
	EnvironmentID     uuid.UUID  `json:"environmentId"`
	Value             string     `json:"value"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
}

type SecretName struct {
	ID         uuid.UUID  `json:"id"`
	ProjectID  uuid.UUID  `json:"projectId"`
	Name       string     `json:"name"`
	NameUpcase *string    `json:"-"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

type Secret struct {
	ID            uuid.UUID  `json:"id"`
	SecretNameID  uuid.UUID  `json:"secretNameId"`
	EnvironmentID uuid.UUID  `json:"environmentId"`
	Value         string     `json:"value"` // Note: This will be encrypted at rest
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt,omitempty"`
}

type CertificateName struct {
	ID         uuid.UUID  `json:"id"`
	ProjectID  uuid.UUID  `json:"projectId"`
	Name       string     `json:"name"`
	NameUpcase *string    `json:"-"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

type Certificate struct {
	ID                uuid.UUID  `json:"id"`
	CertificateNameID uuid.UUID  `json:"certificateNameId"`
	EnvironmentID     uuid.UUID  `json:"environmentId"`
	CertificateData   string     `json:"certificateData"`
	PrivateKeySecret  *string    `json:"privateKeySecret,omitempty"` // Note: Encrypted at rest
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
}

type SSHKeyName struct {
	ID         uuid.UUID  `json:"id"`
	ProjectID  uuid.UUID  `json:"projectId"`
	Name       string     `json:"name"`
	NameUpcase *string    `json:"-"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

type SSHKey struct {
	ID               uuid.UUID  `json:"id"`
	SSHKeyNameID     uuid.UUID  `json:"sshKeyNameId"`
	EnvironmentID    uuid.UUID  `json:"environmentId"`
	PublicKey        string     `json:"publicKey"`
	PrivateKeySecret string     `json:"privateKeySecret"` // Note: Encrypted at rest
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
}
