{{ $fields := .fields }}{{ $structName := .structName }}{{ $genericTypeNames := .genericTypeNames }}
type {{ lower .structName }}Builder{{ genericListWithTypes .genericTypeNames .genericTypes }} struct {
	data {{ .structName }}{{ genericList .genericTypeNames }}
}
{{ range .fieldNames }}
func (this *{{ lower $structName }}Builder{{ genericList $genericTypeNames }}) {{ capitalize . }}({{ . }} {{ index $fields . }}) *{{ lower $structName }}Builder{{ genericList $genericTypeNames }} {
	this.data.{{ . }} = {{ . }}
	return this
}
{{end}}
func (this *{{ lower $structName }}Builder{{ genericList .genericTypeNames }}) Build() {{ .structName }}{{ genericList .genericTypeNames }} {
	return this.data
}

func New{{ capitalize .structName }}Builder{{ genericListWithTypes .genericTypeNames .genericTypes }}() *{{ lower .structName }}Builder{{ genericList .genericTypeNames }} {
	return &{{ lower .structName }}Builder{{ genericList .genericTypeNames }}{
		data: {{ .structName }}{{ genericList .genericTypeNames }}{},
	}
}
