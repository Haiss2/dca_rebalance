package pricing

import "github.com/Haiss2/dca/pkg/storage"

func getPriceAtTs(ps []storage.Price, id, timestamp, interval int64) (storage.Price, int64) {
	if timestamp >= ps[id].Timestamp {
		for index := id; index < int64(len(ps)); index++ {
			if index == int64(len(ps))-1 || ps[index+1].Timestamp > timestamp {
				return storage.Price{
					Timestamp: timestamp,
					Price:     ps[index].Price,
				}, index
			}
		}
	}
	return storage.Price{}, 0
}
