/*
	maintain existed positions
*/

package hunter

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Haiss2/dca/pkg/trade"
)

const (
	monitorBalanceInterval = 1 * time.Hour
)

func (h *Hunter) monitorBalance() {
	longer, err := h.longer.GetAccountInformation()
	if err != nil {
		h.l.Errorw("failed to get longer account information", "err", err)
		return
	}

	shorter, err := h.shorter.GetAccountInformation()
	if err != nil {
		h.l.Errorw("failed to get shorter account information", "err", err)
		return
	}

	lBalance, lPnl, lPos, lMargin := trade.SToF(longer.TotalWalletBalance),
		trade.SToF(longer.TotalUnrealizedProfit), trade.SToF(longer.TotalPositionInitialMargin),
		trade.SToF(longer.TotalMarginBalance)

	sBalance, sPnl, sPos, sMargin := trade.SToF(shorter.TotalWalletBalance),
		trade.SToF(shorter.TotalUnrealizedProfit), trade.SToF(shorter.TotalPositionInitialMargin),
		trade.SToF(shorter.TotalMarginBalance)

	msg := fmt.Sprintf(
		`Balance check
%s LastPrice: %v
Longer Assets:
  - BNB: %.2f | BUSD: %.2f | USDT: %.2f
  - Position (%d): %s
  - PNL: %.2f | Pos: %.2f | Margin: %.2f
Shorter Assets:
  - BNB: %.2f | BUSD: %.2f | USDT: %.2f
  - Position (%d): %s
  - PNL: %.2f | Pos: %.2f | Margin: %.2f
Summary:
  - WalletBalance: %.2f USDT
  - UnrealizedProfit: %.2f USDT
  - PositionInitialMargin: %.2f USDT
  - MarginBalance: %.2f USDT`,
		h.symbol, h.synth.PK.GetLastPrice(),
		getAsset(longer.Asset, "BNB"), getAsset(longer.Asset, "BUSD"), getAsset(longer.Asset, "USDT"),
		len(longer.Position), strings.Join(longer.Position, ", "),
		lPnl, lPos, lMargin,
		getAsset(shorter.Asset, "BNB"), getAsset(shorter.Asset, "BUSD"), getAsset(shorter.Asset, "USDT"),
		len(shorter.Position), strings.Join(shorter.Position, ", "),
		sPnl, sPos, sMargin,
		lBalance+sBalance, lPnl+sPnl, lPos+sPos, lMargin+sMargin,
	)
	h.tele.Notify(msg)
}

func getAsset(assets []trade.Asset, token string) float64 {
	for _, asset := range assets {
		if asset.Symbol == token {
			return trade.SToF(asset.Balance)
		}
	}
	return 0
}

type DcaJob struct {
	LongPrice   float64
	LongAmount  float64
	ShortPrice  float64
	ShortAmount float64
}

func (h *Hunter) checkPositionsWhenStaringApp() (pass bool, job DcaJob, err error) {
	// cancel order for longer
	err = h.longer.CancelAllOrder(h.symbol)
	if err != nil {
		h.l.Errorw("can not cancel orders for longer", "err", err)
		return
	}

	// cancel order for shorter
	err = h.shorter.CancelAllOrder(h.symbol)
	if err != nil {
		h.l.Errorw("can not cancel orders for shorter", "err", err)
		return
	}

	var openedLong, openedShort bool

	// get longer position
	lPos, err := h.longer.GetPosition(h.symbol)
	if err != nil {
		h.l.Errorw("failed to get position for longer", "err", err)
		return
	}
	lp := lPos[0]
	amount, _ := strconv.ParseFloat(lp.PositionAmt, 64)
	openedLong = amount > 0

	// get shorter position
	sPos, err := h.shorter.GetPosition(h.symbol)
	if err != nil {
		h.l.Errorw("failed to get position for shorter", "err", err)
		return
	}
	sp := sPos[0]
	amount, _ = strconv.ParseFloat(sp.PositionAmt, 64)
	openedShort = amount < 0

	h.l.Debugw("opening positions return from checkPositionsWhenStaringApp", "long", lp, "short", sp)

	if openedLong != openedShort {
		err = errors.New("2 ways need to same state")
		return
	}

	// pass means no open position
	return !openedLong, DcaJob{
		LongPrice:   trade.SToF(lp.EntryPrice),
		LongAmount:  trade.SToF(lp.PositionAmt),
		ShortPrice:  trade.SToF(sp.EntryPrice),
		ShortAmount: trade.SToF(sp.PositionAmt),
	}, nil
}
