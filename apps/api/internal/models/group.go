package models

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	NameUpcase  *string    `json:"-"`
	Description *string    `json:"description,omitempty"`
	ImageURI    *string    `json:"imageUri,omitempty"`
	Email       *string    `json:"email,omitempty"`
	EmailUpcase *string    `json:"-"`
	IsActive    bool       `json:"isActive"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type GroupUser struct {
	GroupID uuid.UUID `json:"groupId"`
	UserID  uuid.UUID `json:"userId"`
}

type GroupAdmin struct {
	GroupID uuid.UUID `json:"groupId"`
	UserID  uuid.UUID `json:"userId"`
}

type GroupRole struct {
	GroupID uuid.UUID `json:"groupId"`
	RoleID  uuid.UUID `json:"roleId"`
}
