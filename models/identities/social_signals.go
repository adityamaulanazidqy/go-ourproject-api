package social_signals

import identity "go-ourproject/models/identities"

type Likes struct {
	Id            int                  `json:"id" gorm:"primary_key"`
	MasterpieceID int                  `json:"masterpiece_id" gorm:"foreignkey:masterpiece_id"`
	Count         int                  `json:"count"`
	Masterpiece   identity.Masterpiece `json:"-" gorm:"foreignkey:MasterpieceID"`
}
