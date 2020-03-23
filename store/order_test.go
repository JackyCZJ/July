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
		Seller:      "5e775a8e6eeb4543911c2c03",
		Buyer:       32519,
		Payment:     0,
		PaymentType: 0,
		ShippingTo:  0,
		Item:        nil,
		CreateTime:  time.Now(),
		Status:      "3",
		TrackingNum: "123091823701983",
	}
	err := o.Create()
	if err != nil {
		t.Fatal(err)
	}
}
