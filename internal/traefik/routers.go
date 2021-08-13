package traefik

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"

	"github.com/calvinbui/blackbox-traefik-sd/internal/helpers"
)

type Router struct {
	EntryPoints []string `json:"entryPoints"`
	Middlewares []string `json:"middlewares,omitempty"`
	Service     string   `json:"service"`
	Rule        string   `json:"rule"`
	Status      string   `json:"status"`
	Using       []string `json:"using"`
	Name        string   `json:"name"`
	Provider    string   `json:"provider"`
	TLS         struct {
		CertResolver string `json:"certResolver"`
		Domains      []struct {
			Main string `json:"main"`
		} `json:"domains"`
	} `json:"tls,omitempty"`
	Priority int64 `json:"priority,omitempty"`
}

func (c Client) GetRoutes() ([]Router, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", path.Join(c.Url, "/api/http/routers"), nil)
	if err != nil {
		return []Router{}, err
	}

	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	res, err := client.Do(req)
	if err != nil {
		return []Router{}, err
	}

	if res.StatusCode != 200 {
		return []Router{}, fmt.Errorf("Traefik returned non-success code: %v", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Router{}, err
	}

	var routers []Router
	err = json.Unmarshal(data, &routers)
	if err != nil {
		return []Router{}, err
	}

	return routers, nil
}

func GetHostsFromRouter(routes []Router) []string {
	var re = regexp.MustCompile(`(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9][a-z0-9-]{0,61}[a-z0-9]`)
	var targets []string

	for _, route := range routes {
		hosts := re.FindAllStringSubmatch(route.Rule, -1)
		for _, host := range hosts {
			targets = append(targets, host[0])
		}
	}

	targets = helpers.SliceRemoveDuplicates(targets)

	return targets
}
