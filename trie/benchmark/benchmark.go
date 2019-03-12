// Package benchmark provides internally used benchmark support
package benchmark

// Config defines the variable inputs struct in one benchmark.
type Config struct {
	KeyCnt int64
	KeyLen uint32
	ValLen uint32
}

// SearchResult show the key search result with a constructed data.
// Used to transfer benchmark result currently.
// SearchResult also defines the column titles when output to a chart.
type SearchResult struct {
	KeyCnt                int64
	KeyLen                uint32
	ExsitingKeyNsPerOp    int64
	NonexsitentKeyNsPerOp int64
}
