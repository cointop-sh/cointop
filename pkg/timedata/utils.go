package timedata

import (
	"math"
)

// Copied from https://github.com/haoel/downsampling-algorithm

func calculateTriangleArea(pa, pb, pc Point) float64 {
	area := ((pa.X-pc.X)*(pb.Y-pa.Y) - (pa.X-pb.X)*(pc.Y-pa.Y)) * 0.5
	return math.Abs(area)
}

func calculateAverageDataPoint(points []Point) (avg Point) {

	for _, point := range points {
		avg.X += point.X
		avg.Y += point.Y
	}
	l := float64(len(points))
	avg.X /= l
	avg.Y /= l
	return avg
}

func splitDataBucket(data []Point, threshold int) [][]Point {

	buckets := make([][]Point, threshold)
	for i := range buckets {
		buckets[i] = make([]Point, 0)
	}
	// First and last bucket are formed by the first and the last data points
	buckets[0] = append(buckets[0], data[0])
	buckets[threshold-1] = append(buckets[threshold-1], data[len(data)-1])

	// so we only have N - 2 buckets left to fill
	bucketSize := float64(len(data)-2) / float64(threshold-2)

	//slice remove the first and last point
	d := data[1 : len(data)-1]

	for i := 0; i < threshold-2; i++ {
		bucketStartIdx := int(math.Floor(float64(i) * bucketSize))
		bucketEndIdx := int(math.Floor(float64(i+1)*bucketSize)) + 1
		if i == threshold-3 {
			bucketEndIdx = len(d)
		}
		buckets[i+1] = append(buckets[i+1], d[bucketStartIdx:bucketEndIdx]...)
	}

	return buckets
}

func calculateAveragePoint(points []Point) Point {
	l := len(points)
	var p Point
	for i := 0; i < l; i++ {
		p.X += points[i].X
		p.Y += points[i].Y
	}

	p.X /= float64(l)
	p.Y /= float64(l)
	return p

}

func peakAndTroughPointIndex(points []Point) (int, int) {
	max := -0.1
	min := math.MaxFloat64
	minIdx := 0
	maxIdx := 0
	for i := 0; i < len(points); i++ {
		if points[i].Y > max {
			max = points[i].Y
			maxIdx = i
		}
		if points[i].Y < min {
			min = points[i].Y
			minIdx = i
		}
	}
	return maxIdx, minIdx
}
