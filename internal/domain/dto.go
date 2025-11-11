package domain



// CreateUserReq represents the request body for creating a new user
type CreateUserReq struct {
	FullName string `json:"fullName" validate:"required,min=2,max=120" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,e164" example:"+250788123456"`
	Password string `json:"password" validate:"required,min=8" example:"StrongPassword123"`
	Role     string `json:"role,omitempty" validate:"omitempty,oneof=user admin" example:"user"`
}

// UpdateUserReq represents the request body for updating user info
type UpdateUserReq struct {
	FullName string `json:"fullName,omitempty" validate:"omitempty,min=2,max=120" example:"John Updated"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,e164" example:"+250788987654"`
}

// SetStatusReq represents the request body for updating a user's status
type SetStatusReq struct {
	Status string `json:"status" validate:"required,oneof=active suspended blocked" example:"active"`
}

// LoginReq represents the request body for logging in a user
type LoginReq struct {
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"StrongPassword123"`
}

// ============================
// 📦 Response DTOs
// ============================

// UserResp represents a user returned in responses
type UserResp struct {
	ID       string  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	FullName string  `json:"fullName" example:"John Doe"`
	Email    string  `json:"email" example:"john@example.com"`
	Phone    *string `json:"phone,omitempty" example:"+250788123456"`
	Role     string  `json:"role" example:"user"`
	Status   string  `json:"status" example:"active"`
}

// LoginResp represents a successful login response with a JWT token
type LoginResp struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ErrorResponse represents a standard error message
type ErrorResponse struct {
	Message string `json:"message" example:"Failed to process request"`
	Error   string `json:"error,omitempty" example:"email already exists"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message" example:"Operation successful"`
}

// ListUsersResponse represents a paginated list of users
type ListUsersResponse struct {
	Users []UserResp `json:"users"`
	Total int        `json:"total" example:"25"`
	Page  int        `json:"page" example:"1"`
	Limit int        `json:"limit" example:"10"`
}
