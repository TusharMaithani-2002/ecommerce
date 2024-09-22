package services

import (

	internal "github.com/my_ecommerce/internal/models"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"

)

type UserService struct {
	db *gorm.DB
}

func (user *UserService) InitUserService(database *gorm.DB) {
	// giving an instance of db to service
	user.db = database
	// adding model in table if not present
	user.db.AutoMigrate(&internal.User{})
}

type User struct {
	ID string
	Name string
	Email string
	Address string
}

func (u *UserService) GetUser(id int) (*internal.User,error) {
	var user *internal.User

	if err := u.db.Where("id = ?",id).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) CreateUser(name string, password string, email string, address string) (*internal.User,error) {

	hassedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}
	
	user := &internal.User{
		Name: name,
		Password: string(hassedPassword),
		Email: email,
		Address: address,
	}

	if err := u.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
/*
This will return user while setting a jwt token in cookie
*/
func (u *UserService) LoginUser(email string, password string) (*internal.User,error) {

	var user *internal.User
	if err := u.db.Where("email = ?",email).First(&user).Error; err != nil {
		return nil, err
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password))
	if err != nil {
		return nil, err
	}

	// creating new user
	existingUser := &internal.User{
		Name: user.Name,
		Email:user.Email,
		Address: user.Address,
		ID:user.ID,
	}

	return existingUser, nil
}