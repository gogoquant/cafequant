package backtest

import (
	"errors"
	"fmt"
)

// MetricHandler defines the handling of metrics to a data event
type MetricHandler interface {
	Add(string, float64) error
	Get(string) (float64, bool)
}

// Metric holds metric propertys to a data point.
type Metric map[string]float64

// Add ads a value to the metrics map
func (m Metric) Add(key string, value float64) error {
	if m == nil {
		m = make(map[string]float64)
	}

	if key == "" {
		return errors.New("invalid key given")
	}

	m[key] = value
	return nil
}

// Get return a metric by name, if not found it returns false.
func (m Metric) Get(key string) (float64, bool) {
	value, ok := m[key]
	return value, ok
}

// Mean calculates the average for a slice of float64 values.
func Mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	var total float64

	for _, value := range values {
		total += value
	}

	return total / float64(len(values))
}

// SMA calculates the simple moving average of a given slice of data points.
func SMA(values []float64, period int) ([]float64, error) {
	var result []float64

	if len(values) == 0 {
		return result, fmt.Errorf("no values given")
	}

	// enough values ?
	if len(values) < period {
		return result, fmt.Errorf("invalid length of values, given %v, needs %v", len(values), period)
	}

	for i := range values {
		if i+1 >= period {
			avg := Mean(values[i+1-period : i+1])
			result = append(result, avg)
		}
	}

	return result, nil
}

// EMA calculates the Exponential Moving Average for the
// supplied slice of float64 values for a given period
func EMA(values []float64, period int) ([]float64, error) {
	var result []float64

	sma, err := SMA(values, period)
	if err != nil {
		return result, err
	}

	// multiplier = 2 / (period + 1)
	var multiplier = float64(2) / float64(period+1)

	// use sma as first ema
	result = append(result, sma[0])

	for i := (len(values) - len(sma)) + 1; i < len(values); i++ {
		// current val of values at index i
		currentVal := values[i]

		// last value of result slice - last calculatet ema
		lastEma := result[len(result)-1]

		ema := (currentVal-lastEma)*multiplier + lastEma

		result = append(result, ema)
	}

	return result, nil
}
