package files

import identity "go-ourproject/models/identities"

type FileMasterpiece struct {
	Id            int                  `json:"-" gorm:"primary_key"`
	MasterpieceID int                  `json:"-" gorm:"column:masterpiece_id"`
	FilePath      string               `json:"file_path" gorm:"column:file_path"`
	Masterpiece   identity.Masterpiece `json:"masterpiece,omitempty" gorm:"foreignKey:MasterpieceID"`
}

func (FileMasterpiece) TableName() string {
	return "files_masterpiece"
}
