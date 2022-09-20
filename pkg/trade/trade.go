package trade

import (
	"context"

	"github.com/Haiss2/dca/pkg/common"
	futu "github.com/adshao/go-binance/v2/futures"
	"go.uber.org/zap"
)

type TradeModule struct {
	l       *zap.SugaredLogger
	client  *futu.Client
	configs map[string]common.SymbolConfig
}

func NewTradeModule(client *futu.Client, configs []common.SymbolConfig) *TradeModule {
	configsM := make(map[string]common.SymbolConfig)
	for _, config := range configs {
		configsM[config.Symbol] = config
	}
	return &TradeModule{
		l:       zap.S(),
		client:  client,
		configs: configsM,
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
