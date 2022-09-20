package trade

import (
	"context"

	futu "github.com/adshao/go-binance/v2/futures"
)

type Asset struct {
	Symbol  string
	Balance string
}
type AccountInfo struct {
	Asset                      []Asset
	Position                   []string
	AccountPosition            []*futu.AccountPosition
	TotalWalletBalance         string
	TotalUnrealizedProfit      string
	TotalPositionInitialMargin string
	TotalMarginBalance         string
}

func (t *TradeModule) GetAccountInformation() (*AccountInfo, error) {
	account, err := t.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	assets := make([]Asset, 0, 3)
	for _, asset := range account.Assets {
		if asset.Asset == "BNB" || asset.Asset == "USDT" || asset.Asset == "BUSD" {
			assets = append(assets, Asset{asset.Asset, asset.WalletBalance})
		}
	}

	positions := make([]string, 0)
	accountPositions := make([]*futu.AccountPosition, 0)
	for _, pos := range account.Positions {
		if pos.PositionInitialMargin != "0" {
			positions = append(positions, pos.Symbol)
			accountPositions = append(accountPositions, pos)
		}
	}

	return &AccountInfo{
		Asset:                      assets,
		Position:                   positions,
		AccountPosition:            accountPositions,
		TotalWalletBalance:         account.TotalWalletBalance,
		TotalUnrealizedProfit:      account.TotalUnrealizedProfit,
		TotalPositionInitialMargin: account.TotalPositionInitialMargin,
		TotalMarginBalance:         account.TotalMarginBalance,
	}, nil
}
