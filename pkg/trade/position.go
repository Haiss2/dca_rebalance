package trade

import (
	"context"

	futu "github.com/adshao/go-binance/v2/futures"
)

func (t *TradeModule) GetAllPositions() ([]*futu.PositionRisk, error) {
	return t.client.NewGetPositionRiskService().
		Do(context.Background())
}

func (t *TradeModule) GetPosition(symbol string) ([]*futu.PositionRisk, error) {
	return t.client.NewGetPositionRiskService().
		Symbol(symbol).
		Do(context.Background())
}

func (t *TradeModule) UpdateLeverage(symbol string, leverage int) (*futu.SymbolLeverage, error) {
	return t.client.NewChangeLeverageService().
		Symbol(symbol).
		Leverage(leverage).
		Do(context.Background())
}
