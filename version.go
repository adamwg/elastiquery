package elastiquery

import (
	"encoding/json"
	"errors"
	"net/http"
)

func GetServerVersion(url string) (string, error) {
	versionURL := url + "/_nodes/_all/version"
	resp, err := http.Get(versionURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	version := struct {
		Nodes map[string]struct {
			Version string
		}
	}{}
	if err := dec.Decode(&version); err != nil {
		return "", err
	}

	for _, node := range version.Nodes {
		if node.Version != "" {
			return node.Version, nil
		}
	}

	return "", errors.New("could not find version in server response")
}
