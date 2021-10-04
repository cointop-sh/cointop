package levenshtein

import (
	"math"
	"strings"
)

// DamerauLevenshteinDistance calculates the damerau-levenshtein distance between s1 and s2.
// Reference: [Damerau-Levenshtein Distance](http://en.wikipedia.org/wiki/Damerau%E2%80%93Levenshtein_distance)
// Note that this calculation's result isn't normalized. (not between 0 and 1.)
// and if s1 and s2 are exactly the same, the result is 0.
func DamerauLevenshteinDistance(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}
	s1Array := strings.Split(s1, "")
	s2Array := strings.Split(s2, "")
	lenS1Array := len(s1Array)
	lenS2Array := len(s2Array)
	m := make([][]int, lenS1Array+1)
	var cost int
	for i := range m {
		m[i] = make([]int, lenS2Array+1)
	}
	for i := 0; i < lenS1Array+1; i++ {
		for j := 0; j < lenS2Array+1; j++ {
			if i == 0 {
				m[i][j] = j
			} else if j == 0 {
				m[i][j] = i
			} else {
				cost = 0
				if s1Array[i-1] != s2Array[j-1] {
					cost = 1
				}
				m[i][j] = min(m[i-1][j]+1, m[i][j-1]+1, m[i-1][j-1]+cost)
				if i > 1 && j > 1 && s1Array[i-1] == s2Array[j-2] && s1Array[i-2] == s2Array[j-1] {
					m[i][j] = min(m[i][j], m[i-2][j-2]+cost)
				}
			}
		}
	}
	return m[lenS1Array][lenS2Array]
}

// min returns the minimum number of passed int slices.
func min(is ...int) int {
	min := math.MaxInt32
	for _, v := range is {
		if min > v {
			min = v
		}
	}
	return min
}
