package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	NameUpcase  *string    `json:"-"`
	Slug        string     `json:"slug"`
	Description *string    `json:"description,omitempty"`
	IsActive    bool       `json:"isActive"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type ProjectGroup struct {
	ProjectID   uuid.UUID `json:"projectId"`
	GroupID     uuid.UUID `json:"groupId"`
	Permissions int64     `json:"permissions"`
}
