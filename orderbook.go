package bitmex

import (
	"sort"
	"strconv"
	"sync"
	"time"
)

// OrderBookL2 contains order book l2
type OrderBookL2 struct {
	ID     int64   `json:"id"`
	Price  float64 `json:"price"`
	Side   string  `json:"side"`
	Size   int64   `json:"size"`
	Symbol string  `json:"symbol"`
}

func (o *OrderBookL2) Key() string {
	// return fmt.Sprintf("%s-%s-%d", o.Symbol, o.Side, o.ID)
	return strconv.FormatInt(o.ID, 10)
}

type OrderBookData []*OrderBookL2

// OrderBook10 contains order book 10
type OrderBook10 struct {
	Bids      [][]float64 `json:"bids"`
	Asks      [][]float64 `json:"asks"`
	Timestamp time.Time   `json:"timestamp"`
	Symbol    string      `json:"symbol"`
}

// Item stores the amount and price values
type Item struct {
	Amount float64
	Price  float64
}

type OrderBook struct {
	Bids      []Item    `json:"bids"`
	Asks      []Item    `json:"asks"`
	Timestamp time.Time `json:"timestamp"`
}

func (ob *OrderBook) Valid() bool {
	return len(ob.Bids) > 0 && len(ob.Asks) > 0
}

func (ob *OrderBook) Bid() float64 {
	if len(ob.Bids) < 1 {
		return 0.0
	}
	return ob.Bids[0].Price
}

func (ob *OrderBook) Ask() float64 {
	if len(ob.Asks) < 1 {
		return 0.0
	}
	return ob.Asks[0].Price
}

type OrderBookLocal struct {
	ob map[string]*OrderBookL2
	m  sync.Mutex
}

func NewOrderBookLocal() *OrderBookLocal {
	o := &OrderBookLocal{
		ob: make(map[string]*OrderBookL2),
	}
	return o
}

func (o *OrderBookLocal) GetOrderbook() (ob OrderBook) {
	//ob.Symbol = "XBTUSD"
	for _, v := range o.ob {
		switch v.Side {
		case "Buy":
			ob.Bids = append(ob.Bids, Item{
				Price:  v.Price,
				Amount: float64(v.Size),
			})
		case "Sell":
			ob.Asks = append(ob.Asks, Item{
				Price:  v.Price,
				Amount: float64(v.Size),
			})
		}
	}

	sort.Slice(ob.Bids, func(i, j int) bool {
		return ob.Bids[i].Price > ob.Bids[j].Price
	})

	sort.Slice(ob.Asks, func(i, j int) bool {
		return ob.Asks[i].Price < ob.Asks[j].Price
	})

	ob.Timestamp = time.Now()

	return
}

func (o *OrderBookLocal) LoadSnapshot(newOrderbook []*OrderBookL2) error {
	o.m.Lock()
	defer o.m.Unlock()

	o.ob = make(map[string]*OrderBookL2)

	for _, v := range newOrderbook {
		o.ob[v.Key()] = v
	}

	return nil
}

func (o *OrderBookLocal) Update(orderbook []*OrderBookL2, action string) {
	o.m.Lock()
	defer o.m.Unlock()

	switch action {
	case bitmexActionUpdateData:
		for _, elem := range orderbook {
			if v, ok := o.ob[elem.Key()]; ok {
				// price is same while id is same
				// v.Price = elem.Price
				v.Size = elem.Size
				v.Side = elem.Side
			}
		}
	case bitmexActionDeleteData:
		for _, v := range orderbook {
			delete(o.ob, v.Key())
		}
	case bitmexActionInsertData:
		for _, v := range orderbook {
			o.ob[v.Key()] = v
		}
	}
}
