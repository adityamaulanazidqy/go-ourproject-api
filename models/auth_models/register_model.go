package auth_models

type RegisterRequest struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	RoleId   *int    `json:"-"`
	MajorId  *int    `json:"-"`
	Batch    *int    `json:"-"`
	Photo    *string `json:"-"`
}

type RegisterResponse struct {
	Id       int    `json:"-" gorm:"primary_key"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Batch    int    `json:"batch"`
	Major    string `json:"major"`
	Photo    string `json:"photo"`
}
