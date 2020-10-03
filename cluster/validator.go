package cluster

type Validator interface {
	Validate() error
}
