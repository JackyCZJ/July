package store

import (
	"fmt"
	"testing"
	"time"
)

func TestShop_Create(t *testing.T) {
	shop := Shop{
		Name:        "测试小店",
		Owner:       "JackyTest",
		CreateAt:    time.Now(),
		Description: "测试用小店",
		IsClose:     false,
		IsDelete:    false,
	}
	err := shop.Create()
	if err != nil {
		t.Fatal(err)
	}
}

func TestShop_Delete(t *testing.T) {
	shop := Shop{
		Name: "测试小店",
	}
	err := shop.Delete()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSearchShop(t *testing.T) {
	s, count, err := SearchShop("测试", 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s, count)
}

func TestShopList(t *testing.T) {
	s, count, err := ShopList(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s, count)
}
