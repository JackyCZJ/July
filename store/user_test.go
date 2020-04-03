package store

import (
	"fmt"
	"testing"

	cacheClient "github.com/jackyczj/July/cache"

	"github.com/jackyczj/July/config"
)

func init() {
	Client.Init()
	err := config.Init("../conf/config.yaml")
	if err != nil {
		panic(err)
	}
	cacheClient.InitCache()

}
func TestUserInformation_Create(t *testing.T) {
	ui := UserInformation{
		Username: "JackyTest",
		Password: "wtfIsPassword",
		Email:    "test@test.com",
		Role:     2,
	}
	err := ui.Create()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserInformation_GetUser(t *testing.T) {
	Um := UserInformation{
		Username: "JackyTest",
	}
	u, err := Um.GetUser()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(u)
}

func TestUserInformation_Set(t *testing.T) {
	Um := UserInformation{
		Username: "JackyTest",
	}
	err := Um.Set("email", "test@test")
	if err != nil {
		t.Fatal(err)
	}
	u, err := Um.GetUser()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(u)
}

func TestUserInformation_del(t *testing.T) {

	Um := UserInformation{
		Username: "JackyTest",
	}
	err := Um.Delete()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserExist(t *testing.T) {
	fmt.Println(UserExist("JackyTest"))
}

func TestUserInformation_ChangeRole(t *testing.T) {
	var u UserInformation
	u.Id = 23086
	us, err := u.GetUser()
	if err != nil {
		t.Fatal(err)
	}
	err = us.ChangeRole(1)
	if err != nil {
		t.Fatal(err)

	}
}
