package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"primaryKey;"`
	Telegram     string    `gorm:"column:tg"`
	TelegramLink string    `gorm:"column:tglink"`

	gorm.Model
}

func (u *User) BeforeCreate(tx *gorm.Tx) (err error) {
	u.ID = uuid.New()
	return nil
}
