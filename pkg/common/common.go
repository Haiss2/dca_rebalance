package common

import (
	"math"

	"github.com/Haiss2/dca/pkg/storage"
)

const (
	epsilon = 1e-6
)

func Mean(ps []storage.Price) float64 {
	sum := 0.0
	for _, p := range ps {
		sum += p.Price
	}
	return sum / float64(len(ps))
}

func Std(ps []storage.Price) float64 {
	mean := Mean(ps)
	sum := 0.0
	for _, p := range ps {
		sum += math.Pow(p.Price-mean, 2)
	}
	return math.Pow(sum/float64(len(ps)), 0.5) / mean * 100
}

func Equal(x, y float64) bool {
	return math.Abs(x-y) < epsilon
}
