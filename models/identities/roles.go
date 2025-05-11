package identity

type Roles struct {
	Id   int    `json:"-" gorm:"primary_key"`
	Name string `json:"name" gorm:"unique"`
}
