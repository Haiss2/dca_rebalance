package trade

import (
	"context"
	"fmt"

	futu "github.com/adshao/go-binance/v2/futures"
)

func (t *TradeModule) GetOrder(symbol string, id int64, clientId string) (*futu.Order, error) {
	return t.client.NewGetOrderService().
		Symbol(symbol).
		OrderID(id).
		OrigClientOrderID(clientId).
		Do(context.Background())
}

func (t *TradeModule) GetAllOrders(symbol string) ([]*futu.Order, error) {
	return t.client.NewListOpenOrdersService().
		Symbol(symbol).
		Do(context.Background())
}

func (t *TradeModule) CancelOrder(symbol string, id int64, clientId string) (*futu.CancelOrderResponse, error) {
	return t.client.NewCancelOrderService().
		Symbol(symbol).
		OrderID(id).
		OrigClientOrderID(clientId).
		Do(context.Background())
}

func (t *TradeModule) CancelAllOrder(symbol string) error {
	orders, err := t.GetAllOrders(symbol)
	if err != nil {
		return err
	}
	if len(orders) == 0 {
		return nil
	}
	return t.client.NewCancelAllOpenOrdersService().
		Symbol(symbol).
		Do(context.Background())
}

func (t *TradeModule) CreateMarketOrder(symbol string, qty float64, side futu.SideType) (*futu.Order, error) {
	order, err := t.client.NewCreateOrderService().
		Symbol(symbol).
		Side(side).
		Type(futu.OrderTypeMarket).
		Quantity(fmt.Sprintf("%f", RoundDown(qty, PricePrecision))).
		Do(context.Background())
	if err != nil {
		t.l.Infow("cannot create order", "err", err)
		return nil, err
	}
	t.l.Infow("create order successfully", "id", order.OrderID)
	return createOrderRespToOrder(order), nil
}

func (t *TradeModule) CreateLimitOrder(symbol string, price, qty float64, side futu.SideType) (*futu.Order, error) {
	order, err := t.client.NewCreateOrderService().
		Symbol(symbol).
		Side(side).
		Type(futu.OrderTypeLimit).
		TimeInForce(futu.TimeInForceTypeGTC).
		Price(fmt.Sprintf("%f", RoundDown(price, PricePrecision))).
		Quantity(fmt.Sprintf("%f", RoundDown(qty, QtyPrecision))).
		ReduceOnly(false).
		Do(context.Background())
	if err != nil {
		t.l.Infow("cannot create order", "err", err)
		return nil, err
	}
	t.l.Infow("create order successfully", "id", order.OrderID)
	return createOrderRespToOrder(order), nil
}

func createOrderRespToOrder(orderResp *futu.CreateOrderResponse) *futu.Order {
	return &futu.Order{
		Symbol:           orderResp.Symbol,
		OrderID:          orderResp.OrderID,
		ClientOrderID:    orderResp.ClientOrderID,
		Price:            orderResp.Price,
		OrigQuantity:     orderResp.OrigQuantity,
		ExecutedQuantity: orderResp.ExecutedQuantity,
		CumQuote:         orderResp.CumQuote,
		Status:           orderResp.Status,
		TimeInForce:      orderResp.TimeInForce,
		Type:             orderResp.Type,
		Side:             orderResp.Side,
		Time:             orderResp.UpdateTime,
		UpdateTime:       orderResp.UpdateTime,
		ReduceOnly:       orderResp.ReduceOnly,
		PositionSide:     orderResp.PositionSide,
	}
}
