package identity

import "time"

type Comments struct {
	Id            int         `json:"id" gorm:"primary_key"`
	UserId        int         `json:"-" gorm:"column:user_id"`
	MasterpieceID int         `json:"masterpiece_id" gorm:"column:masterpiece_id"`
	User          Users       `json:"user,omitempty" gorm:"foreignKey:UserId"`
	Masterpiece   Masterpiece `json:"-" gorm:"foreignKey:MasterpieceID"`
	Message       string      `json:"message" gorm:"message"`
	CreatedAt     time.Time   `json:"created_at" gorm:"autoCreateTime"`
}

func (Comments) TableName() string {
	return "comments"
}
