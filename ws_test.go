package bitmex

import (
	"fmt"
	"log"
	"testing"
)

func testBitMeX() *BitMEX {
	b := New(nil, HostTestnet, "", "", true)
	//b.SetHttpProxy("http://127.0.0.1:1081")
	return b
}

func TestBitMEXConnect(t *testing.T) {
	b := testBitMeX()
	subscribeInfos := []SubscribeInfo{
		{Op: BitmexWSOrderBookL2, Param: "XBTUSD"},
	}
	err := b.Subscribe(subscribeInfos)
	if err != nil {
		log.Fatal(err)
	}

	b.On(BitmexWSOrderBookL2, func(ob OrderBookDataL2, symbol string) {
		m := ob.OrderBook()
		fmt.Printf("\rOrderbook Asks: %#v Bids: %#v                            ", m.Asks[0], m.Bids[0])
	})

	b.StartWS()

	select {}
}

func TestBitMEXWS(t *testing.T) {
	b := testBitMeX()
	subscribeInfos := []SubscribeInfo{
		{Op: BitmexWSOrderBookL2, Param: "XBTUSD"},
	}
	err := b.Subscribe(subscribeInfos)
	if err != nil {
		log.Fatal(err)
	}

	b.On(BitmexWSOrderBookL2, func(ob OrderBookDataL2, symbol string) {
		m := ob.OrderBook()
		fmt.Printf("\rOrderbook Asks: %#v Bids: %#v                            ", m.Asks[0], m.Bids[0])
	})

	b.StartWS()

	select {}
}

func TestBitMEX_Subscribe(t *testing.T) {
	b := testBitMeX()
	subscribeInfos := []SubscribeInfo{
		{Op: BitmexWSOrderBookL2_25, Param: "XBTUSD"},
	}
	err := b.Subscribe(subscribeInfos)
	if err != nil {
		log.Fatal(err)
	}

	b.On(BitmexWSOrderBookL2_25, func(ob OrderBookDataL2, symbol string) {
		m := ob.OrderBook()
		fmt.Printf("\rOrderbook Asks: %#v Bids: %#v                            ", m.Asks[0], m.Bids[0])
	})

	b.StartWS()

	select {}
}
