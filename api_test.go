package bitmex

import (
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"testing"
)

var (
	_proxyURL string
	_testnet  bool
	_key      string
	_secret   string
	_key2      string
	_secret2   string
)

func loadConfig() {
	err := config.Load(file.NewSource(file.WithPath("./testdata/config.yaml")))
	if err != nil {
		return
	}
	cfg := config.Map()

	_proxyURL = cfg["proxy_url"].(string)
	_testnet = cfg["testnet"].(bool)
	_key = cfg["key"].(string)
	_secret = cfg["secret"].(string)
	_key2 = cfg["key2"].(string)
	_secret2 = cfg["secret2"].(string)
}

func TestConfig(t *testing.T) {
	err := config.Load(file.NewSource(file.WithPath("./testdata/config.yaml")))
	if err != nil {
		t.Error(err)
	}
	cfg := config.Map()

	proxy := cfg["proxy_url"].(string)
	testnet := cfg["testnet"].(bool)
	key := cfg["key"].(string)
	secret := cfg["secret"].(string)
	key2 := cfg["key2"].(string)
	secret2 := cfg["secret2"].(string)

	t.Log(proxy)
	t.Log(testnet)
	t.Log(key)
	t.Log(secret)
	t.Log(key2)
	t.Log(secret2)
}

func newBitmexForTest() *BitMEX {
	loadConfig()
	var host string
	if _testnet {
		host = HostTestnet
	} else {
		host = HostReal
	}
	bitmex := New(host, _key, _secret)
	if _proxyURL != "" {
		bitmex.SetHttpProxy(_proxyURL)
	}
	return bitmex
}

func newBitmexForTest2() *BitMEX {
	loadConfig()
	var host string
	if _testnet {
		host = HostTestnet
	} else {
		host = HostReal
	}
	bitmex := New(host, _key2, _secret2)
	if _proxyURL != "" {
		bitmex.SetHttpProxy(_proxyURL)
	}
	return bitmex
}
