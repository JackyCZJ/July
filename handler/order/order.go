package order

import (
	"fmt"
	"time"

	"github.com/jackyczj/July/log"

	"github.com/jackyczj/July/handler"
	"github.com/jackyczj/July/handler/cart"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
)

type Order struct {
	OrderList    []cart.Request `json:"order_list" bson:"order_list"`
	PaymentType  int            `json:"payment_type" bson:"payment_type,omitempty" `
	AddressIndex int            `json:"address_index"`
}

const UnPay = "0"
const UnShip = "1"
const SHIPPING = "2"
const RECEIVED = "3"
const COMMENTED = "4"
const CANCEL = "5"

func Create(ctx echo.Context) error {
	id := ctx.Get("user_id").(int32)
	var u store.UserInformation
	u.Id = id
	_, err := u.GetUser()
	if err != nil {
		return handler.ErrorResp(ctx, err, 403)
	}
	var O Order
	err = ctx.Bind(&O)
	if err != nil {
		log.Logworker.Error(err)
		return handler.ErrorResp(ctx, err, 500)
	}
	OList := make(map[string]*Order)
	productCache := map[string]store.Product{}
	for _, v := range O.OrderList {
		var p store.Product
		p.ProductId = v.Product
		err := p.Get()
		if err != nil {
			log.Logworker.Error(err)
			return handler.ErrorResp(ctx, err, 500)
		}
		productCache[p.ProductId] = p //商品缓存
		if OList[p.Owner] == nil {
			OList[p.Owner] = new(Order)
		}
		OList[p.Owner].OrderList = append(OList[p.Owner].OrderList, v) //订单按商铺分类
		go func() {
			_ = p.Set("store", p.Store-v.Count)
		}()
		err = store.CartDel(id, p.ProductId)
		if err != nil {
			log.Logworker.Error(err)
			return handler.ErrorResp(ctx, err, 500)
		}
	}

	if len(OList) > 0 {
		for v := range OList {
			l := len(OList[v].OrderList)
			ItemList := make([]store.Item, l)
			var o store.Order
			o.Payment = 0
			for i := range OList[v].OrderList {
				var item store.Item
				item.ProductId = OList[v].OrderList[i].Product
				item.Count = OList[v].OrderList[i].Count
				o.Payment = o.Payment + productCache[item.ProductId].Price*float64(item.Count)
				ItemList = append(ItemList, item)
			}
			o.Item = ItemList
			o.Buyer = id
			o.Seller = v
			o.Status = UnPay
			o.ShippingTo = O.AddressIndex
			o.CreateAt = time.Now()
			err := o.Create()
			if err != nil {
				return handler.ErrorResp(ctx, err, 500)
			}
		}
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    nil,
	})
}

func List(ctx echo.Context) error {
	u := ctx.Get("user_id").(int32)
	if ctx.Path() == "/api/v1/Shop/order/list" {
		return handler.Response(ctx, handler.ResponseStruct{
			Code:    0,
			Message: "",
			Data:    store.ShopOrderList(u),
		})
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    store.OrderList(u),
	})
}

func Get(ctx echo.Context) error {
	id := ctx.Param("id")
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
	id := ctx.QueryParam("id")
	tn := ctx.QueryParam("tn")
	o := store.Order{}
	o.OrderNo = id
	err := o.Get()
	if err != nil {
		return handler.ErrorResp(ctx, err, 404)
	}
	var s store.Shop
	s.Owner = ctx.Get("user_id").(int32)
	err = s.GetByOwner()
	if err != nil {
		return handler.ErrorResp(ctx, err, 403)
	}
	if o.Seller != s.Id {
		return handler.ErrorResp(ctx, fmt.Errorf("该订单不属于你的店铺 "), 403)
	}

	//o.Status = SHIPPING
	//o.TrackingNum = tn
	//err = o.UpdateAll()
	//if err != nil {
	//	return handler.ErrorResp(ctx, err, 500)
	//}
	err = o.Update("status", SHIPPING)
	err2 := o.Update("tracking_num", tn)
	if err != nil && err2 != nil {
		log.Logworker.Error(err, err2)
		return handler.ErrorResp(ctx, err, 404)
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
	id := ctx.Get("user_id").(int32)
	var u store.UserInformation
	u.Id = id
	_, err := u.GetUser()
	if err != nil {
		return handler.ErrorResp(ctx, err, 403)
	}
	orderId := ctx.Param(":id")
	var O store.Order
	err = ctx.Bind(&O)
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

func Pay(ctx echo.Context) error {
	id := ctx.Param("id")
	u := ctx.Get("user_id").(int32)
	var o store.Order
	o.OrderNo = id
	err := o.Get()
	if err != nil {
		return handler.ErrorResp(ctx, fmt.Errorf("Order not found "), 404)
	}
	if o.Buyer != u {
		return handler.ErrorResp(ctx, fmt.Errorf("Not Your Order "), 403)
	}
	err = o.Update("status", UnShip)
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "支付成功（呸）",
		Data:    nil,
	})
}

func Cancel(ctx echo.Context) error {
	id := ctx.Get("user_id").(int32)
	role := ctx.Get("role").(int)
	var u store.UserInformation
	u.Id = id
	_, err := u.GetUser()
	if err != nil {
		return handler.ErrorResp(ctx, err, 403)
	}
	var O store.Order
	O.OrderNo = ctx.Param("id")
	err = O.Get()
	if err != nil {
		return handler.ErrorResp(ctx, fmt.Errorf("未找到该订单 "), 404)
	}
	if role > 1 {
		var shop store.Shop
		shop.Owner = id
		_ = shop.GetByOwner()
		if O.Buyer != u.Id || O.Seller != shop.Id {
			err = O.Update("status", CANCEL)
			if err != nil {
				return handler.ErrorResp(ctx, fmt.Errorf("未知错误，订单取消失败 "), 500)
			}
			return handler.Response(ctx, handler.ResponseStruct{
				Code:    0,
				Message: "success",
				Data:    nil,
			})
		}
	}
	if O.Buyer != u.Id {
		return handler.ErrorResp(ctx, fmt.Errorf("此订单非你所拥有 "), 403)
	}
	err = O.Update("status", CANCEL)
	if err != nil {
		return handler.ErrorResp(ctx, fmt.Errorf("未知错误，订单取消失败 "), 500)
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "成功",
		Data:    nil,
	})
}

func Confirm(ctx echo.Context) error {
	id := ctx.Param("id")
	u := ctx.Get("user_id").(int32)
	var o store.Order
	o.OrderNo = id
	err := o.Get()
	if err != nil {
		return handler.ErrorResp(ctx, fmt.Errorf("未找到该订单  "), 404)
	}
	if o.Buyer != u {
		return handler.ErrorResp(ctx, fmt.Errorf("此订单非你所拥有 "), 403)
	}
	if o.Status != SHIPPING {
		return handler.ErrorResp(ctx, fmt.Errorf("你tm什么问题? "), 403)
	}
	err = o.Update("status", RECEIVED)
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "成功",
		Data:    nil,
	})
}
