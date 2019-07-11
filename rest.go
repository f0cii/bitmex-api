package bitmex

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sumorf/bitmex-api/swagger"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"
)

const (
	SIDE_BUY  = "Buy"
	SIDE_SELL = "Sell"

	ORD_TYPE_MARKET                         = "Market"          // 市价
	ORD_TYPE_LIMIT                          = "Limit"           // 限价
	ORD_TYPE_STOP                           = "Stop"            // 市价止损
	ORD_TYPE_STOP_LIMIT                     = "StopLimit"       // 限价止损
	ORD_TYPE_MARKET_IF_TOUCHED              = "MarketIfTouched" // 市价止盈
	ORD_TYPE_LIMIT_IF_TOUCHED               = "LimitIfTouched"  // 限价止盈
	ORD_TYPE_MARKET_WITH_LEFT_OVER_AS_LIMIT = "MarketWithLeftOverAsLimit"
)

var (
	NotFound = errors.New("not found")
)

// {"name":"BitMEX API","version":"1.2.0","timestamp":1554709447283}
type Version struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
}

func (b *BitMEX) GetVersion() (version Version, time time.Duration, err error) {
	url := "https://" + b.host + "/api/v1"
	var resp *http.Response
	resp, err = b.httpClient.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	//time = resp.Time()
	err = json.Unmarshal(body, &version)
	return
}

func (b *BitMEX) GetWallet() (wallet swagger.Wallet, err error) {
	var response *http.Response

	params := map[string]interface{}{
		"currency": "",
	}
	wallet, response, err = b.client.UserApi.UserGetWallet(b.ctx, params)
	if err != nil {
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) GetMargin() (margin swagger.Margin, err error) {
	var response *http.Response

	params := map[string]interface{}{
		//"currency": "XBt",
	}
	margin, response, err = b.client.UserApi.UserGetMargin(b.ctx, params)
	if err != nil {
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) getOrderBookL2(depth int, symbol string) (orderbook []swagger.OrderBookL2, err error) {
	var response *http.Response

	params := map[string]interface{}{}
	params["depth"] = float32(depth)

	orderbook, response, err = b.client.OrderBookApi.OrderBookGetL2(symbol, params)
	if err != nil {
		return
	}
	b.onResponsePublic(response)
	return
}

func (b *BitMEX) GetOrderBook(depth int, symbol string) (ob OrderBook, err error) {
	var orderbook []swagger.OrderBookL2
	orderbook, err = b.getOrderBookL2(depth, symbol)
	if err != nil {
		return
	}
	for _, v := range orderbook {
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

func (b *BitMEX) GetBucketed(symbol string, binSize string, partial bool, filter string, columns string, count float32, start float32, reverse bool, startTime time.Time, endTime time.Time) (o []swagger.TradeBin, err error) {
	var response *http.Response

	params := map[string]interface{}{}
	params["binSize"] = binSize
	params["partial"] = partial
	params["symbol"] = symbol
	if filter != "" {
		params["filter"] = filter
	}
	if columns != "" {
		params["columns"] = columns
	}
	if count > 0 {
		params["count"] = count
	}
	if start >= 0 {
		params["start"] = start
	}
	params["reverse"] = reverse
	if !startTime.IsZero() {
		params["startTime"] = startTime
	}
	if !endTime.IsZero() {
		params["endTime"] = endTime
	}
	//params["endTime"] = endTime
	o, response, err = b.client.TradeApi.TradeGetBucketed(params)
	if err != nil {
		return
	}
	b.onResponsePublic(response)
	return
}

func (b *BitMEX) GetPosition(symbol string) (position swagger.Position, err error) {
	var positions []swagger.Position
	positions, err = b.GetPositions(symbol)
	if err != nil {
		return
	}
	if len(positions) == 0 {
		err = NotFound
		return
	}
	if len(positions) != 1 {
		err = errors.New("position error")
		return
	}
	position = positions[0]
	return
}

func (b *BitMEX) GetPositions(symbol string) (positions []swagger.Position, err error) {
	var response *http.Response

	params := map[string]interface{}{}
	if symbol != "" {
		params["filter"] = fmt.Sprintf(`{"symbol":"%s"}`, symbol)
	}

	positions, response, err = b.client.PositionApi.PositionGet(b.ctx, params)
	if err != nil {
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) GetPositionsRaw(filter string, columns string, count int32) (positions []swagger.Position, err error) {
	var response *http.Response

	params := map[string]interface{}{}
	if filter != "" {
		params["filter"] = filter
	}
	if columns != "" {
		params["columns"] = columns
	}
	if count > 0 {
		params["count"] = count
	}

	positions, response, err = b.client.PositionApi.PositionGet(b.ctx, params)
	if err != nil {
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) PositionUpdateLeverage(leverage float64, symbol string) (position swagger.Position, err error) {
	var response *http.Response
	position, response, err = b.client.PositionApi.PositionUpdateLeverage(b.ctx, symbol, leverage)
	if err != nil {
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) GetOrders(symbol string) (orders []swagger.Order, err error) {
	var response *http.Response

	params := map[string]interface{}{}
	params["symbol"] = symbol
	params["filter"] = `{"open":true}`

	orders, response, err = b.client.OrderApi.OrderGetOrders(b.ctx, params)
	if err != nil {
		return
	}
	b.onResponse(response)
	//body, _ := ioutil.ReadAll(response.Body)
	//log.Printf("%v", string(body))
	return
}

func (b *BitMEX) GetOrdersRaw(symbol string, filter string) (orders []swagger.Order, err error) {
	var response *http.Response

	params := map[string]interface{}{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	if filter != "" {
		params["filter"] = filter // `{"open":true}`
	}

	orders, response, err = b.client.OrderApi.OrderGetOrders(b.ctx, params)
	if err != nil {
		return
	}
	b.onResponse(response)
	//body, _ := ioutil.ReadAll(response.Body)
	//log.Printf("%v", string(body))
	return
}

func (b *BitMEX) NewOrder(side string, ordType string, price float64, orderQty int32, postOnly bool, timeInForce string, symbol string) (order swagger.Order, err error) {
	var response *http.Response

	params := map[string]interface{}{}
	params["symbol"] = symbol
	// params["clOrdID"] = ""	// 客户端委托ID
	params["side"] = side
	params["ordType"] = ordType
	params["orderQty"] = float32(orderQty)
	if price > 0.0 {
		params["price"] = price // Limit order only
	}
	params["text"] = `open with bitmex api`

	if timeInForce != "" { // "FillOrKill"	// 全数执行或立刻取消
		params["timeInForce"] = timeInForce
	}

	if postOnly {
		params["execInst"] = "ParticipateDoNotInitiate"
	}

	order, response, err = b.client.OrderApi.OrderNew(b.ctx, symbol, params)
	if err != nil {
		// >= 300 代表有错误
		// 400 Bad Request
		// 503
		// log.Printf("response.StatusCode: %v", response.StatusCode)
		return
	}
	b.onResponse(response)
	return
}

// PlaceOrder 放置委托单
// execInst: MarkPrice = 标记价格 IndexPrice = 指数价格 LastPrice = 最新成交 ParticipateDoNotInitiate = 被动委托
func (b *BitMEX) PlaceOrder(side string, ordType string, stopPx float64, price float64, orderQty int32, timeInForce string, execInst string, symbol string) (order swagger.Order, err error) {
	var response *http.Response

	params := map[string]interface{}{}
	params["symbol"] = symbol
	// params["clOrdID"] = ""	// 客户端委托ID
	params["side"] = side
	params["ordType"] = ordType
	params["orderQty"] = float32(orderQty)
	if stopPx > 0.0 {
		params["stopPx"] = stopPx
	}
	if price > 0.0 {
		params["price"] = price // Limit order only
	}
	params["text"] = `open with bitmex api`

	if timeInForce != "" { // "FillOrKill"	// 全数执行或立刻取消
		params["timeInForce"] = timeInForce
	}

	if execInst != "" {
		params["execInst"] = execInst
	}

	order, response, err = b.client.OrderApi.OrderNew(b.ctx, symbol, params)
	if err != nil {
		// >= 300 代表有错误
		// 400 Bad Request
		// 503
		// log.Printf("response.StatusCode: %v", response.StatusCode)
		return
	}
	b.onResponse(response)
	return
}

// PlaceOrder 放置委托单
// execInst: MarkPrice = 标记价格 IndexPrice = 指数价格 LastPrice = 最新成交 ParticipateDoNotInitiate = 被动委托
func (b *BitMEX) PlaceOrder2(clOrdID string, side string, ordType string, stopPx float64, price float64, orderQty int32,
	displayQty int32, timeInForce string, execInst string, symbol string, text string) (order swagger.Order, err error) {
	var response *http.Response

	params := map[string]interface{}{}
	if clOrdID != "" {
		params["clOrdID"] = clOrdID // 客户端委托ID
	}
	params["symbol"] = symbol
	params["side"] = side
	params["ordType"] = ordType
	params["orderQty"] = float32(orderQty)
	if displayQty >= 0 {
		params["displayQty"] = displayQty
	}
	if stopPx > 0.0 {
		params["stopPx"] = stopPx
	}
	if price > 0.0 {
		params["price"] = price // Limit order only
	}
	if text == "" {
		params["text"] = `open with bitmex api`
	} else {
		params["text"] = text
	}

	if timeInForce != "" { // "FillOrKill"	// 全数执行或立刻取消
		params["timeInForce"] = timeInForce
	}

	if execInst != "" {
		params["execInst"] = execInst
	}

	order, response, err = b.client.OrderApi.OrderNew(b.ctx, symbol, params)
	if err != nil {
		// >= 300 代表有错误
		// 400 Bad Request
		// 503
		// log.Printf("response.StatusCode: %v", response.StatusCode)
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) GetOrder(oid string, symbol string) (order swagger.Order, err error) {
	var response *http.Response
	var orders []swagger.Order

	params := map[string]interface{}{}
	params["symbol"] = symbol
	params["filter"] = fmt.Sprintf(`{"orderID":"%s"}`, oid)

	orders, response, err = b.client.OrderApi.OrderGetOrders(b.ctx, params)
	if err != nil {
		return
	}
	if len(orders) != 1 {
		err = errors.New("order error")
		return
	}
	order = orders[0]
	b.onResponse(response)
	return
}

func (b *BitMEX) GetOrderByClOrdID(clOrdID string, symbol string) (order swagger.Order, err error) {
	var response *http.Response
	var orders []swagger.Order

	params := map[string]interface{}{}
	params["symbol"] = symbol
	params["filter"] = fmt.Sprintf(`{"clOrdID":"%s"}`, clOrdID)

	orders, response, err = b.client.OrderApi.OrderGetOrders(b.ctx, params)
	if err != nil {
		return
	}
	if len(orders) != 1 {
		err = errors.New("order error")
		return
	}
	order = orders[0]
	b.onResponse(response)
	return
}

func (b *BitMEX) AmendOrder(oid string, price float64) (order swagger.Order, err error) {
	var response *http.Response

	params := map[string]interface{}{}
	params["orderID"] = oid
	params["price"] = price

	order, response, err = b.client.OrderApi.OrderAmend(b.ctx, params)
	if err != nil {
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) AmendOrder2(orderID string, origClOrdID string, clOrdID string, simpleOrderQty float64, orderQty float32, simpleLeavesQty float64, leavesQty float32, price float64, stopPx float64, pegOffsetValue float64, text string) (order swagger.Order, err error) {
	var response *http.Response

	params := map[string]interface{}{}
	if orderID != "" {
		params["orderID"] = orderID
	}
	if origClOrdID != "" {
		params["origClOrdID"] = origClOrdID
	}
	if clOrdID != "" {
		params["clOrdID"] = clOrdID
	}
	if simpleOrderQty != 0 {
		params["simpleOrderQty"] = simpleOrderQty
	}
	if orderQty != 0 {
		params["orderQty"] = orderQty
	}
	if simpleLeavesQty != 0 {
		params["simpleLeavesQty"] = simpleLeavesQty
	}
	if leavesQty != 0 {
		params["leavesQty"] = leavesQty
	}
	if price != 0 {
		params["price"] = price
	}
	if stopPx != 0 {
		params["stopPx"] = stopPx
	}
	if pegOffsetValue != 0 {
		params["pegOffsetValue"] = pegOffsetValue
	}
	if text != "" {
		params["text"] = text
	}

	order, response, err = b.client.OrderApi.OrderAmend(b.ctx, params)
	if err != nil {
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) CancelAllOrders(symbol string) (orders []swagger.Order, err error) {
	var response *http.Response

	params := map[string]interface{}{}
	params["symbol"] = symbol
	params["text"] = "cancel order with bitmex api"

	orders, response, err = b.client.OrderApi.OrderCancelAll(b.ctx, params)
	if err != nil {
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) CancelOrder(oid string) (order swagger.Order, err error) {
	var response *http.Response
	var orders []swagger.Order

	params := map[string]interface{}{}
	params["orderID"] = oid
	params["text"] = "cancel order with bitmex api"

	orders, response, err = b.client.OrderApi.OrderCancel(b.ctx, params)
	if err != nil {
		return
	}
	if len(orders) != 1 {
		err = errors.New("order error")
		return
	}
	order = orders[0]
	b.onResponse(response)
	return
}

func (b *BitMEX) CloseOrder(side string, ordType string, price float64, orderQty int32, postOnly bool, timeInForce string, symbol string) (order swagger.Order, err error) {
	var response *http.Response

	params := map[string]interface{}{}
	params["symbol"] = symbol
	params["side"] = side
	params["ordType"] = ordType
	params["orderQty"] = float32(orderQty)
	if price > 0.0 {
		params["price"] = price // Limit order only
	}
	params["text"] = `close with bitmex api`

	//timeInForce := "FillOrKill"	// 全数执行或立刻取消
	if timeInForce != "" {
		params["timeInForce"] = timeInForce
	}

	execInst := "Close"
	if postOnly {
		execInst += ",ParticipateDoNotInitiate"
	}
	params["execInst"] = execInst
	order, response, err = b.client.OrderApi.OrderNew(b.ctx, symbol, params)
	if err != nil {
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) RequestWithdrawal(currency string, amount float32, address string, otpToken string, fee float64) (trans swagger.Transaction, err error) {
	var response *http.Response
	params := map[string]interface{}{}
	if otpToken != "" {
		params["otpToken"] = otpToken
	}
	if fee >= 0 {
		params["fee"] = fee
	}
	trans, response, err = b.client.UserApi.UserRequestWithdrawal(b.ctx, currency, amount, address, params)
	if err != nil {
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) ConfirmWithdrawal(token string) (trans swagger.Transaction, err error) {
	var response *http.Response
	trans, response, err = b.client.UserApi.UserConfirmWithdrawal(token)
	if err != nil {
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) onResponsePublic(response *http.Response) {
	//log.Printf("X-Ratelimit-Limit: %v", response.Header[`X-Ratelimit-Limit`])
	//log.Printf("X-Ratelimit-Remaining: %v", response.Header[`X-Ratelimit-Remaining`])
	//log.Printf("X-Ratelimit-Reset: %v", response.Header[`X-Ratelimit-Reset`])

	xLimit := response.Header.Get(`X-Ratelimit-Limit`)
	xRemaining := response.Header.Get(`X-Ratelimit-Remaining`)
	xReset := response.Header.Get(`X-Ratelimit-Reset`)

	b.rateLimitMutexPublic.Lock()
	defer b.rateLimitMutexPublic.Unlock()

	if xLimit != "" {
		b.rateLimitPublic.Limit, _ = strconv.ParseInt(xLimit, 10, 64)
	}
	if xRemaining != "" {
		b.rateLimitPublic.Remaining, _ = strconv.ParseInt(xRemaining, 10, 64)
	}
	if xReset != "" {
		b.rateLimitPublic.Reset, _ = strconv.ParseInt(xReset, 10, 64)
	}
}

func (b *BitMEX) onResponse(response *http.Response) {
	xLimit := response.Header.Get(`X-Ratelimit-Limit`)
	xRemaining := response.Header.Get(`X-Ratelimit-Remaining`)
	xReset := response.Header.Get(`X-Ratelimit-Reset`)

	b.rateLimitMutex.Lock()
	defer b.rateLimitMutex.Unlock()

	if xLimit != "" {
		b.rateLimit.Limit, _ = strconv.ParseInt(xLimit, 10, 64)
	}
	if xRemaining != "" {
		b.rateLimit.Remaining, _ = strconv.ParseInt(xRemaining, 10, 64)
	}
	if xReset != "" {
		b.rateLimit.Reset, _ = strconv.ParseInt(xReset, 10, 64)
	}
}
