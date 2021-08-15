package helpers

import "encoding/json"

func GetRouteRules(routes []byte) ([]string, error) {
	type Router struct {
		Rule string `json:"rule"`
	}

	// get the json into go
	var routers []Router
	err := json.Unmarshal(routes, &routers)
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
