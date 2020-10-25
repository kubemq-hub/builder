package bridges

const (
	promptSourceStart = `
<cyan>In the next steps, we will configure the Source connections.
We will set:</>
<yellow>Name -</> A unique name for the Source's binding
<yellow>Kind -</> A Source connection type 
<yellow>Connections -</> A list of connections properties based on the selected kind

<cyan>Lets start binding source configuration:</>`
	promptBindingAddConfirmation = "<cyan>Binding %s was added successfully</>"

	promptTargetStart = `
<cyan>In the next steps, we will configure the Target connections.
We will set:</>
<yellow>Name -</> A unique name for the Target's binding
<yellow>Kind -</> A Target connection type 
<yellow>Connections -</> A list of connections properties based on the selected kind

<cyan>Lets start binding target configuration:</>`

	bindingTemplate = `
<red>name:</> {{.Name}}
{{- .SourcesSpec -}}
{{- .TargetsSpec -}}
{{- .PropertiesSpec -}}
`
	promptBindingComplete           = "<cyan>We have completed Source and Target binding configurations\n</>"
	promptShowBinding               = "<cyan>Showing Binding %s configuration:</>"
	promptBindingDeleteConfirmation = "<cyan>Binding %s deleted successfully\n</>"
	promptBindingEditedConfirmation = "<red>Binding %s edited successfully\n</>"
)
