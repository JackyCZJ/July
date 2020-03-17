package store

import (
	"fmt"
	"testing"
	"time"
)

func TestProduct_Add(t *testing.T) {
	p := Product{
		Name:        "wtds342123",
		ImageUri:    []string{"https://via.placeholder.com/150", "https://via.placeholder.com/200"},
		Description: "dfkjalkdfj;a",
		Information: Type{
			"不知道啥",
			"不知道啥",
		},
		Store:    10,
		Price:    1213,
		Off:      1223,
		Owner:    "1231",
		CreateAt: time.Now(),
		Shelves:  true,
	}

	err := p.Add()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSearch(t *testing.T) {
	i, c, e := Search("1222", 1, 10)
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(i, c)
}

func TestProduct_Get(t *testing.T) {
	var p Product
	p.ProductId = "123123"
	err := p.Get()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(p)
}

func TestGetRandom(t *testing.T) {
	i, e := GetRandom()
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(i)
	fmt.Println(len(i))
}

func TestGetListByShop(t *testing.T) {
	p := GetListByShop("1231", true)
	if len(p) == 0 {
		t.Fatal("Get fail")
	}
	fmt.Println(p)
}

func TestSuggestion(t *testing.T) {
	p := Suggestion("wt")
	if len(p) == 0 {
		t.Fatal("Get fail")
	}
	fmt.Println(p)
}
