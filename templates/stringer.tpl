
func (this *{{ .structName }}{{ genericList .genericTypeNames }}) String() string {
	return fmt.Sprintf("{{ .packageName }}.{{ .structName }}{{"{"}}{{ range $i, $e := .fieldNames  }}{{ if $i }}{{", "}}{{ end }}{{ $e }}=%v{{ end }}{{"}"}}", {{ range $i, $e := .fieldNames  }}{{ if $i }}{{", "}}{{ end }}this.{{ $e }}{{ end }})
}
