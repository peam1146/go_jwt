package middlewares

// jwt middleware
import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.AbortWithStatus(401)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte("secret"), nil
		})
		if err != nil {
			c.AbortWithStatus(401)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("name", claims["name"])
			c.Set("id", claims["id"])
			c.Set("email", claims["email"])
			c.Next()
		} else {
			c.AbortWithStatus(401)
			return
		}
	}
}
