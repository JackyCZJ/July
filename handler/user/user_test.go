package user

import "testing"

func Test_isEmail(t *testing.T) {
	email := "test@test.com"
	if isEmail(email) {
		t.Log("Pass")
	} else {
		t.Fatal("Fail")
	}
	email = "test@test"
	if !isEmail(email) {
		t.Log("Pass")
	} else {
		t.Fatal("Fail")
	}
}

func Test_isPhone(t *testing.T) {
	phone := "13010101010"
	if isPhone(phone) {
		t.Log("Pass")
	} else {
		t.Fatal("Fail")
	}
	phone = "1232131"
	if !isPhone(phone) {
		t.Log("Pass")
	} else {
		t.Fatal("Fail")
	}
}
