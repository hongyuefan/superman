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
	return nil
}

type OrderQuery struct {
	Symbol  string
	OrderId string
}

func (o *Order) GetPendingOrderIds() []OrderQuery {
	o.lock.RLock()
	defer o.lock.RUnlock()

	var orderQuerys []OrderQuery

	for _, order := range o.MOrder {

		orderQuerys = append(orderQuerys, OrderQuery{Symbol: order.InstrumentId, OrderId: order.OrderId})
	}

	return orderQuerys
}

func (o *Order) GetOrder(symbol, orderId string) (okex.SpotOrderListResult, error) {

	param := okex.SpotGetOrderParams{
		OrderId:      orderId,
		InstrumentId: symbol,
	}

	return o.client.SpotGetOrder(param)
}

func (o *Order) DoOrder(clId, typ, side, symbol, margin, price, size, notional string) (okex.SpotOrderResult, error) {

	param := okex.SpotOrderParams{
		ClientId:     clId,
		Type:         typ,
		Side:         side,
		InstrumentId: symbol,
		Margin:       margin,
		Price:        price,
		Size:         size,
		Notional:     notional,
	}
	return o.client.SpotDoOrder(param)
}

func (o *Order) CanselOrder(symbol, clId, orderId string) (okex.SpotOrderResult, error) {

	param := okex.SpotCanselOrderParams{
		InstrumentId: symbol,
		ClientId:     clId,
		OrderId:      orderId,
	}
	return o.client.SpotCanselOrder(param)
}

func (o *Order) QueryPendingOrders(symbol, from, to, limit string) ([]okex.SpotOrderListResult, error) {

	param := okex.SpotGetPendingOrderParams{
		InstrumentId: symbol,
		From:         from,
		To:           to,
		Limit:        limit,
	}
	return o.client.SpotGetPendingOrders(param)
}
