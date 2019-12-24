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
	}
	err := ui.Create()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserModel_GetUser(t *testing.T) {
	Um := UserInformation{
		Username: "JackyTest",
	}
	u, err := Um.GetUser()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(u)
}

func TestUserInformation_GetId(t *testing.T) {

	Um := UserInformation{
		Username: "JackyTest",
	}
	u, err := Um.GetId()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
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
