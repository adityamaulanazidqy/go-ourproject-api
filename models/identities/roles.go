package identity

type Roles struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"unique"`
}
