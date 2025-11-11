package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rangira25/user_service/internal/config"
	"github.com/rangira25/user_service/internal/domain"
	"github.com/rangira25/user_service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)



// Predefined errors
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidCreds     = errors.New("invalid credentials")
	ErrEmailAlreadyUsed = errors.New("email already in use")
)

// UserService defines public interface (called by handler)
type UserService interface {
	CreateUser(ctx context.Context, req domain.CreateUserReq) (*domain.UserResp, error)
	GetUser(ctx context.Context, id string) (*domain.UserResp, error)
	ListUsers(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]domain.UserResp, int64, error)
	UpdateUser(ctx context.Context, id string, req domain.UpdateUserReq) (*domain.UserResp, error)
	SetStatus(ctx context.Context, id, status string) error
	DeleteUser(ctx context.Context, id string) error
	RestoreUser(ctx context.Context, id string) error
	AdminResetPassword(ctx context.Context, id, newPassword string) error
	Login(ctx context.Context, req domain.LoginReq) (*domain.LoginResp, error)
}

// Concrete implementation
type userService struct {
	repo repository.UserRepository
}

// Constructor
func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

// ===================== CRUD Operations =====================

func (s *userService) CreateUser(ctx context.Context, req domain.CreateUserReq) (*domain.UserResp, error) {
	if existing, _ := s.repo.GetByEmail(ctx, req.Email); existing != nil {
		return nil, ErrEmailAlreadyUsed
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		FullName:     req.FullName,
		Email:        req.Email,
		Role:         req.Role,
		PasswordHash: ptrString(string(hash)),
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	logAudit("CREATE", user.ID)
	return toResp(user), nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*domain.UserResp, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return toResp(u), nil
}

func (s *userService) ListUsers(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]domain.UserResp, int64, error) {
	users, total, err := s.repo.List(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	out := make([]domain.UserResp, 0, len(users))
	for _, u := range users {
		out = append(out, *toResp(&u))
	}
	return out, total, nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, req domain.UpdateUserReq) (*domain.UserResp, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.FullName != "" {
		u.FullName = req.FullName
	}
	if req.Phone != "" {
		u.Phone = &req.Phone
	}

	u.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}

	logAudit("UPDATE", u.ID)
	return toResp(u), nil
}

func (s *userService) SetStatus(ctx context.Context, id, status string) error {
	err := s.repo.UpdateStatus(ctx, id, status)
	if err == nil {
		logAudit("SET_STATUS:"+status, id)
	}
	return err
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err == nil {
		logAudit("DELETE", id)
	}
	return err
}

func (s *userService) RestoreUser(ctx context.Context, id string) error {
	err := s.repo.Restore(ctx, id)
	if err == nil {
		logAudit("RESTORE", id)
	}
	return err
}

func (s *userService) AdminResetPassword(ctx context.Context, id, newPassword string) error {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	h := string(hash)
	u.PasswordHash = &h

	if err := s.repo.Update(ctx, u); err != nil {
		return err
	}

	logAudit("RESET_PASSWORD", id)
	return nil
}

// ===================== LOGIN (JWT Auth) =====================

func (s *userService) Login(ctx context.Context, req domain.LoginReq) (*domain.LoginResp, error) {
	// 1. Get user by email
	u, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 2. Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(*u.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCreds
	}

	// 3. Use system current time for token
	
	// 4. Create JWT claims
	claims := jwt.MapClaims{
    "user_id": fmt.Sprintf("%d", u.ID),
    "role":    u.Role,
    "exp":     jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
    "iat":     jwt.NewNumericDate(time.Now()),
}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.JWTSecret))

	// 5. Generate JWT token
	
	if err != nil {
		return nil, err
	}

	logAudit("LOGIN", u.ID)
	return &domain.LoginResp{Token: tokenStr}, nil
}

// ===================== Helpers =====================

func toResp(u *domain.User) *domain.UserResp {
	return &domain.UserResp{
		ID:       u.ID,
		FullName: u.FullName,
		Email:    u.Email,
		Phone:    u.Phone,
		Role:     u.Role,
		Status:   u.Status,
	}
}

func ptrString(s string) *string {
	return &s
}
