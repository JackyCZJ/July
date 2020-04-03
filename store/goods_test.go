package store

import (
	"fmt"
	"testing"
	"time"
)

func TestProduct_Add(t *testing.T) {
	p := Product{
		Name:        "awsl",
		ImageUri:    []string{"http://localhost:2333/image/5e808b33b5e26ddfb46164b9", "http://localhost:2333/image/5e8321dfd95740e6bdba977f"},
		Description: "dfkjalkdfj;a",
		Information: Type{
			"不知道啥",
			"不知道啥",
		},
		Store:    10,
		Price:    1213,
		Off:      1223,
		Owner:    "5e85d088626760d4407aa390",
		CreateAt: time.Now(),
		Shelves:  true,
	}
	for i := 1; i < 100; i++ {
		p.Price = 4333
		p.Store = 4333
		p.CreateAt = time.Now()
		p.Name = fmt.Sprintf("魅族%v", i)
		err := p.Add()
		if err != nil {
			t.Fatal(err)
		}
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
	p, i := GetListByShop("5e85d088626760d4407aa390", true, 1)
	if i == 0 {
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
