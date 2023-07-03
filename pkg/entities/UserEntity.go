package entities

import (
	"gorm.io/gorm"
	"time"
)

type UserEntity struct {
	ID         uint           `gorm:"primaryKey" json:"id,omitempty"`
	RoleID     uint           `json:"-"`
	Role       RoleEntity     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Username   string         `gorm:"size:255;not null;unique;constraint:OnUpdate:CASCADE;" json:"username,omitempty"`
	Password   string         `gorm:"size:255;not null;" json:"-"`
	Email      string         `gorm:"email;not null;unique;constraint:OnUpdate:CASCADE" json:"email,omitempty"`
	Avatar     string         `json:"avatar,omitempty"`
	NIM        *string        `gorm:"constraint:OnUpdate:CASCADE;null;unique" json:"nim"`
	NIP        *string        `gorm:"constraint:OnUpdate:CASCADE;null;unique" json:"nip"`
	Major      string         `gorm:"not null" json:"major"`
	ResetToken string         `json:"-"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty"`
}
