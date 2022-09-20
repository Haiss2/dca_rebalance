package storage

import "sync"

type Price struct {
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

type ramDB map[string][]Price

type RamStorage struct {
	mu sync.RWMutex
	db ramDB
}

func NewRamStorage() *RamStorage {
	db := make(ramDB, 0)
	return &RamStorage{
		mu: sync.RWMutex{},
		db: db,
	}
}

func (r *RamStorage) SavePrice(s string, p float64, t int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.db[s] = append(r.db[s], Price{p, t})
}

func (r *RamStorage) GetPricesBySymbol(symbol string) []Price {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.db[symbol]
}

func (r *RamStorage) RemoveExpiredData(timestamp int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for symbol, prices := range r.db {
		start := 0
		for i, p := range prices {
			if p.Timestamp > timestamp {
				start = i
				break
			}
		}
		copy(r.db[symbol], r.db[symbol][start:])
		r.db[symbol] = r.db[symbol][:len(r.db[symbol])-start]
	}
}
