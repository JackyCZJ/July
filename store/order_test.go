package store

import (
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
		Buyer:       29677,
		Payment:     0,
		PaymentType: 0,
		ShippingTo:  0,
		Item:        nil,
		CreateAt:    time.Now(),
		Status:      UnPay,
	}
	err := o.Create()
	if err != nil {
		t.Fatal(err)
	}
}
