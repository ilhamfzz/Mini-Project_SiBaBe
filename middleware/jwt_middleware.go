package middleware

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateToken(username string, name string) (string, error) {
	claims := jwt.MapClaims{}
	claims["username"] = username
	claims["nama"] = name
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_JTW")))
}

func ExtractTokenUsername(e echo.Context) (string, error) {
	user := e.Get("user").(*jwt.Token) //
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	return username, nil
}
