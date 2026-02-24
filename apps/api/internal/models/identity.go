package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID  `json:"id"`
	PrimaryEmail       *string    `json:"primaryEmail,omitempty"`
	PrimaryEmailUpcase *string    `json:"-"`
	PrimaryPhone       *string    `json:"primaryPhone,omitempty"`
	Name               *string    `json:"name,omitempty"`
	NameUpcase         *string    `json:"-"`
	ImageURI           *string    `json:"imageUri,omitempty"`
	IsBanned           bool       `json:"isBanned"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          *time.Time `json:"updatedAt,omitempty"`
}

type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	NameUpcase  string    `json:"-"`
	Description string    `json:"description"`
}

type UserClaim struct {
	ID     int32     `json:"id"`
	UserID uuid.UUID `json:"userId"`
	Type   string    `json:"type"`
	Value  string    `json:"value"`
}

type RoleClaim struct {
	ID     int32     `json:"id"`
	RoleID uuid.UUID `json:"roleId"`
	Type   string    `json:"type"`
	Value  string    `json:"value"`
}

type UserAPIKey struct {
	ID          int32      `json:"id"`
	UserID      uuid.UUID  `json:"userId"`
	Name        string     `json:"name"`
	NameUpcase  string     `json:"-"`
	KeyHint     string     `json:"keyHint"`
	KeyDigest   string     `json:"-"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
	IsLocked    bool       `json:"isLocked"`
	TotpSecret  *string    `json:"-"`
	TotpEnabled bool       `json:"totpEnabled"`
	IsRevoked   bool       `json:"isRevoked"`
	Comment     *string    `json:"comment,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type UserPasswordAuth struct {
	UserID            uuid.UUID  `json:"userId"`
	PasswordDigest    string     `json:"-"`
	PasswordExpiresAt *time.Time `json:"passwordExpiresAt,omitempty"`
	LastAttemptedAt   *time.Time `json:"lastAttemptedAt,omitempty"`
	AttemptCount      int        `json:"attemptCount"`
	OtpDigest         *string    `json:"-"`
	OtpExpiresAt      *time.Time `json:"otpExpiresAt,omitempty"`
	OtpLinkToken      *string    `json:"-"`
	IsLocked          bool       `json:"isLocked"`
	TotpSecret        *string    `json:"-"`
	TotpEnabled       bool       `json:"totpEnabled"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
}

type UserSession struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"userId"`
	ExpiresAt time.Time  `json:"expiresAt"`
	Token     string     `json:"-"`
	IPAddress *string    `json:"ipAddress,omitempty"`
	UserAgent *string    `json:"userAgent,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type UserPasskey struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"userId"`
	CredentialID []byte    `json:"credentialId"` // Stored as binary
	Data         []byte    `json:"data"`         // Stored as JSON
	CreatedAt    time.Time `json:"createdAt"`
}

type UserTotp struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"userId"`
	Secret     string     `json:"-"`
	IsVerified bool       `json:"isVerified"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}
