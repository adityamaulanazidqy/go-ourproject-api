package identity

import (
	"go-ourproject/models/identities/statuses"
	"gorm.io/gorm"
	"time"
)

type Masterpiece struct {
	Id         int `json:"-" gorm:"primary_key"`
	UserID     int `json:"-" gorm:"column:user_id"`
	StatusID   int `json:"-" gorm:"column:status_id"`
	ClassID    int `json:"-" gorm:"column:class_id"`
	SemesterID int `json:"-" gorm:"column:semester_id"`

	// relations
	User         Users                      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Class        Classes                    `json:"-" gorm:"foreignKey:ClassID"`
	ClassName    string                     `json:"class,omitempty" gorm:"-"`
	Status       statuses.MasterpieceStatus `json:"-" gorm:"foreignKey:StatusID"`
	StatusName   string                     `json:"status,omitempty" gorm:"-"`
	Semester     Semesters                  `json:"-" gorm:"foreignKey:SemesterID"`
	SemesterName string                     `json:"semester,omitempty" gorm:"-"`

	Files      []FileMasterpiece `json:"-" gorm:"foreignKey:MasterpieceID"`
	FilesNames []string          `json:"files" gorm:"-"`

	PublicationDate time.Time `json:"publication_date" gorm:"column:publication_date"`
	LinkGithub      string    `json:"link_github" gorm:"column:link_github"`
	ViewerCount     int       `json:"viewer_count" gorm:"column:viewer_count"`
	CreatedAt       time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt       time.Time `json:"-" gorm:"column:updated_at"`
}

type FileMasterpiece struct {
	Id            int         `json:"-" gorm:"primary_key"`
	MasterpieceID int         `json:"-" gorm:"column:masterpiece_id"`
	FilePath      string      `json:"file_path" gorm:"column:file_path"`
	Masterpiece   Masterpiece `json:"masterpiece,omitempty" gorm:"foreignKey:MasterpieceID"`
}

func (FileMasterpiece) TableName() string {
	return "files_masterpiece"
}

func (masterpiece *Masterpiece) BeforeCreate(tx *gorm.DB) (err error) {
	masterpiece.CreatedAt = time.Now()
	masterpiece.UpdatedAt = time.Now()
	return
}

func (masterpiece *Masterpiece) BeforeUpdate(tx *gorm.DB) (err error) {
	masterpiece.UpdatedAt = time.Now()
	return
}
