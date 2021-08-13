package traefik

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

func (c Client) GetRoutes() ([]string, error) {
	// build the URL
	traefikURL, err := url.Parse(c.Url)
	if err != nil {
		return []string{}, err
	}
	traefikURL.Path = path.Join(traefikURL.Path, "/api/http/routers")

	// build the request
	req, err := http.NewRequest("GET", traefikURL.String(), nil)
	if err != nil {
		return []string{}, err
	}
	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	// do the request and check for errors
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []string{}, err
	}
	if res.StatusCode != 200 {
		return []string{}, fmt.Errorf("Traefik returned non-success code: %v", res.StatusCode)
	}

	// parse the response
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []string{}, err
	}

	// get the json into go
	type Router struct {
		Rule string `json:"rule"`
	}
	var routers []Router
	err = json.Unmarshal(data, &routers)
	if err != nil {
		return []string{}, err
	}

	// get all the rules and remove duplicates
	var rules []string
	for _, r := range routers {
		rules = append(rules, r.Rule)
	}

	rules = sliceRemoveDuplicates(rules)

	return rules, nil
}

func sliceRemoveDuplicates(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
