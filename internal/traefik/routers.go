package traefik

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

func GetRoutingRules(traefikUrl, traefikUsername, traefikPassword string) ([]byte, error) {
	// build the URL
	traefikURL, err := url.Parse(traefikUrl)
	if err != nil {
		return nil, err
	}
	traefikURL.Path = path.Join(traefikURL.Path, "/api/http/routers")

	// build the request
	req, err := http.NewRequest("GET", traefikURL.String(), nil)
	if err != nil {
		return nil, err
	}
	if traefikUsername != "" && traefikPassword != "" {
		req.SetBasicAuth(traefikUsername, traefikPassword)
	}

	// do the request and check for errors
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Traefik returned non-success code: %v", res.StatusCode)
	}

	// parse the response
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
