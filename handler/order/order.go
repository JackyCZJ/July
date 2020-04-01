package order

import (
	"time"

	"github.com/jackyczj/July/cache"

	"github.com/jackyczj/July/handler"
	"github.com/jackyczj/July/handler/cart"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
)

type Order struct {
	OrderList    []cart.Request `json:"order_list"`
	PaymentType  int            `json:"payment_type" bson:"payment_type,omitempty" `
	AddressIndex int            `json:"address_index"`
}

const UnPay = "0"
const UnShip = "1"
const SHIPPING = "2"
const RECEIVED = "3"
const COMMENTED = "4"

func Create(ctx echo.Context) error {
	id := ctx.Get("user_id").(int32)
	var u store.UserInformation
	u.Id = id
	_, err := u.GetUser()
	if err != nil {
		return handler.ErrorResp(ctx, err, 403)
	}
	var O Order
	for _, v := range O.OrderList {
		var s store.Order
		s.Buyer = u.Id
		s.CreateAt = time.Now()
		s.Payment = v.Product.Price * float64(v.Count)
		s.Status = UnPay
		s.IsClose = false
		s.ShippingTo = O.AddressIndex
		s.Seller = v.Product.Owner
		cache.SetCc("Order:"+s.OrderNo, s, 12*time.Hour) //存入缓存，12小时内未完成便自动删除
		return s.Create()
	}
	return nil
}

func List(ctx echo.Context) error {
	u := ctx.Get("user_id").(int32)
	r := ctx.Get("role").(int)
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    store.OrderList(u, r),
	})
}

func Get(ctx echo.Context) error {
	id := ctx.Get(":id").(string)
	o := store.Order{}
	o.OrderNo = id
	err := o.Get()
	if err != nil {
		return handler.ErrorResp(ctx, err, 404)
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    o,
	})
}

func Transmit(ctx echo.Context) error {
	id := ctx.Param("id")
	tn := ctx.Param("tn")
	o := store.Order{}
	o.OrderNo = id
	err := o.Get()
	if err != nil {
		return handler.ErrorResp(ctx, err, 404)
	}

	o.Status = SHIPPING
	o.TrackingNum = tn
	err = o.UpdateAll()
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    nil,
	})
}

func Delete(ctx echo.Context) error {
	orderId := ctx.Param(":id")
	o := store.Order{}
	o.OrderNo = orderId
	err := o.Delete()
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    nil,
	})
}

func Edit(ctx echo.Context) error {
	orderId := ctx.Param(":id")
	var O store.Order
	err := ctx.Bind(&O)
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	O.OrderNo = orderId
	err = O.UpdateAll()
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    nil,
	})
}

func Cancel(ctx echo.Context) error {
	return nil
}

func Confirm(ctx echo.Context) error {
	return nil
}
