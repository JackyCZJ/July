package store

import (
	"fmt"
	"testing"
	"time"
)

//
//func TestOrder_Create(t *testing.T) {
//	var s Order
//
//}

func TestOrder_Create(t *testing.T) {
	o := Order{
		Seller:      "5e788d61234cca37ca522ee6",
		Buyer:       31209,
		Payment:     3222,
		PaymentType: 0,
		ShippingTo:  0,
		Item: []Item{
			{
				"5e832e5b041c8d2f73392b25", 1,
			},
		},
		CreateAt: time.Now(),
		Status:   UnPay,
	}
	err := o.Create()
	if err != nil {
		t.Fatal(err)
	}
}

func TestOrder_Get(t *testing.T) {
	var o Order
	o.OrderNo = "5e85c6a1c2dc168f9abe1575"
	err := o.Get()
	fmt.Println(err)
	fmt.Println(o)
}

func TestOrder_Update(t *testing.T) {
	var o Order
	o.OrderNo = "5e831bb6aa9d2a09f68295fe"
	err := o.Update("status", "2")
	if err != nil {
		t.Fatal(err)
	}
}
