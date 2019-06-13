package bitmex

import (
	"fmt"
	"github.com/dgryski/dgoogauth"
	"math"
	"testing"
	"time"
)

func TestBitMEX_GetVersion(t *testing.T) {
	bitmex := newBitmexForTest()
	version, _, err := bitmex.GetVersion()
	if err != nil {
		t.Error(err)
	}

	now := time.Now().UnixNano() / 1000000
	diff := version.Timestamp - now

	t.Logf("Version time: %v now: %v diff: %v", version.Timestamp, now, diff)

	if math.Abs(float64(diff)) >= 2000 { // 时间差多余2s
		t.Error("time error")
	}
}

func TestBitMEX_GetOrderBookL2(t *testing.T) {
	bitmex := newBitmexForTest()
	orderBookL2, err := bitmex.getOrderBookL2(5, "XBTUSD")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(orderBookL2)
	rateLimit := bitmex.GetRateLimit()
	t.Logf("RateLimit: %v", rateLimit)
}

func TestBitMEX_GetOrderBook(t *testing.T) {
	bitmex := newBitmexForTest()
	orderBook, err := bitmex.GetOrderBook(5, "XBTUSD")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(orderBook)
}

func TestBitMEX_GetOrders(t *testing.T) {
	bitmex := newBitmexForTest()
	orders, err := bitmex.GetOrders("XBTUSD")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(orders)
}

func TestBitMEX_GetWallet(t *testing.T) {
	bitmex := newBitmexForTest()
	wallet, err := bitmex.GetWallet()
	if err != nil {
		t.Error(err)
		return
	}
	// XBt
	t.Logf("%#v", wallet)
}

func TestBitMEX_GetMargin(t *testing.T) {
	bitmex := newBitmexForTest()
	margin, err := bitmex.GetMargin()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", margin)
}

func TestBitMEX_GetPositions(t *testing.T) {
	bitmex := newBitmexForTest()
	positions, err := bitmex.GetPositions("XBTUSD")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(positions)
}

func TestBitMEX_PositionUpdateLeverage(t *testing.T) {
	bitmex := newBitmexForTest()
	leverage := 2.0
	position, err := bitmex.PositionUpdateLeverage(leverage, "XBTUSD")
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%#v", position)

	if position.Leverage != leverage {
		t.Error("Leverage error")
		return
	}
}

func TestBitMEX_NewOrder(t *testing.T) {
	bitmex := newBitmexForTest()
	price := 3000.0
	order, err := bitmex.NewOrder(SIDE_BUY, ORD_TYPE_LIMIT, price, 20, true, "", "XBTUSD")
	if err != nil {
		// 403 Forbidden
		t.Error(err)
		return
	}
	if order.Symbol != "XBTUSD" {
		t.Errorf("symbol error [%v]", order.Symbol)
		return
	}
	if order.Price != price {
		t.Errorf("price error [%v]", order.Price)
		return
	}
	t.Logf("%#v", order)
}

func TestBitMEX_CancelOrder(t *testing.T) {
	bitmex := newBitmexForTest()
	oid := `e4c72847-93f9-0304-d666-5f7d6ceb3ade`
	order, err := bitmex.CancelOrder(oid)
	if err != nil {
		// 400 Bad Request
		t.Error(err)
		return
	}
	t.Logf("%#v", order)
}

func TestBitMEX_CancelAllOrders(t *testing.T) {
	bitmex := newBitmexForTest()
	orders, err := bitmex.CancelAllOrders("XBTUSD")
	if err != nil {
		// 400 Bad Request
		t.Error(err)
		return
	}
	t.Logf("%#v", orders)
}

func TestBitMEX_AmendOrder(t *testing.T) {
	bitmex := newBitmexForTest()
	oid := `a17d25a3-6149-3edf-d196-75cc775beb29`
	newPrice := 3001.0
	order, err := bitmex.AmendOrder(oid, newPrice)
	if err != nil {
		t.Error(err)
		return
	}
	if order.Price != newPrice {
		t.Error("price error")
		return
	}
	t.Logf("%#v", order)
}

func TestBitMEX_CloseOrder(t *testing.T) {
	bitmex := newBitmexForTest()
	price := 6000.0
	order, err := bitmex.CloseOrder(SIDE_SELL, ORD_TYPE_LIMIT, price, 20, true, "", "XBTUSD")
	if err != nil {
		// 403 Forbidden
		t.Error(err)
		return
	}
	if order.Symbol != "XBTUSD" {
		t.Errorf("symbol error [%v]", order.Symbol)
		return
	}
	if order.Price != price {
		t.Errorf("price error [%v]", order.Price)
		return
	}
	t.Logf("%#v", order)
}

func TestBitMEX_RequestWithdrawal(t *testing.T) {
	bitmex := newBitmexForTest2()
	currency := "XBt"
	amount := float32(1000000.0)
	address := "3BMEXT9XuBTSkALWTovH1idLSC2tusjKBT"
	optToken := ""
	fee := 0.0
	// 401 Unauthorized
	trans, err := bitmex.RequestWithdrawal(currency, amount, address, optToken, fee)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", trans)
}
