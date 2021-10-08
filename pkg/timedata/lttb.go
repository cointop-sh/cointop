package timedata

import (
	"math"
)

// Copied from https://github.com/haoel/downsampling-algorithm

// Largest triangle three buckets (LTTB) data downsampling algorithm implementation
//  - Require: data . The original data
//  - Require: threshold . Number of data points to be returned
func LTTB(data []Point, threshold int) []Point {

	if threshold >= len(data) || threshold == 0 {
		return data // Nothing to do
	}

	sampledData := make([]Point, 0, threshold)

	// Bucket size. Leave room for start and end data points
	bucketSize := float64(len(data)-2) / float64(threshold-2)

	sampledData = append(sampledData, data[0]) // Always add the first point

	// We have 3 pointers represent for
	// > bucketLow - the current bucket's beginning location
	// > bucketMiddle - the current bucket's ending location,
	//                  also the beginning location of next bucket
	// > bucketHight - the next bucket's ending location.
	bucketLow := 1
	bucketMiddle := int(math.Floor(bucketSize)) + 1

	var prevMaxAreaPoint int

	for i := 0; i < threshold-2; i++ {

		bucketHigh := int(math.Floor(float64(i+2)*bucketSize)) + 1

		// Calculate point average for next bucket (containing c)
		avgPoint := calculateAverageDataPoint(data[bucketMiddle : bucketHigh+1])

		// Get the range for current bucket
		currBucketStart := bucketLow
		currBucketEnd := bucketMiddle

		// Point a
		pointA := data[prevMaxAreaPoint]

		maxArea := -1.0

		var maxAreaPoint int
		for ; currBucketStart < currBucketEnd; currBucketStart++ {

			area := calculateTriangleArea(pointA, avgPoint, data[currBucketStart])
			if area > maxArea {
				maxArea = area
				maxAreaPoint = currBucketStart
			}
		}

		sampledData = append(sampledData, data[maxAreaPoint]) // Pick this point from the bucket
		prevMaxAreaPoint = maxAreaPoint                       // This MaxArea point is the next's prevMAxAreaPoint

		//move to the next window
		bucketLow = bucketMiddle
		bucketMiddle = bucketHigh
	}

	sampledData = append(sampledData, data[len(data)-1]) // Always add last

	return sampledData
}

func LTTB2(data []Point, threshold int) []Point {
	buckets := splitDataBucket(data, threshold)
	samples := LTTBForBuckets(buckets)
	return samples
}

func LTTBForBuckets(buckets [][]Point) []Point {
	bucketCount := len(buckets)
	sampledData := make([]Point, 0)

	sampledData = append(sampledData, buckets[0][0])

	lastSelectedDataPoint := buckets[0][0]
	for i := 1; i < bucketCount-1; i++ {
		bucket := buckets[i]
		averagePoint := calculateAveragePoint(buckets[i+1])

		maxArea := -1.0
		maxAreaIndex := -1
		for j := 0; j < len(bucket); j++ {
			point := bucket[j]
			area := calculateTriangleArea(lastSelectedDataPoint, point, averagePoint)

			if area > maxArea {
				maxArea = area
				maxAreaIndex = j
			}
		}
		lastSelectedDataPoint := bucket[maxAreaIndex]
		sampledData = append(sampledData, lastSelectedDataPoint)
	}
	sampledData = append(sampledData, buckets[len(buckets)-1][0])
	return sampledData
}
