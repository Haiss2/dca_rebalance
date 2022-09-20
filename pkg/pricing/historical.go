package pricing

import (
	"context"
	"strconv"
	"sync"
	"time"

	futu "github.com/adshao/go-binance/v2/futures"
	"go.uber.org/zap"
)

const (
	updateHistoricalInterval = 30 * time.Minute
)

type Historical struct {
	l      *zap.SugaredLogger
	mu     sync.RWMutex
	symbol string
	client *futu.Client

	// 24hr rolling data
	priceChange        float64
	priceChangePercent float64
	weightedAvgPrice   float64
	openPrice          float64
	highPrice          float64
	lowPrice           float64
}

func NewHistorical(symbol string, client *futu.Client) *Historical {
	h := &Historical{
		l:      zap.S(),
		symbol: symbol,
		client: client,
	}

	h.updateHistorical()
	go func() {
		ticker := time.NewTicker(updateHistoricalInterval)
		defer ticker.Stop()

		for {
			<-ticker.C
			h.updateHistorical()
		}
	}()

	return h
}

func (h *Historical) updateHistorical() {
	priceChanges, err := h.client.NewListPriceChangeStatsService().
		Symbol(h.symbol).
		Do(context.Background())
	if err != nil || len(priceChanges) == 0 {
		h.l.Errorw("failed to update historical", "symbol", h.symbol, "err", err)
		return
	}

	h.UpdateWithMutex(priceChanges[0])
}

func (h *Historical) UpdateWithMutex(change *futu.PriceChangeStats) {
	h.mu.Lock()
	defer h.mu.Unlock()

	priceChange, _ := strconv.ParseFloat(change.PriceChange, 64)
	priceChangePercent, _ := strconv.ParseFloat(change.PriceChangePercent, 64)
	weightedAvgPrice, _ := strconv.ParseFloat(change.WeightedAvgPrice, 64)
	openPrice, _ := strconv.ParseFloat(change.OpenPrice, 64)
	highPrice, _ := strconv.ParseFloat(change.HighPrice, 64)
	lowPrice, _ := strconv.ParseFloat(change.LowPrice, 64)

	h.priceChange = priceChange
	h.priceChangePercent = priceChangePercent
	h.weightedAvgPrice = weightedAvgPrice
	h.openPrice = openPrice
	h.highPrice = highPrice
	h.lowPrice = lowPrice
}

func (h *Historical) PriceChange() float64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.priceChange
}

func (h *Historical) PriceChangePercent() float64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.priceChangePercent
}

func (h *Historical) WeightedAvgPrice() float64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.weightedAvgPrice
}

func (h *Historical) OpenPrice() float64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.openPrice
}

func (h *Historical) HighPrice() float64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.highPrice
}

func (h *Historical) LowPrice() float64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.lowPrice
}
