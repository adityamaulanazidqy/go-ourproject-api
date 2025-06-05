package identity

import (
	"gorm.io/gorm"
	"time"
)

type Users struct {
	Id        int       `json:"-" gorm:"primary_key"`
	Username  string    `json:"username"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"`
	RoleID    uint8     `json:"-" gorm:"column:role_id"`
	MajorID   uint8     `json:"-" gorm:"column:major_id"`
	Batch     int       `json:"batch"`
	Photo     string    `json:"photo" gorm:"default:'icon_default.jpg'"`
	CreatedAt time.Time `json:"-,omitempty" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"-,omitempty" gorm:"column:updated_at"`
	Role      Roles     `json:"-" gorm:"foreignKey:RoleID"`
	RoleName  string    `json:"role" gorm:"-"`
	Major     Majors    `json:"-" gorm:"foreignKey:MajorID"`
	MajorName string    `json:"major" gorm:"-"`
}
type UsersResponse struct {
	Id        int       `json:"id" gorm:"primary_key"`
	Username  string    `json:"username"`
	Email     string    `json:"email" gorm:"unique"`
	RoleID    uint8     `json:"-" gorm:"column:role_id"`
	MajorID   uint8     `json:"-" gorm:"column:major_id"`
	Batch     int       `json:"batch"`
	Photo     string    `json:"photo"`
	CreatedAt time.Time `json:"-,omitempty" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"-,omitempty" gorm:"column:updated_at"`
	Role      Roles     `json:"-" gorm:"foreignKey:RoleID"`
	RoleName  string    `json:"role"`
	Major     Majors    `json:"-" gorm:"foreignKey:MajorID"`
	MajorName string    `json:"major"`
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
