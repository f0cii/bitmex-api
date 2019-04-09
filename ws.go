package bitmex

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"log"
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
		case BitmexWSOrderBookL2:
			var orderbooks OrderBookData
			err = json.Unmarshal([]byte(raw), &orderbooks)
			if err != nil {
				return res, err
			}
			res.Data = orderbooks
		case BitmexWSOrder:
			var orders []*Order
			err = json.Unmarshal([]byte(raw), &orders)
			if err != nil {
				return res, err
			}
			res.Data = orders
		case BitmexWSPosition:
			var wallets []*Wallet
			err = json.Unmarshal([]byte(raw), &wallets)
			if err != nil {
				return res, err
			}
			res.Data = wallets
		case BitmexWSWallet:
			log.Printf("Wallet: %v", raw)
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

// StartWS opens the websocket connection, and waits for message events
func (b *BitMEX) StartWS() {
	go func() {
		for {
			_, message, err := b.ws.ReadMessage()
			if err != nil {
				time.Sleep(100 * time.Millisecond)
				log.Println("read:", err)
				continue
			}
			resp, err := decodeMessage(message)
			if err != nil {
				log.Println("decode:", err)
				continue
			}

			switch resp.Table {
			case BitmexWSOrderBookL2:
				b.processOrderbook(&resp, b.symbol)
			case BitmexWSOrder:
				b.processOrder(&resp, b.symbol)
			case BitmexWSPosition:
				b.processPosition(&resp, b.symbol)
			case BitmexWSWallet:
				b.processWallet(&resp, b.symbol)
			default:
				log.Printf("Unknown message Msg=%#v", resp)
			}

			//log.Println("resp:", resp)
		}
	}()
}

func (b *BitMEX) processOrderbook(msg *Response, symbol string) (err error) {
	orderbook, _ := msg.Data.(OrderBookData)
	if len(orderbook) < 1 {
		return errors.New("ws.go error - no orderbook data")
	}

	_, ok := b.snapshotLoaded[symbol]
	if !ok {
		b.snapshotLoaded[symbol] = false
	}

	switch msg.Action {
	case bitmexActionInitialData:
		if !b.snapshotLoaded[symbol] {
			b.orderBook.LoadSnapshot(orderbook)
			b.snapshotLoaded[symbol] = true
		}
	default:
		if b.snapshotLoaded[symbol] {
			b.orderBook.Update(orderbook, msg.Action)
		}
	}

	b.emitter.Emit(BitmexWSOrderBookL2, b.orderBook.GetOrderbook())
	return nil
}

func (b *BitMEX) processOrder(msg *Response, symbol string) (err error) {
	orders, _ := msg.Data.([]*Order)
	if len(orders) < 1 {
		return errors.New("ws.go error - no order data")
	}

	//log.Printf("processOrder Action=%v orders=%v", msg.Action, orders)

	b.emitter.Emit(BitmexWSOrder, orders, msg.Action)
	return nil
}

func (b *BitMEX) processPosition(msg *Response, symbol string) (err error) {
	positions, _ := msg.Data.([]*Position)
	if len(positions) < 1 {
		return errors.New("ws.go error - no position data")
	}

	//log.Printf("processPosition Action=%v positions=%v", msg.Action, positions)

	b.emitter.Emit(BitmexWSPosition, positions, msg.Action)
	return nil
}

func (b *BitMEX) processWallet(msg *Response, symbol string) (err error) {
	wallets, _ := msg.Data.([]*Wallet)
	if len(wallets) < 1 {
		return errors.New("ws.go error - no wallet data")
	}

	//log.Printf("processWallet Action=%v wallets=%v", msg.Action, wallets)

	b.emitter.Emit(BitmexWSWallet, wallets, msg.Action)
	return nil
}
