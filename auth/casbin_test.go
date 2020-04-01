package auth

import (
	"fmt"
	"testing"

	"github.com/casbin/casbin/v2"
)

func TestCasbin(t *testing.T) {
	enforcer, err := casbin.NewEnforcer("../conf/casbin_auth_model.conf", "../conf/casbin_auth_policy.csv")
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println(enforcer.Enforce("2", "/api/v1/Cart/List", "GET"))
}
