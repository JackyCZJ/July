package store

import (
	"fmt"
	"testing"
	"time"
)

func TestShop_Create(t *testing.T) {
	shop := Shop{
		Name:        fmt.Sprintf("测试小店%v", 3),
		Owner:       29677,
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
		Id: "5e8b1f6ba1b1bb082367d698",
	}
	if shop.Get() != nil {
		t.Fatal(shop.Get())
	}
	fmt.Println(shop)

}

func TestShop_ShopModify(t *testing.T) {
	shop := Shop{
		Id:          "5e775a8e6eeb4543911c2c03",
		Name:        "测试小店1",
		Owner:       29677,
		CreateAt:    time.Now(),
		Description: "测试用小店",
		IsClose:     true,
		IsDelete:    false,
	}
	if err := shop.ShopModify(); err != nil {
		t.Fatal(err.Error())
	}
}

func TestShopList1(t *testing.T) {
	data, count, err := ShopList(1, 1)
	fmt.Println(err)
	fmt.Println(count)
	fmt.Println(data)
}

func TestShop_Set(t *testing.T) {
	var s Shop
	s.Id = "5e8b1f6ba1b1bb082367d698"
	err := s.Set("is_close", false)
	if err != nil {
		t.Fatal(err)
	}
}
