package bitmex

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/chuckpreslar/emission"
	"github.com/mariuspass/recws"
	"log"
	"net/url"
	"time"
)

const (
	BitmexHost        = "www.bitmex.com"
	BitmexTestnetHost = "testnet.bitmex.com"
)

// BitMEX describes the API
type BitMEX struct {
	Key            string
	Secret         string
	symbol         string
	ws             recws.RecConn
	emitter        *emission.Emitter
	subscribeCmd   *WSCmd
	orderBook      *OrderBookLocal
	snapshotLoaded map[string]bool // key: symbol
}

// New allows the use of the public or private and websocket api
func New(host string, key string, secret string) *BitMEX {
	b := &BitMEX{}
	b.Key = key
	b.Secret = secret
	b.symbol = "XBTUSD"
	b.emitter = emission.NewEmitter()
	b.orderBook = NewOrderBookLocal()
	b.snapshotLoaded = make(map[string]bool)
	b.ws = recws.RecConn{
		SubscribeHandler: b.subscribeHandler,
	}
	u := url.URL{Scheme: "wss", Host: host, Path: "/realtime"}
	bitmexWSURL := u.String()
	b.ws.Dial(bitmexWSURL, nil)
	return b
}

// sendAuth sends an authenticated subscription
func (b *BitMEX) sendAuth() error {
	if b.Key == "" || b.Secret == "" {
		return nil
	}
	msg := b.getAuthMessage(b.Key, b.Secret)
	log.Println("sendAuth")
	return b.sendWSMessage(msg)
}

func (b *BitMEX) getAuthMessage(key string, secret string) WSCmd {
	nonce := time.Now().Unix() + 412
	req := fmt.Sprintf("GET/realtime%d", nonce)
	sig := hmac.New(sha256.New, []byte(secret))
	sig.Write([]byte(req))
	signature := hex.EncodeToString(sig.Sum(nil))
	var msgKey []interface{}
	msgKey = append(msgKey, key)
	msgKey = append(msgKey, nonce)
	msgKey = append(msgKey, signature)

	return WSCmd{"authKey", msgKey}
}

func (b *BitMEX) Subscribe(subscribeTypes []SubscribeInfo) error {
	message := WSCmd{}
	message.Command = "subscribe"
	for _, v := range subscribeTypes {
		message.Args = append(message.Args, v.Op+":"+v.Param) // "quote:XBTUSD"
	}
	b.subscribeCmd = &message
	b.subscribeHandler()
	return nil
}

func (b *BitMEX) subscribeHandler() error {
	if b.subscribeCmd == nil {
		return nil
	}
	err := b.sendAuth()
	if err != nil {
		return err
	}
	log.Printf("subscribe %v", *b.subscribeCmd)
	return b.sendWSMessage(*b.subscribeCmd)
}
