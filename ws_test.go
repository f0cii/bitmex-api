package bitmex

import (
	"fmt"
	"log"
	"testing"
)

func TestBitMEXWS(t *testing.T) {
	b := New(BitmexTestnetHost, "", "", "XBTUSD")
	subscribeInfos := []SubscribeInfo{
		{Op: BitmexWSOrderBookL2, Param: "XBTUSD"},
	}
	err := b.Subscribe(subscribeInfos)
	if err != nil {
		log.Fatal(err)
	}

	b.On(BitmexWSOrderBookL2, func(m OrderBook) {
		fmt.Printf("\rOrderbook Asks: %#v Bids: %#v                            ", m.Asks[0], m.Bids[0])
	})

	b.StartWS()

	select {}
}
