package connector

type Status struct {
	Replicas int32  `json:"replicas"`
	Type     string `json:"type"`
	Image    string `json:"image"`
	Api      string `json:"api"`
	Status   string `json:"status"`
}
