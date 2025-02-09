package repository

import (
	"github.com/yoshinori0811/chat_app_backend/model"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Insert(user *model.User) error
	GetByEmail(user *model.User, email string) error
	ExistsUserByEmail(email string) (bool, error)
	GetByName(name string, id uint) ([]model.User, error)
	GetUserIDByName(name string) (uint, error)
	GetUserIDsByNames(nameList []string) ([]uint, error)
	GetUserByID(user *model.User) error
	GetUserNameByID(id uint) (string, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{db}
}

func (ur UserRepository) Insert(user *model.User) error {
	sql := `INSERT INTO users (name, email, password, uuid) VALUES (?, ?, ?, ?);`
	if err := ur.db.Exec(sql, user.Name, user.Email, user.Password, user.UUID).Error; err != nil {
		return err
	}
	return nil
}

func (ur UserRepository) GetByEmail(user *model.User, email string) error {
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

func (ur UserRepository) GetByName(name string, id uint) ([]model.User, error) {
	users := []model.User{}
	sql := `SELECT * FROM users WHERE name LIKE ? AND id != ? LIMIT 30`
	if err := ur.db.Raw(sql, "%"+name+"%", id).Scan(&users).Error; err != nil {
		return []model.User{}, err
	}
	return users, nil
}

func (ur UserRepository) GetUserIDByName(name string) (uint, error) {
	var userID uint
	sql := `SELECT id FROM users WHERE name = ?`
	if err := ur.db.Raw(sql, name).First(&userID).Error; err != nil {
		return 0, err
	}
	return userID, nil
}

func (ur UserRepository) GetUserIDsByNames(nameList []string) ([]uint, error) {
	var IDList []uint
	sql := `SELECT id FROM users WHERE name IN (?)`
	if err := ur.db.Raw(sql, nameList).Scan(&IDList).Error; err != nil {
		return IDList, err
	}
	return IDList, nil
}

func (ur UserRepository) GetUserByID(user *model.User) error {
	sql := `SELECT * FROM users WHERE id = ?`
	if err := ur.db.Raw(sql, user.ID).Scan(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur UserRepository) GetUserNameByID(id uint) (string, error) {
	var name string
	sql := `SELECT name FROM users WHERE id = ?`
	if err := ur.db.Raw(sql, id).First(&name).Error; err != nil {
		return "", err
	}
	return name, nil
}
