package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) (string, error) {
	hashedPassword, error := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedPassword), error
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var AccessTokenSecret = []byte(os.Getenv("TOKEN_SECRET"))
var RefreshTokenSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))
var AccessTokenTTL = time.Minute * 15
var RefreshTokenTTL = time.Hour * 24 * 7

type UserClaims struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func NewAccessToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func NewRefreshToken(claims jwt.StandardClaims) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseAccessToken(accessToken string) *UserClaims {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil || !parsedAccessToken.Valid {
		return nil
	}

	return parsedAccessToken.Claims.(*UserClaims)
}

func ParseRefreshToken(refreshToken string) *jwt.StandardClaims {
	parsedRefreshToken, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil || !parsedRefreshToken.Valid {
		return nil
	}

	return parsedRefreshToken.Claims.(*jwt.StandardClaims)
}
