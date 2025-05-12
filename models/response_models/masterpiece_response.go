package response_models

import (
	identity "go-ourproject/models/identities"
	"go-ourproject/models/identities/statuses"
	"time"
)

type MasterpieceResponse struct {
	User            identity.UsersResponse     `json:"user"`
	Status          statuses.MasterpieceStatus `json:"status"`
	Class           identity.Classes           `json:"class"`
	Semester        identity.Semesters         `json:"semester"`
	LinkGithub      string                     `json:"link_github"`
	Files           []FileMasterpieceResponse  `json:"files"`
	PublicationDate time.Time                  `json:"publication_date"`
}

type FileMasterpieceResponse struct {
	Name string `json:"name"`
}
