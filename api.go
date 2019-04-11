package bitmex

import (
	"context"
	"github.com/chuckpreslar/emission"
	"github.com/mariuspass/recws"
	"github.com/vmpartner/bitmex/rest"
	"sync"
)

const (
	BitmexHost        = "www.bitmex.com"
	BitmexTestnetHost = "testnet.bitmex.com"
)

type RateLimit struct {
	Limit     int64
	Remaining int64
	Reset     int64
}

// BitMEX describes the API
type BitMEX struct {
	Key    string
	Secret string
	symbol string
	host   string

	ctx                  context.Context
	rateLimitMutexPublic sync.RWMutex
	rateLimitMutex       sync.RWMutex
	rateLimitPublic      RateLimit
	rateLimit            RateLimit

	ws             recws.RecConn
	emitter        *emission.Emitter
	subscribeCmd   *WSCmd
	orderBook      *OrderBookLocal
	snapshotLoaded map[string]bool // key: symbol
}

// New allows the use of the public or private and websocket api
func New(host string, key string, secret string, symbol string) *BitMEX {
	b := &BitMEX{}
	b.Key = key
	b.Secret = secret
	b.symbol = symbol
	b.emitter = emission.NewEmitter()
	b.orderBook = NewOrderBookLocal()
	b.snapshotLoaded = make(map[string]bool)
	b.ws = recws.RecConn{
		SubscribeHandler: b.subscribeHandler,
	}
	b.host = host
	b.ctx = rest.MakeContext(key, secret, host, 10)
	return b
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
