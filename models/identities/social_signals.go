package identity

type Likes struct {
	Id            int          `json:"id" gorm:"primary_key"`
	MasterpieceID int          `json:"masterpiece_id" gorm:"foreignkey:masterpiece_id"`
	Count         int          `json:"count"`
	Masterpiece   *Masterpiece `json:"-" gorm:"foreignkey:MasterpieceID"`
}

type Dislike struct {
	Id            int          `json:"id" gorm:"primary_key"`
	MasterpieceID int          `json:"masterpiece_id" gorm:"foreignkey:masterpiece_id"`
	Count         int          `json:"count"`
	Masterpiece   *Masterpiece `json:"-" gorm:"foreignkey:MasterpieceID"`
}
