package render

const genericTemplate = `{{emoji .Level}} <b>{{esc .Level}}</b> · {{esc .Data.Type}} {{resolvedIcon .Resolved}}
{{- with (get .Data.Data "name") | str}}{{if .}}
<b>Name:</b> {{esc .}} ({{esc $.Target.Type}}){{end}}{{end}}
{{json .Data.Data}}`

var defaultTemplates = map[string]string{
	"SERVERUNREACHABLE": `{{emoji .Level}} <b>{{esc .Level}}</b> · ServerUnreachable {{resolvedIcon .Resolved}}
<b>Server:</b> {{esc (str (get .Data.Data "name"))}}{{with str (get .Data.Data "region")}}
<b>Region:</b> {{esc .}}{{end}}{{with get .Data.Data "err"}}{{$err := .}}{{with get $err "error"}}
<b>Error:</b> <code>{{esc (str .)}}</code>{{end}}{{with get $err "trace"}}
<pre>{{esc (str .)}}</pre>{{end}}{{end}}`,

	"SERVERCPU": `{{emoji .Level}} <b>{{esc .Level}}</b> · ServerCpu {{resolvedIcon .Resolved}}
<b>Server:</b> {{esc (str (get .Data.Data "name"))}}
<b>CPU:</b> <code>{{esc (str (get .Data.Data "cpu_perc"))}}%</code>`,

	"SERVERMEM": `{{emoji .Level}} <b>{{esc .Level}}</b> · ServerMem {{resolvedIcon .Resolved}}
<b>Server:</b> {{esc (str (get .Data.Data "name"))}}
<b>Memory:</b> <code>{{esc (str (get .Data.Data "mem_used_gb"))}} / {{esc (str (get .Data.Data "mem_total_gb"))}} GB</code>`,

	"SERVERDISK": `{{emoji .Level}} <b>{{esc .Level}}</b> · ServerDisk {{resolvedIcon .Resolved}}
<b>Server:</b> {{esc (str (get .Data.Data "name"))}}
<b>Path:</b> <code>{{esc (str (get .Data.Data "path"))}}</code>
<b>Disk:</b> <code>{{esc (str (get .Data.Data "used_gb"))}} / {{esc (str (get .Data.Data "total_gb"))}} GB</code>`,

	"SERVERTEMP": `{{emoji .Level}} <b>{{esc .Level}}</b> · ServerTemp {{resolvedIcon .Resolved}}
<b>Server:</b> {{esc (str (get .Data.Data "name"))}}
<b>Temp:</b> <code>{{esc (str (get .Data.Data "temp"))}}°C</code>`,

	"CONTAINERSTATECHANGE": `{{emoji .Level}} <b>{{esc .Level}}</b> · ContainerStateChange {{resolvedIcon .Resolved}}
<b>Container:</b> {{esc (str (get .Data.Data "name"))}}
<b>Server:</b> {{esc (str (get .Data.Data "server_name"))}}
<b>State:</b> <code>{{esc (str (get .Data.Data "from"))}}</code> → <code>{{esc (str (get .Data.Data "to"))}}</code>`,

	"STACKSTATECHANGE": `{{emoji .Level}} <b>{{esc .Level}}</b> · StackStateChange {{resolvedIcon .Resolved}}
<b>Stack:</b> {{esc (str (get .Data.Data "name"))}}
<b>Server:</b> {{esc (str (get .Data.Data "server_name"))}}
<b>State:</b> <code>{{esc (str (get .Data.Data "from"))}}</code> → <code>{{esc (str (get .Data.Data "to"))}}</code>`,

	"STACKAUTOUPDATED": `{{emoji .Level}} <b>{{esc .Level}}</b> · StackAutoUpdated {{resolvedIcon .Resolved}}
<b>Stack:</b> {{esc (str (get .Data.Data "name"))}}
<b>Server:</b> {{esc (str (get .Data.Data "server_name"))}}{{with get .Data.Data "images"}}
<b>Images:</b> {{json .}}{{end}}`,

	"DEPLOYMENTSTATECHANGE": `{{emoji .Level}} <b>{{esc .Level}}</b> · DeploymentStateChange {{resolvedIcon .Resolved}}
<b>Deployment:</b> {{esc (str (get .Data.Data "name"))}}
<b>Server:</b> {{esc (str (get .Data.Data "server_name"))}}
<b>State:</b> <code>{{esc (str (get .Data.Data "from"))}}</code> → <code>{{esc (str (get .Data.Data "to"))}}</code>`,

	"BUILDFAILED": `{{emoji .Level}} <b>{{esc .Level}}</b> · BuildFailed {{resolvedIcon .Resolved}}
<b>Build:</b> {{esc (str (get .Data.Data "name"))}}{{with str (get .Data.Data "err")}}
<b>Error:</b> <code>{{esc .}}</code>{{end}}`,

	"RESOURCESYNCPENDINGUPDATES": `{{emoji .Level}} <b>{{esc .Level}}</b> · ResourceSyncPendingUpdates {{resolvedIcon .Resolved}}
<b>Sync:</b> {{esc (str (get .Data.Data "name"))}}`,

	"AWSBUILDERTERMATIONFAILED": `{{emoji .Level}} <b>{{esc .Level}}</b> · AwsBuilderTerminationFailed {{resolvedIcon .Resolved}}
<b>Region:</b> {{esc (str (get .Data.Data "region"))}}
<b>Instance:</b> <code>{{esc (str (get .Data.Data "instance_id"))}}</code>`,

	"NONE": `{{emoji .Level}} <b>{{esc .Level}}</b> · Alert {{resolvedIcon .Resolved}}
{{json .Data.Data}}`,
}
