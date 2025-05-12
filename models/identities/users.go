package identity

import (
	"gorm.io/gorm"
	"time"
)

type Users struct {
	Id        int       `json:"-" gorm:"primary_key"`
	Username  string    `json:"username"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password"`
	RoleID    uint8     `json:"-" gorm:"column:role_id"`
	MajorID   uint8     `json:"-" gorm:"column:major_id"`
	Batch     int       `json:"batch"`
	Photo     string    `json:"photo"`
	CreatedAt time.Time `json:"-,omitempty" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"-,omitempty" gorm:"column:updated_at"`
	Role      Roles     `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Major     Majors    `json:"major,omitempty" gorm:"foreignKey:MajorID"`
}
type UsersResponse struct {
	Id        int       `json:"-" gorm:"primary_key"`
	Username  string    `json:"username"`
	Email     string    `json:"email" gorm:"unique"`
	RoleID    uint8     `json:"-" gorm:"column:role_id"`
	MajorID   uint8     `json:"-" gorm:"column:major_id"`
	Batch     int       `json:"batch"`
	Photo     string    `json:"photo"`
	CreatedAt time.Time `json:"-,omitempty" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"-,omitempty" gorm:"column:updated_at"`
	Role      Roles     `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Major     Majors    `json:"major,omitempty" gorm:"foreignKey:MajorID"`
}

func (user *Users) BeforeCreate(tx *gorm.DB) (err error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return
}

func (user *Users) BeforeUpdate(tx *gorm.DB) (err error) {
	user.UpdatedAt = time.Now()
	return
}
