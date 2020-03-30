package user

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

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
	Username  string    `json:"username"`
}

type LoginModel struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	HashPassword string `json:"hash_password,omitempty"`
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

	claims["name"] = u.Username

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
		Username:  u.Username,
	}
	cache.SetCc("token:"+t.Token, t, time.Hour*76)
	return e.JSON(http.StatusOK, t)

}

func Logout(e echo.Context) error {
	id := e.Get("user_id")
	fmt.Println(id)

	return e.JSON(http.StatusOK, nil)
}

type Reg struct {
	Username        string   `json:"username"`
	Password        string   `json:"password"`
	ConfirmPassword string   `json:"confirm"`
	Phone           string   `json:"phone"`
	Email           string   `json:"email"`
	Residence       []string `json:"residence"`
	Address         string   `json:"address"`
}

func Register(e echo.Context) error {
	r := new(Reg)
	err := e.Bind(r)
	if err != nil {
		return handler.ErrorResp(e, err, 500)
	}
	if store.UserExist(r.Username) {
		return handler.ErrorResp(e, fmt.Errorf("username Exist "), 403)
	}
	if !isPhone(r.Phone) {
		return handler.ErrorResp(e, fmt.Errorf("Not Vail phone fomat . "), 403)
	}
	if !isEmail(r.Email) {
		return handler.ErrorResp(e, fmt.Errorf("Not Vail Email fomat . "), 403)
	}
	if store.EmailExist(r.Email) {
		return handler.ErrorResp(e, fmt.Errorf("Email alreay registed "), 403)
	}
	switch r.Password {
	case "":
		return handler.ErrorResp(e, fmt.Errorf("password should not be empty "), 403)
	default:
		return handler.ErrorResp(e, fmt.Errorf("Password not match . "), 403)
	case r.ConfirmPassword:
	}
	var u store.UserInformation
	u.Username = r.Username
	ep, _ := auth.Encrypt(r.Password)
	u.Password = ep
	u.Email = r.Email
	u.Phone = r.Phone

	a := new(store.Address)
	a.Phone = u.Phone
	a.Name = "默认地址"
	a.Area = r.Residence[2]
	a.Address = r.Address
	a.UserID = u.Id
	u.Addresses = append(u.Addresses, *a)
	a.IsDefault = true
	u.Role = 1
	err = u.Create()
	if err != nil {
		return handler.ErrorResp(e, err, 500)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = u.Username

	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	_, _ = token.SignedString([]byte(viper.GetString("jwt_secret")))
	t := &Token{
		Token:     utils.NewUUID(),
		ExpiresAt: time.Now().Add(time.Hour * 76),
		Role:      u.Role,
		UserID:    strconv.FormatUint(uint64(u.Id), 10),
		Username:  u.Username,
	}
	cache.SetCc("token:"+t.Token, t, time.Hour*76)

	return e.JSON(http.StatusOK, t)
}

func CheckUsername(e echo.Context) error {
	key := e.Param("key")
	return handler.Response(e, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    store.UserExist(key),
	})
}

func CheckEmail(e echo.Context) error {
	key := e.Param("key")
	return handler.Response(e, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    store.EmailExist(key),
	})
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
	u.Id = id.(int32)
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
	u.Id = id.(int32)
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

func Address(ctx echo.Context) error {
	id := ctx.Get("user_id")
	var u store.UserInformation
	u.Id = id.(int32)
	_, err := u.GetUser()
	if err != nil {
		return handler.ErrorResp(ctx, err, 404)
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    u.Addresses,
	})
}
