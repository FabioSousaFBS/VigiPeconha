package users

type CreateUserRequest struct {
	Name     string  `json:"name" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Phone    *string `json:"phone"`
	Password string  `json:"password" binding:"required,min=8"`
	Role     string  `json:"role" binding:"required"`
}

type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type CreateUserParams struct {
	Name           string
	Email          string
	Phone          *string
	PasswordHash   string
	Role           string
	OrganizationID string
}
