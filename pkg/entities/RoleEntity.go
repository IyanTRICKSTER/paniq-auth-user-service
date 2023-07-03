package entities

import (
	"gorm.io/gorm"
	"time"
)

type RoleEntity struct {
	ID          uint               `gorm:"primaryKey" json:"id,omitempty"`
	Name        string             `gorm:"unique" json:"name,omitempty"`
	Users       []UserEntity       `gorm:"foreignKey:RoleID;references:ID" json:"users,omitempty"`
	Permissions []PermissionEntity `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt     `gorm:"index" json:"deleted_at,omitempty"`
}
