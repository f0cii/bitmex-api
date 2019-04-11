package bitmex

import "os"

const (
	proxyURL = "http://127.0.0.1:1080"
)

func useProxy() {
	os.Setenv("HTTP_PROXY", proxyURL)
	os.Setenv("HTTPS_PROXY", proxyURL)
}

func newBitmexForTest() *BitMEX {
	useProxy()

	key := `8K2Oi0bnRRZ7GK4UJnY-38oj`
	secret := `9EmGvk8mKX5nWa11y1KyPPGn78Lv2ZEiLx3TH0YasE_oE06y`
	bitmex := New(BitmexTestnetHost, key, secret, "XBTUSD")
	return bitmex
}
