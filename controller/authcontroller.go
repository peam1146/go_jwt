package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"main.go/services"
)

type AuthController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	User(c *gin.Context)
	DeleteUser(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type authController struct {
	AuthService services.AuthServices
}

// AuthController constructor
func NewAuthController(authService *services.AuthServices) AuthController {
	return &authController{AuthService: *authService}
}

// Register
func (a *authController) Register(c *gin.Context) {
	var data map[string]string
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := a.AuthService.CreateUser(data["email"], password, data["name"]); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User created successfully"})
	return

}

// Login
func (a *authController) Login(c *gin.Context) {
	var data map[string]string
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := a.AuthService.FindByEmail(data["email"])
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Println(user)
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		log.Println(err)
		return
	}
	token, err, expire := a.AuthService.TokenGenerator(user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"token": token, "expire_date": expire})
	return
}

// user
func (a *authController) User(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	user, err := a.AuthService.FindByEmail(email.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"user": user})
	return
}

// DeleteUser
func (a *authController) DeleteUser(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	user, err := a.AuthService.FindByEmail(email.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err := a.AuthService.DeleteUser(user); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted successfully"})
	return
}

// RefreshToken
func (a *authController) RefreshToken(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	user, err := a.AuthService.FindByEmail(email.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	token, err, expire := a.AuthService.TokenGenerator(user)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"token": token, "expire_date": expire})
}
