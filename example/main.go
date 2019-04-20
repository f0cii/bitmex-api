package main

import (
	"fmt"
	"github.com/sumorf/bitmex-api"
	"github.com/sumorf/bitmex-api/swagger"
	"log"
)

func main() {
	b := bitmex.New(bitmex.BitmexTestnetHost, "8K2Oi0bnRRZ7GK4UJnY-38oj", "9EmGvk8mKX5nWa11y1KyPPGn78Lv2ZEiLx3TH0YasE_oE06y", "XBTUSD")
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

	orders := map[string]*swagger.Order{}

	b.On(bitmex.BitmexWSOrderBookL2, func(m bitmex.OrderBook) {
		fmt.Printf("\rOrderbook Asks: %#v Bids: %#v                            ", m.Asks[0], m.Bids[0])
	}).On(bitmex.BitmexWSOrder, func(m []*swagger.Order, action string) {
		fmt.Printf("Order action=%v orders=%#v\n", action, m)

		switch action {
		case "partial":
			for _, v := range m {
				orders[v.OrderID] = &(*v)
			}
		case "insert":
			for _, v := range m {
				orders[v.OrderID] = &(*v)
			}
		case "update":
			for _, v := range m {
				if old, ok := orders[v.OrderID]; ok {
					if v.Price > 0.0 {
						old.Price = v.Price
					}
					if v.OrderQty > 0 {
						old.OrderQty = v.OrderQty
					}
					if v.OrdStatus != "" {
						old.OrdStatus = v.OrdStatus
					}
					if v.AvgPx > 0 {
						old.AvgPx = v.AvgPx
					}
					if v.CumQty > 0 {
						old.CumQty = v.CumQty
					}
					if v.SimpleCumQty > 0 {
						old.SimpleCumQty = v.SimpleCumQty
					}
					if v.Text != "" {
						old.Text = v.Text
					}
					old.Timestamp = v.Timestamp // 2018-10-12T02:33:18.886Z
				}
			}
		case "delete":
			log.Printf("delete?")
		}
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
