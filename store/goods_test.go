package store

import (
	"fmt"
	"testing"
	"time"
)

func TestProduct_Add(t *testing.T) {
	p := Product{
		Name:        "wtds332",
		ImageUri:    "https://via.placeholder.com/150",
		Description: "dfkjalkdfj;a",
		Information: Type{
			"不知道啥",
			"不知道啥",
		},
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

func TestGetRandom(t *testing.T) {
	i, e := GetRandom()
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(i)
	fmt.Println(len(i))
}
