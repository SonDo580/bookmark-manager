package model

import (
	"time"
	"uuid"

	"gorm.io/gorm"
)

type User struct {
	ID          string    `gorm:"type:uuid;primaryKey;not null;column:id" json:"id"`
	Username    string    `gorm:"type:varchar(100);uniqueIndex;not null;column:username" json:"username"`
	Password    string    `gorm:"type:text;not null;column:password" json:"-"`
	Email       string    `gorm:"type:varchar(255);uniqueIndex;not null;column:email" json:"email"`
	DisplayName string    `gorm:"type:varchar(100);not null;column:display_name" json:"display_name"`
	CreatedAt   time.Time `gorm:"type:timestamptz;column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;column:updated_at;not null;autoUpdateTime" json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
