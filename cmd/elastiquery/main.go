package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/adamwg/elastiquery"
	"github.com/adamwg/elastiquery/v2"
	"github.com/ianschenck/envflag"
)

func main() {
	// Global options can be set by either envvar or flag, with the flag taking
	// precedence.
	var esURL string
	envflag.StringVar(&esURL, "ES_URL", "", "ElasticSearch server URL")
	flag.StringVar(&esURL, "es-url", "", "ElasticSearch server URL")
	var esIndex string
	envflag.StringVar(&esIndex, "ES_INDEX", "", "ElasticSearch index")
	flag.StringVar(&esIndex, "index", "", "ElasticSearch index")

	timeout := flag.Duration("timeout", 30*time.Second, "Timeout for ElasticSearch queries")
	offset := flag.Int("offset", 0, "Number of results to skip")
	limit := flag.Int("limit", 0, "Number of results to return")
	sortField := flag.String("sort-by", "", "Field name to sort results by")
	sortReverse := flag.Bool("reverse", false, "Sort in reverse order")
	or := flag.Bool("or", false, "Require only one of the given queries to match, rather than all of them")

	rawQuery := flag.String("raw", "", "Raw ElasticSearch JSON query")
	termQueries := flag.String("terms", "", "Semicolon-separated ElasticSearch term queries in the form field=term")
	prefixQueries := flag.String("prefixes", "", "Semicolon-separated ElasticSearch term queries in the form field=term")

	envflag.Parse()
	flag.Parse()

	if esURL == "" {
		log.Fatal("ElasticSearch server address must be provided via flag or environment")
	}
	if esIndex == "" {
		log.Fatal("ElasticSearch index must be provided via flag or environment")
	}

	esVersion, err := elastiquery.GetServerVersion(esURL)
	if err != nil {
		log.Fatalf("Could not determine ES server version: %v", err)
	}
	log.Printf("ES server version is %s", esVersion)

	var client elastiquery.ESClient
	if strings.HasPrefix(esVersion, "2.") {
		client, err = v2.NewClient(esURL)
	}
	if err != nil {
		log.Fatalf("Error creating ElasticSearch client: %v", err)
	}
	if client == nil {
		log.Fatalf("Could not create client for ES version %s", esVersion)
	}

	var queries []elastiquery.ESQuery
	if *rawQuery != "" {
		queries = append(queries, client.RawQuery(*rawQuery))
	}
	if *termQueries != "" {
		qs := strings.Split(*termQueries, ";")
		for _, q := range qs {
			parts := strings.SplitN(q, "=", 2)
			if len(parts) != 2 {
				log.Fatalf("Invalid term query %q", q)
			}
			queries = append(queries, client.TermQuery(parts[0], parts[1]))
		}
	}
	if *prefixQueries != "" {
		qs := strings.Split(*prefixQueries, ";")
		for _, q := range qs {
			parts := strings.SplitN(q, "=", 2)
			if len(parts) != 2 {
				log.Fatalf("Invalid prefix query %q", q)
			}
			queries = append(queries, client.PrefixQuery(parts[0], parts[1]))
		}
	}

	if len(queries) == 0 {
		log.Fatal("No queries specified")
	}

	opts := parseQueryOpts(*offset, *limit, *sortField, *sortReverse)
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	var result elastiquery.ESResult
	if len(queries) == 1 {
		result, err = queries[0].Do(ctx, esIndex, opts...)
	} else if *or {
		result, err = client.OrQuery(queries...).Do(ctx, esIndex, opts...)
	} else {
		result, err = client.AndQuery(queries...).Do(ctx, esIndex, opts...)
	}

	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	if result.TotalHits() == 0 {
		log.Fatalf("Query returned no results")
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(result.RawHits())
}

func parseQueryOpts(offset, limit int, sortField string, sortReverse bool) []elastiquery.QueryOpt {
	var ret []elastiquery.QueryOpt

	if offset != 0 {
		ret = append(ret, elastiquery.WithOffset(offset))
	}
	if limit != 0 {
		ret = append(ret, elastiquery.WithLimit(limit))
	}
	if sortField != "" {
		ret = append(ret, elastiquery.WithSortField(sortField))
	}
	if sortReverse {
		ret = append(ret, elastiquery.WithReverseSort())
	}

	return ret
}
