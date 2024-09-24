package services

import (
	"github.com/my_ecommerce/internal/dto"
	"github.com/my_ecommerce/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func (user *UserService) InitUserService(database *gorm.DB) {
	// giving an instance of db to service
	user.db = database
	// adding model in table if not present
	user.db.AutoMigrate(&models.User{})
}

func (u *UserService) GetUser(id int) (*dto.UserResponse, error) {
	var user *models.User

	if err := u.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	userResponse := &dto.UserResponse{
		Name:    user.Name,
		Email:   user.Email,
		Address: user.Address,
		ID:      user.ID,
		Role:    user.Role,
	}

	return userResponse, nil
}

func (u *UserService) CreateUser(name string, password string, email string, address string, role string) (*dto.UserResponse, error) {

	hassedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     name,
		Password: string(hassedPassword),
		Email:    email,
		Address:  address,
		Role:     role,
	}

	if err := u.db.Create(user).Error; err != nil {
		return nil, err
	}

	userResponse := &dto.UserResponse{
		Name:    user.Name,
		Email:   user.Email,
		Address: user.Address,
		ID:      user.ID,
		Role:    user.Role,
	}

	return userResponse, nil
}

/*
This will return user while setting a jwt token in cookie
*/
func (u *UserService) LoginUser(email string, password string) (*dto.UserResponse, error) {

	var user *models.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	// creating new user
	existingUser := &dto.UserResponse{
		Name:    user.Name,
		Email:   user.Email,
		Address: user.Address,
		ID:      user.ID,
		Role:    user.Role,
	}

	return existingUser, nil
}
