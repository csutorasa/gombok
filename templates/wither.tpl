{{ $fields := .fields }}{{ $fieldName := .fieldName }}
func (this *{{ .structName }}{{ genericList .genericTypeNames }}) {{ .w }}ith{{ capitalize .fieldName }}({{ .fieldName }} {{ .fieldType }}) *{{ .structName }}{{ genericList .genericTypeNames }} {
	return &{{ .structName }}{{ genericList .genericTypeNames }}{
{{ range .fieldNames }}		{{ . }}: {{ if eq . $fieldName }}{{ . }}{{ else }}this.{{ . }}{{ end }},
{{ end }}	}
}
