package bitmex

import (
	"context"
	"github.com/chuckpreslar/emission"
	"github.com/sumorf/bitmex-api/recws"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"time"

	//"github.com/mariuspass/recws"
	"github.com/sumorf/bitmex-api/swagger"
	"net/http"
	"net/url"
	"sync"
)

const (
	HostReal    = "www.bitmex.com"
	HostTestnet = "testnet.bitmex.com"
)

type RateLimit struct {
	Limit     int64
	Remaining int64
	Reset     int64
}

// BitMEX describes the API
type BitMEX struct {
	Key      string
	Secret   string
	host     string
	proxyURL string

	ctx                  context.Context
	timeout              time.Duration
	httpClient           *http.Client
	cfg                  *swagger.Configuration
	client               *swagger.APIClient
	rateLimitMutexPublic sync.RWMutex
	rateLimitMutex       sync.RWMutex
	rateLimitPublic      RateLimit
	rateLimit            RateLimit

	ws              recws.RecConn
	emitter         *emission.Emitter
	subscribeCmd    *WSCmd
	orderBookLocals map[string]*OrderBookLocal // key: symbol
	orderLocals     map[string]*swagger.Order  // key: OrderID
	orderBookLoaded map[string]bool            // key: symbol
}

// New allows the use of the public or private and websocket api
func New(host string, key string, secret string) *BitMEX {
	b := &BitMEX{}
	b.Key = key
	b.Secret = secret
	b.emitter = emission.NewEmitter()
	b.orderBookLocals = make(map[string]*OrderBookLocal)
	b.orderLocals = make(map[string]*swagger.Order)
	b.orderBookLoaded = make(map[string]bool)
	b.ws = recws.RecConn{
		SubscribeHandler: b.subscribeHandler,
	}
	b.host = host
	b.ctx = MakeContext(key, secret, host, 10)
	b.timeout = 10 * time.Second
	b.cfg = GetConfiguration(b.ctx)
	b.httpClient = &http.Client{
		Timeout: b.timeout,
	}
	b.cfg.HTTPClient = b.httpClient
	b.client = swagger.NewAPIClient(b.cfg)
	return b
}

// SetHttpProxy proxyURL: http://127.0.0.1:1080
func (b *BitMEX) SetHttpProxy(proxyURL string) error {
	proxyURL_, err := url.Parse(proxyURL)
	if err != nil {
		return err
	}

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL_),
	}

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
		Timeout:   b.timeout,
	}
	b.httpClient = client
	b.cfg.HTTPClient = client
	b.proxyURL = proxyURL
	return nil
}

// SetProxy proxyURL: 127.0.0.1:1080
func (b *BitMEX) SetProxy(socks5Proxy string) error {
	// socks5Proxy := "127.0.0.1:1080"
	dialer, err := proxy.SOCKS5("tcp", socks5Proxy, nil, proxy.Direct)
	if err != nil {
		log.Fatal("Error creating dialer, aborting.")
	}

	dialFunc := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.Dial(network, addr)
	}

	tr := &http.Transport{DialContext: dialFunc} // Dial: dialer.Dial,
	client := &http.Client{Transport: tr, Timeout: b.timeout}

	b.httpClient = client
	b.cfg.HTTPClient = client
	b.proxyURL = "http://" + socks5Proxy
	return nil
}

func (b *BitMEX) GetRateLimit() RateLimit {
	b.rateLimitMutex.RLock()
	defer b.rateLimitMutex.RUnlock()
	return b.rateLimit
}

func (b *BitMEX) GetRateLimitPublic() RateLimit {
	b.rateLimitMutexPublic.RLock()
	defer b.rateLimitMutexPublic.RUnlock()
	return b.rateLimitPublic
}

func MakeContext(key string, secret string, host string, timeout int64) context.Context {
	return context.WithValue(context.TODO(), swagger.ContextAPIKey, swagger.APIKey{
		Key:     key,
		Secret:  secret,
		Host:    host,
		Timeout: timeout,
	})
}

func GetConfiguration(ctx context.Context) *swagger.Configuration {
	c := ctx.Value(swagger.ContextAPIKey).(swagger.APIKey)
	cfg := &swagger.Configuration{
		BasePath:      "https://" + c.Host + "/api/v1",
		DefaultHeader: make(map[string]string),
		UserAgent:     "Swagger-Codegen/1.0.0/go",
		ExpireTime:    5, //seconds
	}
	return cfg
}

func GetClient(ctx context.Context) *swagger.APIClient {
	c := ctx.Value(swagger.ContextAPIKey).(swagger.APIKey)
	cfg := &swagger.Configuration{
		BasePath:      "https://" + c.Host + "/api/v1",
		DefaultHeader: make(map[string]string),
		UserAgent:     "Swagger-Codegen/1.0.0/go",
		ExpireTime:    5, //seconds
	}

	return swagger.NewAPIClient(cfg)
}
