package web

type AnnotationSchema struct {
	Display interface{} `json:"x-display,omitempty"`
	Props   interface{} `json:"x-props,omitempty"`
	Class   interface{} `json:"x-class,omitempty"`
}

func NewAnnotationSchema() *AnnotationSchema {
	return &AnnotationSchema{}
}

func (a *AnnotationSchema) SetDisplay(value interface{}) *AnnotationSchema {
	a.Display = value
	return a
}

func (a *AnnotationSchema) SetProps(value interface{}) *AnnotationSchema {
	a.Props = value
	return a
}
func (a *AnnotationSchema) SetClass(value interface{}) *AnnotationSchema {
	a.Class = value
	return a
}
