package user

import (
	"regexp"
	"time"

	"github.com/jackyczj/NoGhost/pkg/auth"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func Login(e echo.Context) error {

	username := e.FormValue("username")
	password := e.FormValue("password")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	isPhone, err := regexp.MatchString("/[1][358][0-9]{9}/", username)
	if err != nil {
		return err
	}
	isUsername, err := regexp.MatchString("/[A-Za-z][A-Za-z0-9]{12}/", username)
	if err != nil {
		return err
	}
	isEmail, err := regexp.MatchString("/[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,4}/", username)
	if isPhone {

	} else if isEmail {

	}
	if a := auth.Compare(password, password); a != nil {
		return a
	}

	claims["name"] = username

	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	_, err = token.SignedString([]byte(viper.GetString("jwt_secret")))
	if err != nil {
		return err
	}

	return nil

}

func Register(e echo.Context) error {
	return nil
}
