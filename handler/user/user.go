package user

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/rs/xid"

	"github.com/jackyczj/July/cache"

	"github.com/jackyczj/July/handler"

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
	Role      int       `json:"role"`
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
	lm := new(LoginModel)
	err := e.Bind(&lm)
	if err != nil {
		return echo.NewHTTPError(401, "AuthFailed,username or password can't be empty")
	}

	if lm.Username == "" || lm.Password == "" {
		return echo.NewHTTPError(401, "AuthFailed,username or password can't be empty")
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	s := store.UserInformation{}
	//if isPhone(username) {
	//	s = store.UserInformation{
	//		Phone: username,
	//	}
	//} else
	if isEmail(lm.Username) {
		s = store.UserInformation{
			Email: lm.Username,
		}
	} else {
		s = store.UserInformation{
			Username: lm.Username,
		}
	}
	u, err := s.GetUser()
	if err != nil {
		return echo.NewHTTPError(401, "AuthFailed,username not found.")
	}
	if a := auth.Compare(u.Password, lm.Password); a != nil {
		return echo.NewHTTPError(401, "AuthFailed,password incorrect.")
	}

	claims["name"] = lm.Username

	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	_, err = token.SignedString([]byte(viper.GetString("jwt_secret")))
	if err != nil {
		return echo.NewHTTPError(401, "Can't get user id.")
	}
	t := &Token{
		Token:     utils.NewUUID(),
		ExpiresAt: time.Now().Add(time.Hour * 76),
		Role:      u.Role,
		UserID:    strconv.FormatUint(uint64(u.Id), 10),
	}
	cache.SetCc("token:"+t.Token, t, time.Hour*76)
	return e.JSON(http.StatusOK, t)

}

func Logout(e echo.Context) error {
	id := e.Get("user_id")
	fmt.Println(id)
	return e.JSON(http.StatusOK, nil)
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
	//phone := e.FormValue("phone")
	//if !isPhone(phone) {
	//	return echo.NewHTTPError(401, "Phone format error")
	//}
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
	user.Id = xid.New().Pid()
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
	u.Id = id.(uint16)
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
	u.Id = id.(uint16)
	u, err := u.GetUser()
	if err != nil {
		return echo.NewHTTPError(401, "Edit user information fail.")
	}
	err = e.Bind(&u)
	if err != nil {
		return echo.NewHTTPError(401, "Edit user information fail.")
	}
	if isEmail(u.Email) {
		err = u.Set("email", u.Email)
		if err != nil {
			return echo.NewHTTPError(401, "Edit user information fail.")
		}
	}
	return e.JSON(200, "Edit success")
}

func UsernameCheck(e echo.Context) error {
	user := e.Get("username").(string)
	var um store.UserInformation
	um.Username = user
	_, err := um.GetUser()
	res := handler.ResponseStruct{
		Code:    0,
		Message: "username don't exist",
		Data:    nil,
	}
	if err != nil {
		return handler.Response(e, res)
	}
	res.Code = 1
	res.Message = "username exist"
	return handler.Response(e, res)
}

func VailPhone(phone string) bool {
	//todo: send message
	return true
}

func VailEmail(email string) bool {
	//todo:send message
	return true
}
