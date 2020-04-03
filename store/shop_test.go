package store

import (
	"fmt"
	"testing"
	"time"
)

func TestShop_Create(t *testing.T) {
	shop := Shop{
		Name:        fmt.Sprintf("测试小店%v", 11),
		Owner:       31209,
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

func TestShop_Get(t *testing.T) {
	shop := Shop{
		Id: "5e775a8e6eeb4543911c2c03",
	}
	if shop.Get() != nil {
		t.Fatal(shop.Get())
	}
	fmt.Println(shop)

}

func TestShop_ShopModify(t *testing.T) {
	shop := Shop{
		Id:          "5e775a8e6eeb4543911c2c03",
		Name:        "测试小店2",
		Owner:       32519,
		CreateAt:    time.Now(),
		Description: "测试用小店",
		IsClose:     false,
		IsDelete:    false,
	}
	if err := shop.ShopModify(); err != nil {
		t.Fatal(err.Error())
	}
}

func TestShopList1(t *testing.T) {
	_, count, err := ShopList(1, 1)
	fmt.Println(err)
	fmt.Println(count)
}
