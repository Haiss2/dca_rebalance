package trade

import (
	"context"

	futu "github.com/adshao/go-binance/v2/futures"
	"go.uber.org/zap"
)

const (
	QtyPrecision   = 4
	PricePrecision = 1
)

type TradeModule struct {
	l      *zap.SugaredLogger
	client *futu.Client
}

func NewTradeModule(client *futu.Client) *TradeModule {
	return &TradeModule{
		l:      zap.S(),
		client: client,
	}
}

func (t *TradeModule) GetTradeHistories(symbol string, limit int) ([]*futu.AccountTrade, error) {
	return t.client.NewListAccountTradeService().
		Symbol(symbol).
		Limit(limit).
		Do(context.Background())
}

func (t *TradeModule) CalcCommissionPNL(symbol string, orderID int64) (comm, pnl float64) {
	trades, err := t.GetTradeHistories(symbol, 10)
	if err != nil {
		t.l.Debugw("failed to get trades histories", "symbol", symbol, "limit", 10, "err", err)
		return
	}

	for _, trade := range trades {
		if trade.OrderID == orderID {
			comm += SToF(trade.Commission)
			pnl += SToF(trade.RealizedPnl)
		}
	}
	return
}
