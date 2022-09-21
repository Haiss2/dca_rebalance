package hunter

import (
	"sync"
	"time"

	"github.com/Haiss2/dca/pkg/pricing"
	"github.com/Haiss2/dca/pkg/telegram"
	"github.com/Haiss2/dca/pkg/trade"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

type Hunter struct {
	sync.Mutex
	l       *zap.SugaredLogger
	symbol  string
	synth   *pricing.Synthetic
	longer  *trade.TradeModule
	shorter *trade.TradeModule
	tele    *telegram.TelegramBot
}

func NewHunter(
	c *cli.Context, symbol string,
	synth *pricing.Synthetic,
	longer, shorter *trade.TradeModule,
	tele *telegram.TelegramBot,
) *Hunter {
	h := &Hunter{
		l:       zap.S(),
		symbol:  symbol,
		synth:   synth,
		longer:  longer,
		shorter: shorter,
		tele:    tele,
	}

	time.Sleep(1 * time.Second)
	h.monitorBalance()
	go func() {
		ticker := time.NewTicker(monitorBalanceInterval)
		defer ticker.Stop()

		for {
			<-ticker.C
			h.monitorBalance()
		}
	}()

	return h
}

func (h *Hunter) Hunt() {
	// Firstly we have to maintain old position when staring app
	pass, job, err := h.checkPositionsWhenStaringApp()
	if err != nil {
		h.l.Errorw("Failed to start hunter", "err", err)
		return
	}

	if pass { // no opened positions
		h.HuntAndTrade()
	} else { // there are 2 opposite side positions
		h.Dca(job)
	}
}

func (h *Hunter) HuntAndTrade() {

}
