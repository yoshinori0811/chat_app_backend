package repository

import (
	"errors"

	"github.com/yoshinori0811/chat_app/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user *model.User) error
	GetUserByEmail(user *model.User, email string) error
	UpdateSession(user *model.User, sessionID string) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db}
}

func (ur UserRepository) CreateUser(user *model.User) error {
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
		// if err := ur.db.Where("email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			*user = model.User{}
			return nil
		}
		return err
	}
	return nil
}

func (ur UserRepository) UpdateSession(user *model.User, sessionID string) error {
	sql := `UPDATE users SET session_id = ? WHERE email = ?`
	if err := ur.db.Raw(sql, sessionID, user.Email).Scan(&user).Error; err != nil {
		return err
	}
	return nil
}
