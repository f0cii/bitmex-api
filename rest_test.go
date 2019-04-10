package bitmex

import (
	"testing"
)

func TestBitMEX_GetOrderBookL2(t *testing.T) {
	useProxy()

	bitmex := New(BitmexTestnetHost, "8K2Oi0bnRRZ7GK4UJnY-38oj", "9EmGvk8mKX5nWa11y1KyPPGn78Lv2ZEiLx3TH0YasE_oE06y", "XBTUSD")
	orderBookL2, err := bitmex.GetOrderBookL2(5)
	if err != nil {
		t.Error(err)
	}
	t.Log(orderBookL2)
	rateLimit := bitmex.GetRateLimit()
	t.Logf("RateLimit: %v", rateLimit)
}
