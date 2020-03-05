package store

import (
	"fmt"
	"testing"
	"time"
)

func TestProduct_Add(t *testing.T) {
	p := Product{
		Name:        "wtds33112",
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
		IsDelete: false,
	}

	err := p.Add()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSearch(t *testing.T) {
	i, e := Search("1222", 1, 10)
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(i)
}

func TestProduct_Get(t *testing.T) {
	var p Product
	p.ProductId = 3864
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
