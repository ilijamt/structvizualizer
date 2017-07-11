package structvizualizer

const Template string = `
digraph {{ .GraphName }} {

edge [color=gray50, fontname=Calibri, fontsize=11]
node [shape=record, fontname=Calibri, fontsize=11]

# Define all the models first
{{ range $key, $value := .Objects }}
{{- $key }} [ label="{ { {{ $value.Name }} }|{{ join $value.Label "|" }} }" ]
{{ end }}

# Connect all the models to where they are tied to
{{ range $key, $value := .Connections }}
{{- $value.Embedded }} -> {{ $value.Base }} [ label="{{ $value.Label }}" ]
{{ end }}

}
`
