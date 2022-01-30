{{ $fields := .fields }}
func New{{ .structName }}{{ genericListWithTypes .genericTypeNames .genericTypes }}({{ range $i, $e := .fieldNames  }}{{ if $i }}{{", "}}{{ end }}{{ $e }} {{ index $fields $e }}{{ end }}) *{{ .structName }}{{ genericList .genericTypeNames }} {
	return &{{ .structName }}{{ genericList .genericTypeNames }}{
{{ range .fieldNames }}		{{ . }}: {{ . }},
{{ end }}	}
}
