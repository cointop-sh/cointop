package timedata

import (
	"math"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
)

// Resample the [timestamp,value] data given to numsteps between start-end (returns numSteps+1 points).
// If the data does not extend past start/end then there will likely be NaN in the output data.
func ResampleTimeSeriesData(data [][]float64, start float64, end float64, numSteps int) [][]float64 {
	var newData [][]float64
	l := len(data)
	step := (end - start) / float64(numSteps)
	for pos := start; pos <= end; pos += step {
		idx := sort.Search(l, func(i int) bool { return data[i][0] >= pos })
		var val float64
		if idx == 0 {
			val = math.NaN() // off the left
		} else if idx == l {
			val = math.NaN() // off the right
		} else {
			// between two points - linear interpolation
			left := data[idx-1]
			right := data[idx]
			dvdt := (right[1] - left[1]) / (right[0] - left[0])
			val = left[1] + (pos-left[0])*dvdt
		}
		newData = append(newData, []float64{pos, val})
	}
	return newData
}

// Assuming that the [timestamp,value] data provided is roughly evenly spaced, calculate that interval.
func CalculateTimeQuantum(data [][]float64) time.Duration {
	if len(data) > 1 {
		minTime := time.UnixMilli(int64(data[0][0]))
		maxTime := time.UnixMilli(int64(data[len(data)-1][0]))
		return time.Duration(int64(maxTime.Sub(minTime)) / int64(len(data)-1))
	}
	return 0
}

// Print out all the [timestamp,value] data provided
func DebugLogPriceData(data [][]float64) {
	for i := range data {
		log.Debugf("%s %.2f", time.Unix(int64(data[i][0]/1000), 0), data[i][1])
	}
}
