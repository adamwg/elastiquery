package elastiquery

import (
	"context"
	"encoding/json"
)

// ESClient is a version-independent ElasticSearch client.
type ESClient interface {
	RawQuery(string) ESQuery
	PrefixQuery(string, string) ESQuery
	TermQuery(string, string) ESQuery
	AndQuery(...ESQuery) ESQuery
	OrQuery(...ESQuery) ESQuery
}

// ESQuery is a version-independent ElasticSearch query.
type ESQuery interface {
	Do(ctx context.Context, index string, opts ...QueryOpt) (ESResult, error)
}

// ESResult is a version-independent ElasticSearch query result.
type ESResult interface {
	TotalHits() int64
	RawHits() []*json.RawMessage
}

type QueryOpt func(*QueryParams)

func WithOffset(offset int) QueryOpt {
	return func(qp *QueryParams) {
		qp.Offset = offset
	}
}

func WithLimit(limit int) QueryOpt {
	return func(qp *QueryParams) {
		qp.Limit = limit
	}
}

func WithSortField(field string) QueryOpt {
	return func(qp *QueryParams) {
		qp.SortField = field
	}
}

func WithReverseSort() QueryOpt {
	return func(qp *QueryParams) {
		qp.SortReverse = true
	}
}

// QueryParams are common parameters for queries. It should not normally be used
// directly - use the `With...` functions to modify the defaults instead.
type QueryParams struct {
	Offset      int
	Limit       int
	SortField   string
	SortReverse bool
}
