package identity

type Classes struct {
	Id    int    `json:"id" gorm:"primary_key"`
	Class string `json:"name" gorm:"unique"`
}
