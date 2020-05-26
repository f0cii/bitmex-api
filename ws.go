package bitmex

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/frankrap/bitmex-api/swagger"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"log"
	"net/url"
	"time"
)

const (
	// Bitmex websocket op
	BitmexWSAnnouncement        = "announcement"        // 网站公告
	BitmexWSChat                = "chat"                // Trollbox 聊天室
	BitmexWSConnected           = "connected"           // 已连接用户/机器人的统计数据
	BitmexWSFunding             = "funding"             // 永续产品的资金费率更新 每个资金时段发送（通常是 8 小时）
	BitmexWSInstrument          = "instrument"          // 产品更新，包括交易量以及报价
	BitmexWSInsurance           = "insurance"           // 每日保险基金的更新
	BitmexWSLiquidation         = "liquidation"         // 进入委托列表的强平委托
	BitmexWSOrderBookL2_25      = "orderBookL2_25"      // 前 25 层的 Level 2 委托列表
	BitmexWSOrderBookL2         = "orderBookL2"         // 完整的 level 2 委托列表
	BitmexWSOrderBook10         = "orderBook10"         // 前 10 层的委托列表，用传统的完整委托列表推送
	BitmexWSPublicNotifications = "publicNotifications" // 全系统的告示（用于段时间的消息）
	BitmexWSQuote               = "quote"               // 最高层的委托列表
	BitmexWSQuoteBin1m          = "quoteBin1m"          // 每分钟报价数据
	BitmexWSQuoteBin5m          = "quoteBin5m"          // 每 5 分钟报价数据
	BitmexWSQuoteBin1h          = "quoteBin1h"          // 每小时报价数据
	BitmexWSQuoteBin1d          = "quoteBin1d"          // 每天报价数据
	BitmexWSSettlement          = "settlement"          // 结算信息
	BitmexWSTrade               = "trade"               // 实时交易
	BitmexWSTradeBin1m          = "tradeBin1m"          // 每分钟交易数据
	BitmexWSTradeBin5m          = "tradeBin5m"          // 每 5 分钟交易数据
	BitmexWSTradeBin1h          = "tradeBin1h"          // 每小时交易数据
	BitmexWSTradeBin1d          = "tradeBin1d"          // 每天交易数据

	// Bitmex websocket private op
	BitmexWSAffiliate            = "affiliate"            // 邀请人状态，已邀请用户及分红比率
	BitmexWSExecution            = "execution"            // 个别成交，可能是多个成交
	BitmexWSOrder                = "order"                // 你委托的更新
	BitmexWSMargin               = "margin"               // 你账户的余额和保证金要求的更新
	BitmexWSPosition             = "position"             // 你仓位的更新
	BitmexWSPrivateNotifications = "privateNotifications" // 个人的通知，现时并未使用
	BitmexWSTransact             = "transact"             // 资金提存更新
	BitmexWSWallet               = "wallet"               // 比特币余额更新及总提款存款

	bitmexActionInitialData = "partial"
	bitmexActionInsertData  = "insert"
	bitmexActionDeleteData  = "delete"
	bitmexActionUpdateData  = "update"
)

type SubscribeInfo struct {
	Op    string
	Param string
}

type WSCmd struct {
	Command string        `json:"op"`
	Args    []interface{} `json:"args"`
}

type Response struct {
	Success   bool        `json:"success,omitempty"`
	Subscribe string      `json:"subscribe,omitempty"`
	Request   interface{} `json:"request,omitempty"`
	Table     string      `json:"table,omitempty"`
	Action    string      `json:"action,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func decodeMessage(message []byte) (Response, error) {
	var res Response
	err := json.Unmarshal(message, &res)
	if err != nil {
		return res, err
	}

	ret := gjson.ParseBytes(message)

	if ret.Get("table").Exists() {
		raw := ret.Get("data").Raw
		switch res.Table {
		case BitmexWSInstrument:
			var instruments []*swagger.Instrument
			err = json.Unmarshal([]byte(raw), &instruments)
			if err != nil {
				return res, err
			}
			res.Data = instruments
		case BitmexWSOrderBookL2:
			var orderbooks OrderBookData
			err = json.Unmarshal([]byte(raw), &orderbooks)
			if err != nil {
				return res, err
			}
			res.Data = orderbooks
		case BitmexWSQuote:
			var quotes []*swagger.Quote
			err = json.Unmarshal([]byte(raw), &quotes)
			if err != nil {
				return res, err
			}
			res.Data = quotes
		case BitmexWSTradeBin1m, BitmexWSTradeBin5m, BitmexWSTradeBin1h, BitmexWSTradeBin1d:
			var tradeBins []*swagger.TradeBin
			err = json.Unmarshal([]byte(raw), &tradeBins)
			if err != nil {
				return res, err
			}
			res.Data = tradeBins
		case BitmexWSTrade:
			var trades []*swagger.Trade
			err = json.Unmarshal([]byte(raw), &trades)
			if err != nil {
				return res, err
			}
			res.Data = trades
		case BitmexWSExecution:
			var executions []*swagger.Execution
			err = json.Unmarshal([]byte(raw), &executions)
			if err != nil {
				return res, err
			}
			res.Data = executions
		case BitmexWSOrder:
			var orders []*swagger.Order
			err = json.Unmarshal([]byte(raw), &orders)
			if err != nil {
				return res, err
			}
			res.Data = orders
		case BitmexWSMargin:
			var margins []*swagger.Margin
			err = json.Unmarshal([]byte(raw), &margins)
			if err != nil {
				return res, err
			}
			res.Data = margins
		case BitmexWSPosition:
			var positions []*swagger.Position
			err = json.Unmarshal([]byte(raw), &positions)
			if err != nil {
				return res, err
			}
			res.Data = positions
		case BitmexWSWallet:
			var wallets []*swagger.Wallet
			err = json.Unmarshal([]byte(raw), &wallets)
			if err != nil {
				return res, err
			}
			res.Data = wallets
		}
	}
	return res, err
}

func (b *BitMEX) sendWSMessage(msg interface{}) error {
	msgs, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "marshalling WSmessage failed")
	}
	log.Println(string(msgs))

	err = b.ws.WriteMessage(websocket.TextMessage, msgs)
	if err != nil {
		return errors.Wrap(err, "sending WSmessage failed")
	}
	return nil
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
		args := v.Op
		if v.Param != "" {
			args += ":" + v.Param
		}
		message.Args = append(message.Args, args) // v.Op+":"+v.Param "quote:XBTUSD"
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

// StartWS opens the websocket connection, and waits for message events
func (b *BitMEX) StartWS() {
	u := url.URL{Scheme: "wss", Host: b.host, Path: "/realtime"}
	bitmexWSURL := u.String()
	b.ws.SetProxyURL(b.proxyURL)
	b.ws.Dial(bitmexWSURL, nil)

	go func() {
		t := time.NewTicker(time.Second * 5)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				err := b.ws.WriteMessage(websocket.TextMessage, []byte("ping"))
				if err != nil {
					// The connection has disconnected if ping errors
					// and everything will automatically tear down.
					log.Printf("%v", err)
				}
			}
		}
	}()

	go func() {
		for {
			messageType, message, err := b.ws.ReadMessage()
			if err != nil {
				if b.ws.IsClosed() {
					log.Println("StartWS done")
					return
				}
				time.Sleep(500 * time.Millisecond)
				log.Println("read:", err)
				continue
			}
			if messageType == websocket.TextMessage {
				if string(message) == "pong" {
					continue
				}
			}
			resp, err := decodeMessage(message)
			if err != nil {
				log.Println("decode:", err)
				continue
			}

			if resp.Success {
				if b.debugMode {
					log.Println(string(message))
				}
				continue
			}

			switch resp.Table {
			case BitmexWSInstrument:
				b.processInstrument(&resp)
			case BitmexWSOrderBookL2_25:
				b.processOrderbook(&resp)
			case BitmexWSOrderBookL2:
				b.processOrderbook(&resp)
			case BitmexWSQuote:
				b.processQuote(&resp)
			case BitmexWSTradeBin1m, BitmexWSTradeBin5m, BitmexWSTradeBin1h, BitmexWSTradeBin1d:
				b.processTradeBin(&resp, resp.Table)
			case BitmexWSTrade:
				b.processTrade(&resp)
			case BitmexWSExecution:
				b.processExecution(&resp)
			case BitmexWSOrder:
				b.processOrder(&resp)
			case BitmexWSMargin:
				b.processMargin(&resp)
			case BitmexWSPosition:
				b.processPosition(&resp)
			case BitmexWSWallet:
				b.processWallet(&resp)
			default:
				if resp.Subscribe != "" {
					if b.debugMode {
						log.Printf("Subscribe message Msg=%#v", resp)
					}
				} else {
					if b.debugMode {
						log.Printf("Unknown message Msg=%#v", resp)
						log.Println("resp:", string(message))
					}
				}
			}

			//log.Println("resp:", resp)
		}
	}()
}

// CloseWS closes the websocket connection
func (b *BitMEX) CloseWS() {
	b.ws.CloseWS()
}

func (b *BitMEX) processInstrument(msg *Response) (err error) {
	instruments, _ := msg.Data.([]*swagger.Instrument)
	if len(instruments) < 1 {
		return errors.New("ws.go error - no instrument data")
	}

	b.emitter.Emit(BitmexWSInstrument, instruments, msg.Action)
	return nil
}

func (b *BitMEX) processOrderbook(msg *Response) (err error) {
	orderbook, _ := msg.Data.(OrderBookData)
	if len(orderbook) < 1 {
		return errors.New("ws.go error - no orderbook data")
	}

	symbol := orderbook[0].Symbol

	_, ok := b.orderBookLoaded[symbol]
	if !ok {
		b.orderBookLoaded[symbol] = false
	}

	_, ok = b.orderBookLocals[symbol]
	if !ok {
		b.orderBookLocals[symbol] = NewOrderBookLocal()
	}

	switch msg.Action {
	case bitmexActionInitialData:
		if !b.orderBookLoaded[symbol] {
			b.orderBookLocals[symbol].LoadSnapshot(orderbook)
			b.orderBookLoaded[symbol] = true
		}
	default:
		if b.orderBookLoaded[symbol] {
			b.orderBookLocals[symbol].Update(orderbook, msg.Action)
		}
	}

	ob := b.orderBookLocals[symbol].GetOrderbookL2()
	b.emitter.Emit(BitmexWSOrderBookL2_25, ob, symbol)
	b.emitter.Emit(BitmexWSOrderBookL2, ob, symbol)
	return nil
}

func (b *BitMEX) processQuote(msg *Response) (err error) {
	quotes, _ := msg.Data.([]*swagger.Quote)
	if len(quotes) < 1 {
		return errors.New("ws.go error - no quote data")
	}

	b.emitter.Emit(BitmexWSQuote, quotes, msg.Action)
	return nil
}

func (b *BitMEX) processTradeBin(msg *Response, name string) (err error) {
	tradeBins, _ := msg.Data.([]*swagger.TradeBin)
	if len(tradeBins) < 1 {
		return errors.New("ws.go error - no tradeBin data")
	}

	b.emitter.Emit(name, tradeBins, msg.Action)
	return nil
}

func (b *BitMEX) processTrade(msg *Response) (err error) {
	trades, _ := msg.Data.([]*swagger.Trade)
	if len(trades) < 1 {
		return errors.New("ws.go error - no trade data")
	}
	b.emitter.Emit(BitmexWSTrade, trades, msg.Action)
	return nil
}

func (b *BitMEX) processExecution(msg *Response) (err error) {
	executions, _ := msg.Data.([]*swagger.Execution)
	if len(executions) < 1 {
		return errors.New("ws.go error - no execution data")
	}

	b.emitter.Emit(BitmexWSExecution, executions, msg.Action)
	return nil
}

func (b *BitMEX) processOrder(msg *Response) (err error) {
	orders, _ := msg.Data.([]*swagger.Order)
	if len(orders) < 1 {
		return errors.New("ws.go error - no order data")
	}

	switch msg.Action {
	case bitmexActionInitialData, bitmexActionInsertData:
		for _, v := range orders {
			b.orderLocals[v.OrderID] = v
		}
	case bitmexActionUpdateData:
		for _, v := range orders {
			if old, ok := b.orderLocals[v.OrderID]; ok {
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
				if v.Side != "" {
					old.Side = v.Side
				}
				old.Timestamp = v.Timestamp // 2018-10-12T02:33:18.886Z
			}
		}
	case bitmexActionDeleteData:
	}

	var result []*swagger.Order
	for _, v := range orders {
		order, ok := b.orderLocals[v.OrderID]
		if ok {
			newOrder := *order
			result = append(result, &newOrder)
		}
	}

	//b.emitter.Emit(BitmexWSOrder, orders, msg.Action)
	b.emitter.Emit(BitmexWSOrder, result, msg.Action)
	return nil
}

func (b *BitMEX) processMargin(msg *Response) (err error) {
	margins, _ := msg.Data.([]*swagger.Margin)
	if len(margins) < 1 {
		return errors.New("ws.go error - no margin data")
	}

	b.emitter.Emit(BitmexWSMargin, margins, msg.Action)
	return nil
}

func (b *BitMEX) processPosition(msg *Response) (err error) {
	positions, _ := msg.Data.([]*swagger.Position)
	if len(positions) < 1 {
		return errors.New("ws.go error - no position data")
	}

	b.emitter.Emit(BitmexWSPosition, positions, msg.Action)
	return nil
}

func (b *BitMEX) processWallet(msg *Response) (err error) {
	wallets, _ := msg.Data.([]*swagger.Wallet)
	if len(wallets) < 1 {
		return errors.New("ws.go error - no wallet data")
	}

	b.emitter.Emit(BitmexWSWallet, wallets, msg.Action)
	return nil
}
