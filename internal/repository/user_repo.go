package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/rangira25/user_service/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]domain.User, int64, error)
	Update(ctx context.Context, u *domain.User) error
	UpdateStatus(ctx context.Context, id, status string) error
	Delete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

// Implement methods (Create, GetByID, etc.)
func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) List(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]domain.User, int64, error) {
	var users []domain.User
	q := r.db.WithContext(ctx).Model(&domain.User{})
	// apply filters (status, role, search)
	for k, v := range filter {
		switch k {
		case "status", "role":
			q = q.Where(k+" = ?", v)
		case "search":
			// search across name/email/phone
			q = q.Where("full_name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", "%"+v.(string)+"%", "%"+v.(string)+"%", "%"+v.(string)+"%")
		}
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepo) Update(ctx context.Context, u *domain.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *userRepo) UpdateStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("status", status).Error
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.User{}, "id = ?", id).Error
}

func (r *userRepo) Restore(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil).Error
}
