package pricing

import (
	"time"

	"github.com/Haiss2/dca/pkg/storage"
	futu "github.com/adshao/go-binance/v2/futures"
)

type Synthetic struct {
	symbol string
	client *futu.Client
	PK     *PriceKeeper
	H      *Historical
	K      *Kline
}

func NewSynthetic(symbol string, client *futu.Client, db *storage.RamStorage, duration time.Duration) *Synthetic {
	return &Synthetic{
		symbol: symbol,
		client: client,
		PK:     NewPriceKeeper(symbol, db, duration),
		H:      NewHistorical(symbol, client),
		K:      NewKline(symbol, client, false),
	}
}
