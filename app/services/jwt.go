package services

import (
	"k3s-client/global"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type jwtService struct {
}

var JwtService = new(jwtService)

type JwtUser interface {
	GetUid() string
}
type CustomClaims struct {
	UserId uint `json:"user_id"`
	// Role   string `json:"role"`
	jwt.StandardClaims
}

const (
	TokenType    = "bearer"
	AppGuardName = "app"
)

type TokenOutPut struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// 给
func (jwtService *jwtService) CreateToken(guardName string, userId uint /*, role string*/) (tokenData TokenOutPut, err error, token *jwt.Token) {

	token = jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		CustomClaims{
			UserId: userId,
			// Role:   role,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() + int64(global.App.Config.Jwt.JwTtl),
				// Id:        user.GetUid(),
				Issuer:    guardName,
				NotBefore: time.Now().Unix() - 1000,
				Subject:   "access_token",
			},
		},
	)

	tokenStr, err := token.SignedString([]byte(global.App.Config.Jwt.Secret))

	tokenData = TokenOutPut{
		tokenStr,
		int64(global.App.Config.Jwt.JwTtl),
		TokenType,
	}
	return
}
