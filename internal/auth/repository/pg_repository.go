package auth_repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"hson98/app-chat/internal/models"
)

type authRepo struct {
	db *gorm.DB
}

func (a *authRepo) FindUser(context context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.User, error) {
	var user models.User
	db := a.db.Table(models.User{}.TableName())
	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}
	if err := db.Where(conditions).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (a *authRepo) CreateUser(context context.Context, user *models.User) (*models.User, error) {
	db := a.db.Begin()
	var userInsert models.User
	if err := db.Table(user.TableName()).Create(&user).Scan(&userInsert).Error; err != nil {
		db.Rollback()
		return nil, err
	}
	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return nil, err
	}

	return &userInsert, nil
}

type Repository interface {
	CreateUser(context context.Context, user *models.User) (*models.User, error)
	FindUser(context context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.User, error)
}

func NewAuthRepo(db *gorm.DB) Repository {
	return &authRepo{db: db}
}
