package pricing

import (
	"context"
	"sync"
	"time"

	futu "github.com/adshao/go-binance/v2/futures"
	"go.uber.org/zap"
)

const (
	updateKlineInterval = 30 * time.Minute
	klineInterval       = "15m"
	klineDuration       = 4 * time.Hour
)

type Kline struct {
	l      *zap.SugaredLogger
	mu     sync.RWMutex
	symbol string
	client *futu.Client
	klines []*futu.Kline
}

func NewKline(symbol string, client *futu.Client, klineEnabled bool) *Kline {
	if !klineEnabled {
		return nil
	}
	k := &Kline{
		l:      zap.S(),
		symbol: symbol,
		client: client,
		klines: make([]*futu.Kline, 0),
	}

	k.updateKline()
	go func() {
		ticker := time.NewTicker(updateHistoricalInterval)
		defer ticker.Stop()

		for {
			<-ticker.C
			k.updateKline()
		}
	}()

	return k
}

func (k *Kline) updateKline() {
	endTime := time.Now().UnixMilli()
	klines, err := k.client.NewKlinesService().
		Symbol(k.symbol).
		Interval(klineInterval).
		EndTime(endTime).
		StartTime(endTime - klineDuration.Milliseconds()).
		Do(context.Background())
	if err != nil || len(klines) == 0 {
		k.l.Errorw("failed to update kline", "symbol", k.symbol, "err", err)
		return
	}

	k.l.Debugw("update kline successfully", "symbol", k.symbol, "data", klines)

	k.UpdateWithMutex(klines)
}

func (k *Kline) UpdateWithMutex(klines []*futu.Kline) {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.klines = klines
}

func (k *Kline) GetKline() []*futu.Kline {
	k.mu.Lock()
	defer k.mu.Unlock()

	return k.klines
}
