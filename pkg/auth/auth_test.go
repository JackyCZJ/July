package auth

import "testing"

func TestAuth(t *testing.T) {
	password := "testPassword"
	enPwd, err := Encrypt(password)
	if err != nil {
		t.Fatal(err)
	}
	err = Compare(enPwd, password)
	if err != nil {
		t.Fatal(err)
	}
}
