package identity

type Semesters struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"unique"`
}
