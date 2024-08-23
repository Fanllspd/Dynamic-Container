package middleware

import (
	"k3s-client/app/services"
	"k3s-client/global"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuth(GuardName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is required"})
			c.Abort()
			return
		}
		tokenStr = tokenStr[len(services.TokenType)+1:]

		token, err := jwt.ParseWithClaims(tokenStr, &services.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(global.App.Config.Jwt.Secret), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*services.CustomClaims)

		if claims.Issuer != GuardName || !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is invalid"})
			c.Abort()
			return
		}
		if claims.UserId <= 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id is required"})
			c.Abort()
			return
		}
		log.Printf("claims.UserId: %d", claims.UserId)
		// c.Set("token", token)

		// if claims.Role != "admin" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "permission denied"})
		// 	return
		// }
		// c.Set("role", claims.Role)

		// if claims.UserId  == 0 {
		c.Set("userId", claims.UserId)
	}
}
