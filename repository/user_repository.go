package repository

import (
	"github.com/yoshinori0811/chat_app/model"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	InsertUser(user *model.User) error
	GetUserByEmail(user *model.User, email string) error
	ExistsUserByEmail(email string) (bool, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{db}
}

func (ur UserRepository) InsertUser(user *model.User) error {
	sql := `INSERT INTO users (name, email, password, uuid) VALUES (?, ?, ?, ?);`
	if err := ur.db.Exec(sql, user.Name, user.Email, user.Password, user.UUID).Error; err != nil {
		// if err := ur.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (ur UserRepository) GetUserByEmail(user *model.User, email string) error {
	sql := `SELECT * FROM users WHERE email = ?`
	if err := ur.db.Raw(sql, email).First(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur UserRepository) ExistsUserByEmail(email string) (bool, error) {
	var count int64
	sql := `SELECT COUNT(*) FROM users WHERE email = ?`
	if err := ur.db.Raw(sql, email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
