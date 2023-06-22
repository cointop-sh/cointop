package timedata

import (
	"math"
	"sort"
	"time"

	"github.com/cointop-sh/cointop/pkg/humanize"
)

// ResampleTimeSeriesData resamples the given [timestamp,value] data to numsteps between start-end (returns numSteps+1 points).
// If the data does not extend past start/end then there will likely be NaN in the output data.
func ResampleTimeSeriesData(data [][]float64, start float64, end float64, numSteps int) [][]float64 {
	var newData [][]float64
	l := len(data)
	step := (end - start) / float64(numSteps)
	for pos := start; pos <= end; pos += step {
		idx := sort.Search(l, func(i int) bool { return data[i][0] >= pos })
		var val float64
		if idx == 0 {
			if data[0][0] == pos {
				val = data[0][1] // exactly left
			} else {
				val = math.NaN() // off the left
			}
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

// CalculateTimeQuantum determines the given [timestamp,value] data
func CalculateTimeQuantum(data [][]float64) time.Duration {
	if len(data) > 1 {
		minTime := time.UnixMilli(int64(data[0][0]))
		maxTime := time.UnixMilli(int64(data[len(data)-1][0]))
		return time.Duration(int64(maxTime.Sub(minTime)) / int64(len(data)-1))
	}
	return 0
}

// BuildTimeSeriesLabels returns a list of short labels representing time values from the given [timestamp,value] data
func BuildTimeSeriesLabels(data [][]float64) []string {
	minTime := time.UnixMilli(int64(data[0][0]))
	maxTime := time.UnixMilli(int64(data[len(data)-1][0]))
	timeRange := maxTime.Sub(minTime)

	var timeFormat string
	if timeRange.Hours() < 24 {
		timeFormat = "15:04"
	} else if timeRange.Hours() < 24*7 {
		timeFormat = "Mon 15:04"
	} else if timeRange.Hours() < 24*365 {
		timeFormat = "02-Jan"
	} else {
		timeFormat = "Jan 2006"
	}
	var labels []string
	for i := range data {
		labelTime := time.UnixMilli(int64(data[i][0]))
		labels = append(labels, humanize.FormatTime(labelTime, timeFormat))
	}
	return labels
}
