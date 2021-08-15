package prometheus

// example of a TargetGroup:
// [
//   {
//     "targets":
//       [ "<host>", ... ],
//     "labels":
//       { "<labelname>": "<labelvalue>", ... }
//   },
//   {
//     ...
//   },
//   ...
// ]

type TargetGroup struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels,omitempty"`
}

type TargetGroups []TargetGroup
