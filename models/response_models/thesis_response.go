package response_models

import (
	identity "go-ourproject/models/identities"
	"go-ourproject/models/identities/statuses"
)

type ThesisResponse struct {
	Users      identity.Users        `json:"users"`
	Status     statuses.ThesisStatus `json:"-"`
	StatusName string                `json:"status"`
}
