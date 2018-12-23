package base

import (
	"sync"

	"github.com/okcoin-okex/open-api-v3-sdk/okex-go-sdk-api"
)

type Order struct {
	client *okex.Client
	MOrder map[string]okex.SpotOrderListResult
	lock   sync.RWMutex
}

func NewOrder(client *okex.Client) *Order {
	return &Order{
		client: client,
		MOrder: make(map[string]okex.SpotOrderListResult),
	}
}

func (o *Order) LoadOrders() error {
	o.lock.Lock()
	defer o.lock.Unlock()

	results, err := o.client.SpotGetAllPendingOrders()

	if err != nil {
		return err
	}

	for _, result := range results {
		o.MOrder[result.OrderId] = result
	}
	return
}

func (o *Order) GetOrder(orderId, insId string) (okex.SpotOrderListResult, error) {

	param := okex.SpotGetOrderParams{
		OrderId:      orderId,
		InstrumentId: insId,
	}

	return o.client.SpotGetOrder(param)
}

func (o *Order) DoOrder(oid, typ, side, insId, margin, price, size, notional string) (okex.SpotOrderResult, error) {

	param := okex.SpotOrderParams{
		ClientId:     oid,
		Type:         typ,
		Side:         side,
		InstrumentId: insId,
		Margin:       margin,
		Price:        price,
		Size:         size,
		Notional:     notional,
	}
	return o.client.SpotDoOrder(param)
}

func (o *Order) CanselOrder(insId, clId, orderId string) (okex.SpotOrderResult, error) {

	param := okex.SpotCanselOrderParams{
		InstrumentId: insId,
		ClientId:     clId,
		OrderId:      orderId,
	}
	return o.client.SpotCanselOrder(param)
}

func (o *Order) QueryPendingOrders(insId, from, to, limit string) ([]okex.SpotOrderListResult, error) {

	param := okex.SpotGetPendingOrderParams{
		InstrumentId: insId,
		From:         from,
		To:           to,
		Limit:        limit,
	}
	return o.client.SpotGetPendingOrders(param)
}
