package auth_models

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleId   int    `json:"role_id"`
	MajorId  int    `json:"major_id"`
	Batch    int    `json:"batch"`
	Photo    string `json:"photo"`
}
