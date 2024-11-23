package es

import (
	"encoding/json"
	"net/http"
)

// ErrorDetails encapsulate error details from Elasticsearch.
// It is used in e.g. elastic.Error and elastic.BulkResponseItem.
type ErrorDetails struct {
	Type      string                 `json:"type"`
	Reason    string                 `json:"reason"`
	Line      int32                  `json:"line,omitempty"`
	Col       int32                  `json:"col,omitempty"`
	CausedBy  map[string]interface{} `json:"caused_by,omitempty"`
	RootCause []*ErrorDetails        `json:"root_cause,omitempty"`
}

// TotalHits specifies total number of hits and its relation
type TotalHits struct {
	Value    int64  `json:"value"`    // value of the total hit count
	Relation string `json:"relation"` // how the value should be interpreted: accurate ("eq") or a lower bound ("gte")
}

// SearchHit is a single hit.
type SearchHit struct {
	Score  *float64        `json:"_score,omitempty"`  // computed score
	Index  string          `json:"_index,omitempty"`  // index name
	Type   string          `json:"_type,omitempty"`   // type meta field
	Id     string          `json:"_id,omitempty"`     // external or internal
	Source json.RawMessage `json:"_source,omitempty"` // stored document source
}

// SearchHits specifies the list of search hits.
type SearchHits struct {
	TotalHits *TotalHits   `json:"total,omitempty"`     // total number of hits found
	MaxScore  *float64     `json:"max_score,omitempty"` // maximum score of all hits
	Hits      []*SearchHit `json:"hits,omitempty"`      // the actual hits returned
}

// SearchResult is the result of a search in Elasticsearch.
type SearchResult struct {
	Header       http.Header   `json:"-"`
	TookInMillis int64         `json:"took,omitempty"`      // search time in milliseconds
	Hits         *SearchHits   `json:"hits,omitempty"`      // the actual search hits
	TimedOut     bool          `json:"timed_out,omitempty"` // true if the search timed out
	Error        *ErrorDetails `json:"error,omitempty"`     // only used in MultiGet
	Status       int           `json:"status,omitempty"`    // used in MultiSearch
}
