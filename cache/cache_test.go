package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/jackyczj/July/config"
)

func init() {
	err := config.Init("../conf/config.yaml")
	if err != nil {
		panic(err)
	}
}

func TestCache_Test(t *testing.T) {
	InitCache()
	s, err := Cluster.Ping().Result()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
	SetCc("cc", "cache", 3*time.Second)
	err = GetCc("cc", &s)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s)
}
