package identity

import (
	"time"
)

type Supervision struct {
	Id        int       `json:"id" gorm:"primary_key"`
	ThesisID  int       `json:"thesis_id"`
	Thesis    Thesis    `json:"-" gorm:"foreignkey:ThesisID"`
	TeacherID int       `json:"teacher_id"`
	Teacher   Users     `json:"-" gorm:"foreignkey:TeacherID"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type SupervisionRequest struct {
	ThesisID  int       `json:"thesis_id"`
	TeacherID int       `json:"-,omitempty"`
	Notes     string    `json:"notes"`
	StatusID  int       `json:"status_id"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type SupervisionResponse struct {
	Thesis      Thesis      `json:"thesis"`
	Supervision Supervision `json:"supervision"`
	Status      string      `json:"status"`
}

func (Supervision) TableName() string {
	return "supervision"
}
