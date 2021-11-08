package services

import (
	"os"
	"time"

	"gorm.io/gorm"
	"main.go/databases"
	"main.go/models"

	// jwt
	"github.com/dgrijalva/jwt-go"
)

type AuthServices interface {
	CreateUser(email string, password []byte, name string) error
	FindByEmail(email string) (*models.User, error)
	DeleteUser(*models.User) error
	TokenGenerator(*models.User) (string, error, string)
}

type authServices struct {
	DB gorm.DB
}

// constructor
func NewAuthServices() AuthServices {
	return &authServices{DB: *databases.DB}
}

// create user
func (a *authServices) CreateUser(email string, password []byte, name string) error {
	user := &models.User{Email: email, Password: password, Name: name}
	return a.DB.Create(user).Error
}

// find user by email
func (a *authServices) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := a.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

// delete user
func (a *authServices) DeleteUser(user *models.User) error {
	return a.DB.Delete(user).Error
}

// generate token with expire time
func (a *authServices) TokenGenerator(user *models.User) (string, error, string) {
	exp := time.Now().Add(time.Second * 60)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"name":  user.Name,
		"exp":   exp.Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	return tokenString, err, exp.Format(time.RFC3339)
}
