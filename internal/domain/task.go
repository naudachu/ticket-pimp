package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskEntity struct {
	ID          uuid.UUID `gorm:"primaryKey;"`
	Title       string    `gorm:"column:title;not null"`
	Description string    `gorm:"column:description;not null"`

	Creator     string `gorm:"not null"`
	Responsible string `gorm:"not null"`

	Status Status `gorm:"column:status;not null;embeded"`

	gorm.Model
}

func (t *TaskEntity) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}

type Status int

const (
	NEW        Status = iota //0
	INPROGRESS               //1
	REVIEW                   //2
	DONE                     //3
)

func (s Status) String() string {
	switch s {
	case NEW:
		return "NEW"
	case INPROGRESS:
		return "INPROGRESS"
	case REVIEW:
		return "REVIEW"
	case DONE:
		return "DONE"
	}
	return "DRAFT"
}
