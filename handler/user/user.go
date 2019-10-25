package user

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/jackyczj/July/cache"

	uuid "github.com/satori/go.uuid"

	"github.com/jackyczj/July/pkg/auth"
	"github.com/jackyczj/July/store"
	"github.com/jackyczj/July/utils"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type Token struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    string    `json:"-"`
}

type LoginModel struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	HashPassword string `json:"hash_password,omitempty"`
}

func init() {
	store.Client.Init()
}

func Login(e echo.Context) error {
	username := e.FormValue("username")
	password := e.FormValue("password")

	if username == "" || password == "" {
		return echo.NewHTTPError(401, "AuthFailed,username or password can't be empty")
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	s := store.UserInformation{}
	if isPhone(username) {
		s = store.UserInformation{
			Phone: username,
		}
	} else if isEmail(username) {
		s = store.UserInformation{
			Email: username,
		}
	} else {
		s = store.UserInformation{
			Username: username,
		}
	}
	u, err := s.GetUser()
	if err != nil {
		return echo.NewHTTPError(401, "AuthFailed,username not found.")
	}
	if a := auth.Compare(u.Password, password); a != nil {
		return echo.NewHTTPError(401, "AuthFailed,password incorrect.")
	}

	claims["name"] = username

	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	_, err = token.SignedString([]byte(viper.GetString("jwt_secret")))
	if err != nil {
		return echo.NewHTTPError(401, "Can't get user id.")
	}
	id, err := u.GetId()
	if err != nil {
		return err
	}
	t := &Token{
		Token:     utils.NewUUID(),
		ExpiresAt: time.Now().Add(time.Hour * 76),
		UserID:    id,
	}
	cache.SetCc("Token"+t.Token, t, time.Hour*76)
	return e.JSON(http.StatusOK, t)

}

func Register(e echo.Context) error {

	user := store.UserInformation{}
	user.Lock()
	username := e.FormValue("username")
	if username == "" {
		return echo.NewHTTPError(401, "Username can't be empty")
	}
	password := e.FormValue("password")
	if password == "" {
		return echo.NewHTTPError(401, "Password can't be empty")
	}
	phone := e.FormValue("phone")
	if !isPhone(phone) {
		return echo.NewHTTPError(401, "Phone format error")
	}
	email := e.FormValue("mail")
	if !isEmail(email) {
		return echo.NewHTTPError(401, "Email format error")
	}
	g := e.FormValue("gander")
	gander, err := strconv.Atoi(g)
	if err != nil {
		return echo.NewHTTPError(401, "Gander invalid")
	}
	user.Username = username
	user.Email = email
	user.Password, _ = auth.Encrypt(password)
	user.Id = uuid.NewV1().String()
	user.Gander = gander

	return nil
}

func isPhone(phone string) bool {
	ok, err := regexp.MatchString("[1][358][0-9]{9}", phone)
	if err != nil {
		return false
	}
	return ok
}

func isEmail(mail string) bool {
	isEmail, err := regexp.MatchString("[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,4}", mail)
	if err != nil {
		return false
	}
	return isEmail
}

func Get(e echo.Context) error {
	id := e.Get("user_id")
	u := new(store.UserInformation)
	u.Id = id.(string)
	u, err := u.GetUser()
	if err != nil {
		return echo.NewHTTPError(401, "Get user information fail")
	}
	u.Password = ""
	return e.JSON(200, u)
}

func Set(e echo.Context) error {
	id := e.Get("user_id")
	u := new(store.UserInformation)
	u.Id = id.(string)
	u, err := u.GetUser()
	if err != nil {
		return echo.NewHTTPError(401, "Edit user information fail.")
	}
	phone := e.FormValue("phone")
	switch phone {
	case "":
	default:
		if isPhone(phone) {
			u.Phone = phone
			break
		}
		return echo.NewHTTPError(401, "Phone format error")
	}

	email := e.FormValue("mail")
	switch email {
	case "":
	default:
		if isEmail(email) {
			u.Email = email
			break
		}
		return echo.NewHTTPError(401, "Email format error")
	}
	g := e.FormValue("gander")
	if g != "" {
		gander, err := strconv.Atoi(g)
		if err != nil {
			return echo.NewHTTPError(401, "Gander invalid")
		}
		u.Gander = gander
	}
	err = u.Set()
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(401, "Edit fail , may system error.")
	}
	message := struct {
		Message string `json:"message"`
	}{
		Message: "Edit success",
	}
	return e.JSON(200, message)
}

func VailPhone(phone string) bool {
	//todo: send message
	return true
}

func VailEmail(email string) bool {
	//todo:send message
	return true
}
