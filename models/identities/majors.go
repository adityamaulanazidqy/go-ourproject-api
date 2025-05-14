package identity

type Majors struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"unique"`
}
