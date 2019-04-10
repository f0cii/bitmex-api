package bitmex

import "os"

const (
	proxyURL = "http://127.0.0.1:1080"
)

func useProxy() {
	os.Setenv("HTTP_PROXY", proxyURL)
	os.Setenv("HTTPS_PROXY", proxyURL)
}
