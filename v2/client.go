// Package v2 implements a client for ElasticSearch 2.x servers.
package v2

import (
	"github.com/adamwg/elastiquery"
	"gopkg.in/olivere/elastic.v3"
)

type client struct {
	wrap *elastic.Client
}

var _ elastiquery.ESClient = &client{}

func (c *client) RawQuery(query string) elastiquery.ESQuery {
	return &rawQuery{
		client: c.wrap,
		query:  query,
	}
}

func (c *client) PrefixQuery(field, prefix string) elastiquery.ESQuery {
	return &prefixQuery{
		client: c.wrap,
		field:  field,
		search: prefix,
	}
}

func (c *client) TermQuery(field, term string) elastiquery.ESQuery {
	return &termQuery{
		client: c.wrap,
		field:  field,
		search: term,
	}
}

func (c *client) AndQuery(queries ...elastiquery.ESQuery) elastiquery.ESQuery {
	return &boolQuery{
		client:     c.wrap,
		op:         AND,
		subQueries: queries,
	}
}

func (c *client) OrQuery(queries ...elastiquery.ESQuery) elastiquery.ESQuery {
	return &boolQuery{
		client:     c.wrap,
		op:         OR,
		subQueries: queries,
	}
}

func NewClient(url string) (elastiquery.ESClient, error) {
	ecl, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)

	if err != nil {
		return nil, err
	}

	return &client{
		wrap: ecl,
	}, nil
}