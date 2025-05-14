package response_models

import (
	identity "go-ourproject/models/identities"
	"go-ourproject/models/identities/statuses"
)

type MasterpieceResponse struct {
	User            identity.UsersResponse     `json:"user"`
	Status          statuses.MasterpieceStatus `json:"-"`
	StatusName      string                     `json:"status"`
	Class           identity.Classes           `json:"-"`
	ClassName       string                     `json:"class"`
	Semester        identity.Semesters         `json:"-"`
	SemesterName    string                     `json:"semester"`
	LinkGithub      string                     `json:"link_github"`
	Files           []string                   `json:"files"`
	PublicationDate string                     `json:"publication_date"`
}

type FileMasterpieceResponse struct {
	Name string `json:"name"`
}
