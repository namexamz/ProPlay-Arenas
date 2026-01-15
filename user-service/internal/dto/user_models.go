package dto


type Role string

const (
	RoleOwner  Role = "Owner"
	RoleClient Role = "Client"
	RoleAdmin  Role = "Admin"
)

type CreateUserRequest struct {
    FullName string `json:"full_name" binding:"required,min=2,max=100"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Role     Role `json:"role" binding:"omitempty,oneof=client owner admin"`
}

type UpdateUserRequest struct {
    FullName *string `json:"full_name" binding:"omitempty,min=2,max=100"`
    Email    *string `json:"email" binding:"omitempty,email"`
    Role     *Role `json:"role" binding:"omitempty,oneof=client owner admin"`
}

type UserResponse struct {
    ID       uint   `json:"id"`
    FullName string `json:"full_name"`
    Email    string `json:"email"`
    Role     Role `json:"role"`
}