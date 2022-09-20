package hunter

import (
	"github.com/Haiss2/dca/pkg/pricing"
	"github.com/Haiss2/dca/pkg/telegram"
	"github.com/Haiss2/dca/pkg/trade"
	"github.com/urfave/cli"
)

type Hunter struct {
}

func NewHunter(
	c *cli.Context, synth *pricing.Synthetic,
	longTrade, shortTrade *trade.TradeModule,
	tele *telegram.TelegramBot,
) *Hunter {
	return &Hunter{}
}

func (h *Hunter) Hunt() {}
