package user

import (
	"net/http"
	"regexp"
	"time"

	"github.com/jackyczj/NoGhost/pkg/auth"
	"github.com/jackyczj/NoGhost/store"
	"github.com/jackyczj/NoGhost/utils"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type Token struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    string    `json:"-"`
}

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
	var s store.UserInformation
	isEmail, err := regexp.MatchString("/[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,4}/", username)
	if err != nil {
		return err
	}
	if isPhone {
		s.Phone = username
	} else if isEmail {
		s.Email = username
	} else if isUsername {
		s.Username = username
	}
	u, err := s.GetUser()
	if err != nil {
		return echo.NewHTTPError(401, "AuthFailed", "username not found.")
	}
	if a := auth.Compare(u.Password, password); a != nil {
		return echo.NewHTTPError(401, "AuthFailed", "password incorrect.")
	}

	claims["name"] = username

	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	_, err = token.SignedString([]byte(viper.GetString("jwt_secret")))
	if err != nil {
		return err
	}
	id, err := u.GetId()
	if err != nil {
		return err
	}
	t := &Token{
		Token:     utils.NewUUID(),
		ExpiresAt: time.Now().Add(time.Hour * 96),
		// 这个userid应该是检索出来的，这里为demo写死。
		UserID: id,
	}
	return e.JSON(http.StatusOK, t)

}

func Register(e echo.Context) error {
	return nil
}
