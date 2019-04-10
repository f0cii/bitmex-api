package bitmex

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty"
	"github.com/vmpartner/bitmex/rest"
	"github.com/vmpartner/bitmex/swagger"
	"net/http"
	"strconv"
	"time"
)

const (
	SIDE_BUY        = "Buy"
	SIDE_SELL       = "Sell"
	ORD_TYPE_MARKET = "Market"
	ORD_TYPE_LIMIT  = "Limit"
)

// {"name":"BitMEX API","version":"1.2.0","timestamp":1554709447283}
type Version struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
}

func (b *BitMEX) GetVersion() (version Version, time time.Duration, err error) {
	var resp *resty.Response

	url := "https://" + b.host + "/api/v1"
	resp, err = resty.R().Get(url)
	if err != nil {
		return
	}
	time = resp.Time()
	err = json.Unmarshal(resp.Body(), &version)
	return
}

func (b *BitMEX) GetWallet() (wallet swagger.Wallet, err error) {
	var response *http.Response

	wallet, response, err = rest.GetWallet(b.ctx)
	if err != nil {
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) GetOrderBookL2(depth int) (orderbook []swagger.OrderBookL2, err error) {
	var response *http.Response

	client := rest.GetClient(b.ctx)
	localVarOptionals := map[string]interface{}{}
	localVarOptionals["depth"] = float32(depth)

	orderbook, response, err = client.OrderBookApi.OrderBookGetL2(b.symbol, localVarOptionals)
	if err != nil {
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) GetPositions() (positions []swagger.Position, err error) {
	var response *http.Response

	client := rest.GetClient(b.ctx)
	localVarOptionals := map[string]interface{}{}
	localVarOptionals["filter"] = fmt.Sprintf(`{"symbol": "%s"}`, b.symbol)

	positions, response, err = client.PositionApi.PositionGet(b.ctx, localVarOptionals)
	if err != nil {
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) GetOrders() (orders []swagger.Order, err error) {
	var response *http.Response

	params := map[string]interface{}{}
	params["symbol"] = b.symbol
	params["filter"] = `{"open":true}`

	orders, response, err = rest.GetOrder(b.ctx, params)
	if err != nil {
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) NewOrder(side string, ordType string, price float64, orderQty int32, postOnly bool) (order swagger.Order, err error) {
	var response *http.Response

	params := map[string]interface{}{}
	params["symbol"] = b.symbol
	// params["clOrdID"] = ""	// 客户端委托ID
	params["side"] = side
	params["ordType"] = ordType
	params["orderQty"] = float32(orderQty)
	if price > 0.0 {
		params["price"] = price // Limit order only
	}
	params["text"] = `open with bitmex api`

	//timeInForce := "FillOrKill"	// 全数执行或立刻取消
	//if timeInForce != "" {
	//	params["timeInForce"] = timeInForce
	//}

	if postOnly {
		params["execInst"] = "ParticipateDoNotInitiate"
	}

	order, response, err = rest.NewOrder(b.ctx, params)
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

func (b *BitMEX) GetOrder(oid string) (order swagger.Order, err error) {
	var response *http.Response
	var orders []swagger.Order

	params := map[string]interface{}{}
	params["symbol"] = b.symbol
	params["filter"] = fmt.Sprintf(`{"orderID":"%s"}`, oid)

	orders, response, err = rest.GetOrder(b.ctx, params)
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

	order, response, err = rest.AmendOrder(b.ctx, params)
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
	params["symbol"] = b.symbol
	params["orderID"] = oid

	orders, response, err = rest.CancelOrder(b.ctx, params)
	if err != nil {
		return
	}
	if len(orders) != 1 {
		err = errors.New("order error")
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) CloseOrder(side string, ordType string, price float64, orderQty int32, postOnly bool) (order swagger.Order, err error) {
	var response *http.Response

	params := map[string]interface{}{}
	params["symbol"] = b.symbol
	params["side"] = side
	params["ordType"] = ordType
	params["orderQty"] = float32(orderQty)
	if price > 0.0 {
		params["price"] = price // Limit order only
	}
	params["text"] = `close with bitmex api`

	//timeInForce := "FillOrKill"	// 全数执行或立刻取消
	//if timeInForce != "" {
	//	params["timeInForce"] = timeInForce
	//}

	execInst := "Close"
	if postOnly {
		execInst += ",ParticipateDoNotInitiate"
	}
	params["execInst"] = execInst
	order, response, err = rest.NewOrder(b.ctx, params)
	if err != nil {
		return
	}
	b.onResponse(response)
	return
}

func (b *BitMEX) onResponse(response *http.Response) {
	//log.Printf("X-Ratelimit-Limit: %v", response.Header[`X-Ratelimit-Limit`])
	//log.Printf("X-Ratelimit-Remaining: %v", response.Header[`X-Ratelimit-Remaining`])
	//log.Printf("X-Ratelimit-Reset: %v", response.Header[`X-Ratelimit-Reset`])

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
