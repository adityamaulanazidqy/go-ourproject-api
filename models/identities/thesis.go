package identity

import (
	"go-ourproject/models/identities/statuses"
	"time"
)

type Thesis struct {
	Id          int                   `json:"id" gorm:"primary_key"`
	UserID      int                   `json:"user_id,omitempty"`
	User        Users                 `json:"-" gorm:"foreignkey:UserID"`
	TeacherID   int                   `json:"teacher_id,omitempty" gorm:"column:teacher_id"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	StatusID    int                   `json:"status_id,omitempty"`
	Status      statuses.ThesisStatus `json:"-" gorm:"foreignkey:StatusID"`
	CreatedAt   time.Time             `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time             `json:"updated_at" gorm:"autoUpdateTime"`
}

type ThesisRequest struct {
	UserID      int    `json:"user_id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ThesisID struct {
	ThesisID int `json:"thesis_id"`
}

func (Thesis) TableName() string {
	return "thesis"
}
