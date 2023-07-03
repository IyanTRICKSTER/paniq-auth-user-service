package entities

import (
	"gorm.io/gorm"
	"paniq-auth-user-service/pkg/contracts/apiResources"
	"paniq-auth-user-service/pkg/contracts/permissionCodes"
	"time"
)

type PermissionEntity struct {
	ID        uint                           `gorm:"primaryKey" json:"-"`
	Name      string                         `json:"name"`
	Code      permissionCodes.PermissionCode `json:"code"`
	Resource  apiResources.RESOURCE          `json:"resource"`
	Roles     []RoleEntity                   `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
	CreatedAt time.Time                      `json:"-"`
	UpdatedAt time.Time                      `json:"-"`
	DeletedAt gorm.DeletedAt                 `gorm:"index" json:"-"`
}
