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
		s.Buyer = u.Username
		s.CreateTime = time.Now().UTC()
		s.Payment = v.Product.Price * float64(v.Count)
		s.Status = UnPay
		s.IsClose = false
		s.ShippingTo = u.Addresses[O.AddressIndex]
		s.Seller = v.Product.Owner
		cache.SetCc("Order:"+s.OrderNo, s, 12*time.Hour) //存入缓存，12小时内未完成便自动删除
		return s.Create()
	}
	return nil
}

func List() {

}

func Search() {

}
