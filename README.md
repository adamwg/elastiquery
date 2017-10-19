# elastiquery - A CLI tool for querying ElasticSearch

elastiquery is a CLI tool that can query ElasticSearch servers and print results
as JSON. The resulting JSON can then be manipulated with your favorite CLI
tools, e.g. `jq`.

## Compatibility

Currently only ElasticSearch version 2.x is supported, however elastiquery is
built in such a way that adding support for additional versions is simple.

## Usage

```
Usage of elastiquery:
  -es-url string
        ElasticSearch server URL
  -index string
        ElasticSearch index
  -limit int
        Number of results to return
  -offset int
        Number of results to skip
  -or
        Require only one of the given queries to match, rather than all of them
  -prefixes string
        Semicolon-separated ElasticSearch term queries in the form field=term
  -raw string
        Raw ElasticSearch JSON query
  -reverse
        Sort in reverse order
  -sort-by string
        Field name to sort results by
  -terms string
        Semicolon-separated ElasticSearch term queries in the form field=term
  -timeout duration
        Timeout for ElasticSearch queries (default 30s)
```

For convenience, the ElasticSearch URL and index can also be specified via the
environment variables `ES_URL` and `ES_INDEX`.

## Installation

`go get github.com/adamwg/elastiquery/cmd/elastiquery`

## Examples

Find all records whose `app` field is `myapp` and whose `message` field starts
with `error`:

```console
$ elastiquery -terms 'app=myapp' -prefixes 'message=error'
```

Find all records whose 'app' field is 'myapp' or 'yourapp':

```console
$ elastiquery -or -terms 'app=myapp;app=yourapp'
```
