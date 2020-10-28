package cluster

type Status struct {
	Replicas      int32  `json:"replicas"`
	Version       string `json:"version"`
	Ready         int32  `json:"ready"`
	Grpc          string `json:"grpc"`
	Rest          string `json:"rest"`
	Api           string `json:"api"`
	Selector      string `json:"selector"`
	LicenseType   string `json:"license_type"`
	LicenseTo     string `json:"license_to"`
	LicenseExpire string `json:"license_expire"`
	Status        string `json:"status"`
}
