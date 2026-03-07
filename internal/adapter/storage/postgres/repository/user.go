package repository

import (
	"context"

	"github.com/yehezkiel1086/go-redis-job-listings/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-redis-job-listings/internal/core/domain"
)

type UserRepository struct {
	db *postgres.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	db := r.db.GetDB()

	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	db := r.db.GetDB()

	var user domain.User
	if err := db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	db := r.db.GetDB()

	var user domain.User
	if err := db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	db := r.db.GetDB()

	if err := db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) DeleteUserByID(ctx context.Context, id uint) error {
	db := r.db.GetDB()

	if err := db.WithContext(ctx).Delete(&domain.User{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	db := r.db.GetDB()

	var users []domain.User
	if err := db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
