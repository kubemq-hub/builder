package survey

type Entry interface {
	Complete() error
	Render(target interface{}) error
}
