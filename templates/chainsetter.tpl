
func (this *{{ .structName }}{{ genericList .genericTypeNames }}) {{ .s }}et{{ capitalize .fieldName }}({{ .fieldName }} {{ .fieldType }}) *{{ .structName }}{{ genericList .genericTypeNames }} {
	this.{{ .fieldName }} = {{ .fieldName }}
	return this
}
