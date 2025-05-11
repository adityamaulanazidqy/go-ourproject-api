package auth_models

import identity "go-ourproject/models/identities"

type LoginRequest struct {
	Email    string `json:"email" example:"user@gmail.com"`
	Password string `json:"password" example:"password"`
}

type LoginResponse struct {
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Role     identity.Roles  `json:"role"`
	Major    identity.Majors `json:"major"`
}
