package bitmex

import (
	"encoding/json"
	"github.com/sumorf/bitmex-api/swagger"
	"testing"
)

func TestParseOrders(t *testing.T) {
	raw := `[{"orderID":"fe492d85-9bd6-4147-c286-c9149163a3be","clOrdID":"","clOrdLinkID":"","account":149029,"symbol":"XBTUSD","side":"Buy","simpleOrderQty":null,"orderQty":20,"price":2100,"displayQty":null,"stopPx":null,"pegOffsetValue":null,"pegPriceType":"","currency":"USD","settlCurrency":"XBt","ordType":"Limit","timeInForce":"GoodTillCancel","execInst":"ParticipateDoNotInitiate","contingencyType":"","exDestination":"XBME","ordStatus":"New","triggered":"","workingIndicator":true,"ordRejReason":"","simpleLeavesQty":null,"leavesQty":20,"simpleCumQty":null,"cumQty":0,"avgPx":null,"multiLegReportingType":"SingleSecurity","text":"open long with bitmex api","transactTime":"2019-04-08T08:53:00.460Z","timestamp":"2019-04-08T08:53:00.460Z"},{"orderID":"f4c8e673-fab9-b41d-51b1-35b064c27e26","clOrdID":"","clOrdLinkID":"","account":149029,"symbol":"XBTUSD","side":"Buy","simpleOrderQty":null,"orderQty":20,"price":2100,"displayQty":null,"stopPx":null,"pegOffsetValue":null,"pegPriceType":"","currency":"USD","settlCurrency":"XBt","ordType":"Limit","timeInForce":"GoodTillCancel","execInst":"ParticipateDoNotInitiate","contingencyType":"","exDestination":"XBME","ordStatus":"New","triggered":"","workingIndicator":true,"ordRejReason":"","simpleLeavesQty":null,"leavesQty":20,"simpleCumQty":null,"cumQty":0,"avgPx":null,"multiLegReportingType":"SingleSecurity","text":"open long with bitmex api","transactTime":"2019-04-08T08:53:22.279Z","timestamp":"2019-04-08T08:53:22.279Z"},{"orderID":"f31997b8-809c-1ce2-20ef-9ce5fba4886c","clOrdID":"","clOrdLinkID":"","account":149029,"symbol":"XBTUSD","side":"Buy","simpleOrderQty":null,"orderQty":30,"price":3600,"displayQty":null,"stopPx":null,"pegOffsetValue":null,"pegPriceType":"","currency":"USD","settlCurrency":"XBt","ordType":"Limit","timeInForce":"GoodTillCancel","execInst":"","contingencyType":"","exDestination":"XBME","ordStatus":"New","triggered":"","workingIndicator":true,"ordRejReason":"","simpleLeavesQty":null,"leavesQty":30,"simpleCumQty":null,"cumQty":0,"avgPx":null,"multiLegReportingType":"SingleSecurity","text":"Submission from testnet.bitmex.com","transactTime":"2019-04-09T08:15:12.704Z","timestamp":"2019-04-09T08:15:12.704Z"}]`

	var orders []*swagger.Order
	err := json.Unmarshal([]byte(raw), &orders)
	if err != nil {
		t.Error(err)
	}

	for _, v := range orders {
		t.Logf("%v", *v)
	}
}
