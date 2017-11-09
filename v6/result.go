package v6

import (
	"encoding/json"

	"github.com/adamwg/elastiquery"
	elastic "gopkg.in/olivere/elastic.v6"
)

type result struct {
	wrap *elastic.SearchResult
}

var _ elastiquery.ESResult = &result{}

func (r *result) TotalHits() int64 {
	return r.wrap.TotalHits()
}

func (r *result) RawHits() []*json.RawMessage {
	var ret []*json.RawMessage

	if r.wrap.TotalHits() == 0 {
		return ret
	}

	// Avoid re-allocation in the loop
	ret = make([]*json.RawMessage, 0, r.wrap.TotalHits())
	for _, hit := range r.wrap.Hits.Hits {
		ret = append(ret, hit.Source)
	}

	return ret
}
