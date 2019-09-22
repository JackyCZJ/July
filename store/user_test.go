package store

import (
	"testing"
)

func TestUserInformation_Create(t *testing.T) {
	Client.Init()
	defer Client.Close()
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
	Client.Init()
	defer Client.Close()

	Um := UserInformation{
		Username: "JackyTest",
	}
	u, err := Um.GetUser()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u.Password)
}

func TestUserInformation_GetId(t *testing.T) {
	Client.Init()
	defer Client.Close()

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
	Client.Init()
	defer Client.Close()
	Um := UserInformation{
		Username: "JackyTest",
	}
	err := Um.del()
	if err != nil {
		t.Fatal(err)
	}
}
