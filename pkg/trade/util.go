package trade

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// RoundingNumberDown ...
func RoundDown(n float64, precision int) float64 {
	rounding := math.Pow10(precision)
	return math.Floor(n*rounding) / rounding
}

// RoundingNumberAuto round up or down depend on the number
func RoundAuto(n float64, precision int) float64 {
	d := decimal.NewFromFloat(n)
	dr := d.Round(int32(precision))
	f, _ := dr.Float64()
	return f
}

// RoundingNumberUp ...
func RoundUp(n float64, precision int) float64 {
	rounding := math.Pow10(precision)
	return math.Ceil(n*rounding) / rounding
}

func GetPrecision(precisionString string) (int, error) {
	f, err := strconv.ParseFloat(precisionString, 64)
	if err != nil {
		return 0, err
	}
	if f == 0 {
		return 0, errors.New("precision string is zero") // should not happen
	}
	precision := math.Log10(1 / f)
	return int(precision), nil
}

func FloatToString(f float64) string {
	df := decimal.NewFromFloat(f)
	return df.String()
}

func FloatToPointer(f float64) *float64 {
	return &f
}

func StringToPointer(s string) *string {
	return &s
}

func Int64ToPointer(i int64) *int64 {
	return &i
}

func BoolToPointer(b bool) *bool {
	return &b
}

func TimeToPointer(t time.Time) *time.Time {
	return &t
}

func ParseStringToFloat(l *zap.SugaredLogger, name, s string) float64 {
	if l == nil {
		l = zap.S()
	}
	if s == "" {
		l.Warnw("variable is empty string", "name", name)
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		l.Errorw("cannot parse number", "name", name, "raw", s, "err", err)
	}
	return f
}

func SToF(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
