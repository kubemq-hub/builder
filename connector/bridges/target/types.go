package target

const (
	promptTargetFirstConnection = "<cyan>Lets add our first connection for kind %s:</>"
	promptShowTarget            = "<cyan>Showing Targets configuration:</>"
)
const targetTemplate = `
<red>targets:</>
  <red>kind:</> {{.Kind}}
  <red>connections:</>
{{ .ConnectionSpec | indent 2}}
`
