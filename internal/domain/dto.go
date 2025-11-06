package domain

type CreateUserReq struct {
	FullName string `json:"fullName" validate:"required,min=2,max=120"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,e164"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role,omitempty" validate:"omitempty,oneof=user admin"`
}

type UpdateUserReq struct {
	FullName string `json:"fullName,omitempty" validate:"omitempty,min=2,max=120"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,e164"`
}

type SetStatusReq struct {
	Status string `json:"status" validate:"required,oneof=active suspended blocked"`
}

// Response DTOs
type UserResp struct {
	ID       string  `json:"id"`
	FullName string  `json:"fullName"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone,omitempty"`
	Role     string  `json:"role"`
	Status   string  `json:"status"`
}
