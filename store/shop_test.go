package store

import (
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
		Name:        "测试小店",
		Owner:       "JackyTest",
		CreateAt:    time.Now(),
		Description: "测试用小店",
		IsClose:     false,
		IsDelete:    false,
	}
	err := shop.Delete()
	if err != nil {
		t.Fatal(err)
	}
}
