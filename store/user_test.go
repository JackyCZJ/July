package store

import (
	"testing"
)

func TestUserInformation_Create(t *testing.T) {
	Client.Init()
	defer Client.Close()
	ui := UserInformation{
		Username: "Jacky",
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
		Username: "Jacky",
	}
	u, err := Um.GetUser()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
}
