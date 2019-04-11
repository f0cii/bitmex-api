package bitmex

import (
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
	orderBookL2, err := bitmex.GetOrderBookL2(5)
	if err != nil {
		t.Error(err)
	}
	t.Log(orderBookL2)
	rateLimit := bitmex.GetRateLimit()
	t.Logf("RateLimit: %v", rateLimit)
}

func TestBitMEX_GetOrders(t *testing.T) {
	bitmex := newBitmexForTest()
	orders, err := bitmex.GetOrders()
	if err != nil {
		t.Error(err)
	}
	t.Log(orders)
}

func TestBitMEX_GetPositions(t *testing.T) {
	bitmex := newBitmexForTest()
	positions, err := bitmex.GetPositions()
	if err != nil {
		t.Error(err)
	}
	t.Log(positions)
}
