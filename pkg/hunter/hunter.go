package hunter

import (
	"sync"

	"github.com/Haiss2/dca/pkg/pricing"
	"github.com/Haiss2/dca/pkg/telegram"
	"github.com/Haiss2/dca/pkg/trade"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

type Hunter struct {
	sync.Mutex
	l          *zap.SugaredLogger
	symbol     string
	synth      *pricing.Synthetic
	longTrade  *trade.TradeModule
	shortTrade *trade.TradeModule
	tele       *telegram.TelegramBot
}

func NewHunter(
	c *cli.Context, symbol string,
	synth *pricing.Synthetic,
	longTrade, shortTrade *trade.TradeModule,
	tele *telegram.TelegramBot,
) *Hunter {
	return &Hunter{
		l:          zap.S(),
		symbol:     symbol,
		synth:      synth,
		longTrade:  longTrade,
		shortTrade: shortTrade,
		tele:       tele,
	}
}

func (h *Hunter) Hunt() {
	h.longTrade.CreateLimitOrder(h.symbol, 18000.123131, 0.012313, futures.SideTypeBuy)
}
