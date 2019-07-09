package main

import (
	"fmt"
	"github.com/sumorf/bitmex-api"
	"github.com/sumorf/bitmex-api/swagger"
	"log"
)

func main() {
	b := bitmex.New(bitmex.HostTestnet, "8K2Oi0bnRRZ7GK4UJnY-38oj", "9EmGvk8mKX5nWa11y1KyPPGn78Lv2ZEiLx3TH0YasE_oE06y")
	subscribeInfos := []bitmex.SubscribeInfo{
		{Op: bitmex.BitmexWSOrderBookL2, Param: "XBTUSD"},
		{Op: bitmex.BitmexWSOrder, Param: "XBTUSD"},
		{Op: bitmex.BitmexWSPosition, Param: "XBTUSD"},
		{Op: bitmex.BitmexWSWallet, Param: "XBTUSD"},
	}
	err := b.Subscribe(subscribeInfos)
	if err != nil {
		log.Fatal(err)
	}

	b.On(bitmex.BitmexWSOrderBookL2, func(m bitmex.OrderBookDataL2, symbol string) {
		ob := m.OrderBook()
		fmt.Printf("\rOrderbook Asks: %#v Bids: %#v                            ", ob.Asks[0], ob.Bids[0])
	}).On(bitmex.BitmexWSOrder, func(m []*swagger.Order, action string) {
		fmt.Printf("Order action=%v orders=%#v\n", action, m)
		//for _, v := range m {
		//	fmt.Printf("Order=%#v %v\n", *v, (*v).Timestamp.String())
		//}
	}).On(bitmex.BitmexWSPosition, func(m []*swagger.Position, action string) {
		fmt.Printf("Position action=%v positions=%#v\n", action, m)
	}).On(bitmex.BitmexWSWallet, func(m []*swagger.Wallet, action string) {
		fmt.Printf("Wallet action=%v wallets=%#v\n", action, m)
	})

	b.StartWS()

	select {}
}
