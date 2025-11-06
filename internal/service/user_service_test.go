package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/rangira25/user_service/internal/domain"
	"github.com/rangira25/user_service/internal/service"
)


// Mock Repository implementing all methods

type MockRepo struct {
	mock.Mock
}

// GetByID implements repository.UserRepository.
func (m *MockRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	panic("unimplemented")
}

// List implements repository.UserRepository.
func (m *MockRepo) List(ctx context.Context, filter map[string]interface{}, limit int, offset int) ([]domain.User, int64, error) {
	panic("unimplemented")
}

// Restore implements repository.UserRepository.
func (m *MockRepo) Restore(ctx context.Context, id string) error {
	panic("unimplemented")
}

// UpdateStatus implements repository.UserRepository.
func (m *MockRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	panic("unimplemented")
}

func (m *MockRepo) Create(ctx context.Context, u *domain.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockRepo) Update(ctx context.Context, u *domain.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}


// Test CreateUser

func TestCreateUser(t *testing.T) {
	mockRepo := &MockRepo{}
	svc := service.NewUserService(mockRepo)

	req := domain.CreateUserReq{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "user",
	}

	// Set mock expectations
	mockRepo.On("GetByEmail", mock.Anything, req.Email).Return(nil, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	resp, err := svc.CreateUser(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, req.Email, resp.Email)

	mockRepo.AssertExpectations(t)
}
