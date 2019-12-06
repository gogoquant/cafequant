package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/xiyanxiyan10/quantcore/gobacktest"
	"github.com/xiyanxiyan10/quantcore/gobacktest/data"
	"github.com/xiyanxiyan10/quantcore/gobacktest/strategy"
)

// Result bundles the result of a single backtest
type Result struct {
	smaShort int
	smaLong  int
	result   float64
}

func main() {
	// create intervals for the short and long range
	shortRange := linspace(5, 50, 10)
	longRange := linspace(150, 250, 10)
	// create a slice for different test results
	results := []Result{}
	for _, short := range shortRange {
		for _, long := range longRange {
			results = append(results, Result{smaShort: short, smaLong: long})
		}
	}

	// initiate new backtester
	test := gobacktest.New()

	// define and load symbols
	symbols := []string{"SZG.DE"}
	test.SetSymbols(symbols)

	// create data provider and load data into the backtest
	data := &data.BarEventFromCSVFile{FileDir: "../testdata/bar/"}
	data.Load(symbols)
	test.SetData(data)

	startTest := time.Now()
	// iterate over every field in the matrix
	for i := range results {
		// create strategy provider and load into the backtest
		strategy := strategy.MovingAverageCross(results[i].smaShort, results[i].smaLong)
		test.SetStrategy(strategy)

		// run the backtest
		test.Run()

		// get the result and save to slice
		result, _ := test.Stats().TotalEquityReturn()
		fmt.Printf("backtest sma%d / sma%d with result %f%%\n", results[i].smaShort, results[i].smaLong, result*100)
		results[i].result = result

		test.Reset()
	}

	stopTest := time.Now()
	fmt.Printf("Complete backtest took %v sec\n", stopTest.Sub(startTest).Seconds())

	sortedResults := sortResults(results)
	// print best results
	fmt.Println("Best results:")
	for k := 0; k < 3; k++ {
		result := sortedResults[k]
		fmt.Printf("%v. SMA %v / SMA %v: %2f%%\n", k+1, result.smaShort, result.smaLong, result.result*100)
	}
	// print worst results
	fmt.Println("Worst results:")
	for k := 0; k < 3; k++ {
		result := sortedResults[len(sortedResults)-1-k]
		fmt.Printf("%v. SMA %v / SMA %v: %2f%%\n", k+1, result.smaShort, result.smaLong, result.result*100)
	}

}

// linspace returns a slice of n evenly spaced integers within a given range
// this function is similar to linspace() in the Python package Numpy
func linspace(start, stop, n int) []int {
	delta := stop - start
	step := delta / (n - 1)
	slice := make([]int, n, n)
	i := 0
	for x := start; x < stop; x += step {
		slice[i] = x
		i++
	}
	slice[n-1] = stop
	return slice
}

// sortResults sorts the result slice descending - highest result first
func sortResults(r []Result) []Result {
	sort.Slice(r, func(i, j int) bool {
		r1 := r[i]
		r2 := r[j]

		// if result is equal sort by smaShort
		if r1.result == r2.result {
			// if smaShort is equal sort by smaLong
			if r1.smaShort == r2.smaShort {
				// sort by smaLong
				return r1.smaLong < r2.smaLong
			}
			// sort by smaShort
			return r1.smaShort < r2.smaShort
		}
		// else sort by result
		return r1.result > r2.result
	})
	return r
}
