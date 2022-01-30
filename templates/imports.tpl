{{ $imports := .imports }}
import (
{{ range .imports }}	{{ if .HasName }}{{ .Name }} {{ end }}"{{ .Path }}"
{{ end }})
