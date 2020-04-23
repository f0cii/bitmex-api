# bitmex-api
A Go wrapper to the BitMEX API

### Example

```
package main

import (
	"fmt"
	"github.com/frankrap/bitmex-api"
	"github.com/frankrap/bitmex-api/swagger"
	"log"
)

func main() {
    b := bitmex.New(nil,
        bitmex.HostTestnet, "8K2Oi0bnRRZ7GK4UJnY-38oj", "9EmGvk8mKX5nWa11y1KyPPGn78Lv2ZEiLx3TH0YasE_oE06y", true)
	subscribeInfos := []bitmex.SubscribeInfo{
		{Op: bitmex.BitmexWSOrderBookL2, Param: "XBTUSD"},
		{Op: bitmex.BitmexWSOrder, Param: "XBTUSD"},
		{Op: bitmex.BitmexWSPosition, Param: "XBTUSD"},
		{Op: bitmex.BitmexWSMargin, Param: "XBTUSD"},
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
	}).On(bitmex.BitmexWSPosition, func(m []*swagger.Position, action string) {
		fmt.Printf("Position action=%v positions=%#v\n", action, m)
	}).On(bitmex.BitmexWSMargin, func(m []*swagger.Margin, action string) {
		fmt.Printf("Wallet action=%v margins=%#v\n", action, m)
	})

	b.StartWS()

	// Get orderbook by rest api
	b.GetOrderBook(10, "XBTUSD")
	// Place a limit buy order
	b.PlaceOrder(bitmex.SIDE_BUY, bitmex.ORD_TYPE_LIMIT, 0, 6000.0, 1000, "", "", "XBTUSD")
	b.GetOrders("XBTUSD")
	b.GetOrder("{OrderID}", "XBTUSD")
	b.AmendOrder("{OrderID}", 6000.5)
	b.CancelOrder("{OrderID}")
	b.CancelAllOrders("XBTUSD")
	b.GetPosition("XBTUSD")
	b.GetMargin()

	forever := make(chan bool)
	<-forever
}
```