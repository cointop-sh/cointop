package table

// SortOrder int
type SortOrder int

// SortFn sort function
type SortFn func(interface{}, interface{}) bool

const (
	// SortNone sort none
	SortNone SortOrder = iota
	// SortAsc sort ascendinge
	SortAsc
	// SortDesc sort descending
	SortDesc
)

// SortBy sort by
type SortBy struct {
	index  int
	order  SortOrder
	sortFn SortFn
}
