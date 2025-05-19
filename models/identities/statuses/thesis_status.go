package statuses

type ThesisStatus struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"unique"`
}

func (ThesisStatus) TableName() string {
	return "thesis_statuses"
}
