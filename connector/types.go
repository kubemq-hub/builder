package connector

const (
	promptBindingConfirm     = "<cyan>Here is the binding configuration:</>%s"
	promptBindingReconfigure = "<cyan>Lets reconfigure the binding:</>"

	promptBindingStart = `
<cyan>In the next steps, we will add Bindings for the connector %s.
Bindings represent a set of links between Sources and Targets.
Each link (Binding) consists of:</>
<yellow>Source:</> A connection which receives data from an external service
<yellow>Target:</> A connection which sends data to an external service
<yellow>Middlewares:</> Allows setting logging, retries, and rate-limiting functions between Source and Target
<cyan>Lets add our first binding:</>`

	promptConnectorContinue = "<cyan>Lets continue with connector settings:</>"

	connectorTemplate = `
<red>name:</> {{.Name}}
<red>namespace:</> {{.Namespace }}
<red>type:</> {{.Type}}
<red>replicas:</> {{.Replicas}}
<red>serviceType:</> {{.ServiceType}}
<red>config:</> |- 
<yellow>{{ .Config | indent 2}}</>`

	promptConnectorConfirm     = "<cyan>Here is the connector configuration:</>%s"
	promptConnectorReconfigure = "<cyan>Lets reconfigure the connector:</>"
)
