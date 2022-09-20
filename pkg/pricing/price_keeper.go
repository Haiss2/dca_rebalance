package pricing

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/Haiss2/dca/pkg/storage"
	futu "github.com/adshao/go-binance/v2/futures"
	"go.uber.org/zap"
)

const (
	defaultChSize          = 20
	removeOldPriceInterval = time.Second
	ignoreLastPriceRange   = 10 * time.Second
)

type PriceKeeper struct {
	l           *zap.SugaredLogger
	mu          sync.RWMutex
	db          *storage.RamStorage
	symbol      string
	lastPrice   float64
	duration    time.Duration
	subscribeCh map[string]chan *storage.Price
}

func NewPriceKeeper(symbol string, db *storage.RamStorage, duration time.Duration) *PriceKeeper {
	storageCh := make(chan *storage.Price)
	k := &PriceKeeper{
		l:        zap.S(),
		mu:       sync.RWMutex{},
		db:       db,
		symbol:   symbol,
		duration: duration,
		subscribeCh: map[string]chan *storage.Price{
			"storage": storageCh,
		},
	}
	go k.subscribeWsAggTradeServe()
	go k.storageData(storageCh)
	go k.removeOldPriceRoutine()
	return k
}

func (k *PriceKeeper) subscribeWsAggTradeServe() {
	k.l.Infow("start subscribe to binance ws", "symbol", k.symbol)

	handler := func(event *futu.WsAggTradeEvent) {
		p, _ := strconv.ParseFloat(event.Price, 64)
		pricePtr := &storage.Price{Price: p, Timestamp: event.Time}
		for _, ch := range k.subscribeCh {
			ch <- pricePtr
		}
		k.UpdateLastPrice(p)
	}

	errHandler := func(err error) {
		k.l.Infow("error while subscribe to binance ws, reconnecting...", "err", err)
		k.subscribeWsAggTradeServe()
	}

	doneC, _, err := futu.WsAggTradeServe(k.symbol, handler, errHandler)
	if err != nil {
		k.l.Infow("error while subscribe to binance ws, reconnecting...", "err", err)
		k.subscribeWsAggTradeServe()
	}

	<-doneC
	k.l.Infow("Closed ws connection!", "symbol", k.symbol)
}

func (k *PriceKeeper) UpdateLastPrice(price float64) {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.lastPrice = price
}

func (k *PriceKeeper) GetLastPrice() float64 {
	k.mu.Lock()
	defer k.mu.Unlock()

	return k.lastPrice
}

func (k *PriceKeeper) Subscribe(name string, ch chan *storage.Price) {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.subscribeCh[name] = ch
}

func (k *PriceKeeper) Unsubscribe(name string) {
	k.mu.Lock()
	defer k.mu.Unlock()

	ch := k.subscribeCh[name]
	delete(k.subscribeCh, name)
	go func() {
		time.Sleep(100 * time.Millisecond)
		close(ch)
	}()
}

func (k *PriceKeeper) storageData(storageCh chan *storage.Price) {
	for {
		p := <-storageCh
		k.db.SavePrice(k.symbol, p.Price, p.Timestamp)
	}
}

func (k *PriceKeeper) removeOldPriceRoutine() {
	ticker := time.NewTicker(removeOldPriceInterval)
	defer ticker.Stop()
	for {
		<-ticker.C
		expiredTs := time.Now().UTC().UnixMilli() - k.duration.Milliseconds()
		k.db.RemoveExpiredData(expiredTs)
	}
}

func (k *PriceKeeper) GetPrices(symbol string, interval int64) ([]storage.Price, float64, error) {
	k.mu.Lock()
	defer k.mu.Unlock()

	ps := k.db.GetPricesBySymbol(symbol)
	if len(ps) == 0 {
		return nil, 0, errors.New("empty price")
	}

	ps = k.smoothPrices(ps, interval)

	return ps, k.lastPrice, nil
}

func (k *PriceKeeper) smoothPrices(ps []storage.Price, interval int64) []storage.Price {
	var index int64 = 0
	result := make([]storage.Price, 0)
	timestamp := ps[0].Timestamp
	for {
		if timestamp > time.Now().Add(-ignoreLastPriceRange).UTC().UnixMilli() {
			break
		}
		p, newIndex := getPriceAtTs(ps, index, timestamp, interval)
		result = append(result, p)
		timestamp += interval
		index = newIndex
	}
	return result
}
