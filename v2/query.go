package v2

import (
	"context"

	"github.com/adamwg/elastiquery"
	elastic "gopkg.in/olivere/elastic.v3"
)

type v2Query interface {
	ElasticQuery() elastic.Query
}

type rawQuery struct {
	client *elastic.Client
	query  string
}

var _ elastiquery.ESQuery = &rawQuery{}

func (rq *rawQuery) Do(ctx context.Context, index string, opts ...elastiquery.QueryOpt) (elastiquery.ESResult, error) {
	q := rq.ElasticQuery()
	return doSearch(ctx, rq.client, index, q, opts...)
}

func (rq *rawQuery) ElasticQuery() elastic.Query {
	return elastic.NewRawStringQuery(rq.query)
}

type prefixQuery struct {
	client *elastic.Client
	field  string
	search string
}

var _ elastiquery.ESQuery = &prefixQuery{}

func (pq *prefixQuery) Do(ctx context.Context, index string, opts ...elastiquery.QueryOpt) (elastiquery.ESResult, error) {
	q := pq.ElasticQuery()
	return doSearch(ctx, pq.client, index, q, opts...)
}

func (pq *prefixQuery) ElasticQuery() elastic.Query {
	return elastic.NewPrefixQuery(pq.field, pq.search)
}

type termQuery struct {
	client *elastic.Client
	field  string
	search string
}

var _ elastiquery.ESQuery = &termQuery{}

func (tq *termQuery) Do(ctx context.Context, index string, opts ...elastiquery.QueryOpt) (elastiquery.ESResult, error) {
	q := tq.ElasticQuery()
	return doSearch(ctx, tq.client, index, q, opts...)
}

func (tq *termQuery) ElasticQuery() elastic.Query {
	return elastic.NewTermQuery(tq.field, tq.search)
}

type boolOp int

const (
	AND boolOp = iota
	OR
)

type boolQuery struct {
	client     *elastic.Client
	op         boolOp
	subQueries []elastiquery.ESQuery
}

var _ elastiquery.ESQuery = &boolQuery{}

func (bq *boolQuery) Do(ctx context.Context, index string, opts ...elastiquery.QueryOpt) (elastiquery.ESResult, error) {
	q := bq.ElasticQuery()
	return doSearch(ctx, bq.client, index, q, opts...)
}

func (bq *boolQuery) ElasticQuery() elastic.Query {
	var queries []elastic.Query
	for _, sq := range bq.subQueries {
		queries = append(queries, sq.(v2Query).ElasticQuery())
	}

	q := elastic.NewBoolQuery()
	switch bq.op {
	case AND:
		q = q.Must(queries...)

	case OR:
		q = q.Should(queries...)
	}

	return q
}

type rangeQuery struct {
	client *elastic.Client
	field  string
	from   interface{}
	to     interface{}
}

func (rq *rangeQuery) Do(ctx context.Context, index string, opts ...elastiquery.QueryOpt) (elastiquery.ESResult, error) {
	q := rq.ElasticQuery()
	return doSearch(ctx, rq.client, index, q, opts...)
}

func (rq *rangeQuery) ElasticQuery() elastic.Query {
	return elastic.NewRangeQuery(rq.field).From(rq.from).To(rq.to)
}

func doSearch(ctx context.Context, client *elastic.Client, index string, query elastic.Query,
	opts ...elastiquery.QueryOpt) (elastiquery.ESResult, error) {

	params := elastiquery.QueryParams{
		Offset:      0,
		Limit:       500,
		SortField:   "@timestamp",
		SortReverse: false,
	}
	for _, opt := range opts {
		opt(&params)
	}

	search := client.Search(index).
		Sort(params.SortField, params.SortReverse).
		From(params.Offset).
		Size(params.Limit).
		Query(query)

	res, err := search.DoC(ctx)
	if err != nil {
		return nil, err
	}

	return &result{res}, nil
}
