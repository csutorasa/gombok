package main

import (
	"fmt"
	"io"
	"text/template"
)

var builderTemplateString string = `{{ $fields := .fields }}{{ $structName := .structName }}
type {{ lower .structName }}Builder struct {
	data {{ .structName }}
}
{{ range .fieldNames }}
func (this *{{ lower $structName }}Builder) {{ capitalize . }}({{ . }} {{ index $fields . }}) *{{ lower $structName }}Builder {
	this.data.{{ . }} = {{ . }}
	return this
}
{{end}}
func (this *{{ lower $structName }}Builder) Build() {{ .structName }} {
	return this.data
}

func New{{ capitalize .structName }}Builder() *{{ lower .structName }}Builder {
	return &{{ lower .structName }}Builder{
		data: {{ .structName }}{},
	}
}
`
var builderTemplate *template.Template

func init() {
	builderTemplate = loadTemplate("Builder", builderTemplateString)
}

func processBuilder(data *typeProcessorData) error {
	commands, found := hasComment(data.typeComments, "Builder")
	if found {
		_, err := parseBuilderConfig(commands, data.structName)
		if err != nil {
			return err
		}
		fieldNames := []string{}
		for _, fieldName := range data.fieldNames {
			commands, found := hasComment(data.fieldComments[fieldName], "Builder")
			if found {
				fieldConfig, err := parseBuilderFieldConfig(commands, data.structName, fieldName)
				if err != nil {
					return err
				}
				if !fieldConfig.exclude {
					fieldNames = append(fieldNames, fieldName)
				}
			} else {
				fieldNames = append(fieldNames, fieldName)
			}
		}
		debugLogger.Printf("Generating Builder for %s", data.structName)
		data.addCodeWriter(func(wr io.Writer) error {
			return builderTemplate.Execute(wr, map[string]interface{}{
				"structName": data.structName,
				"fieldNames": fieldNames,
				"fields":     data.fields,
			})
		})
	}
	return nil
}

type builderConfig struct {
}

func parseBuilderConfig(commands []string, structName string) (*builderConfig, error) {
	for _, command := range commands {
		return nil, fmt.Errorf("invalid command %s on %s Builder", command, structName)
	}
	return &builderConfig{}, nil
}

type builderFieldConfig struct {
	exclude bool
}

func parseBuilderFieldConfig(commands []string, structName, fieldName string) (*builderFieldConfig, error) {
	exclude := false
	excludeSet := false
	for _, command := range commands {
		switch command {
		case "exclude":
			if excludeSet {
				return nil, fmt.Errorf("nvalid command %s on %s.%s Builder", command, structName, fieldName)
			}
			exclude = true
			excludeSet = true
		case "include":
			if excludeSet {
				return nil, fmt.Errorf("invalid command %s on %s.%s Builder", command, structName, fieldName)
			}
			exclude = false
			excludeSet = true
		default:
			return nil, fmt.Errorf("invalid command %s on %s.%s Builder", command, structName, fieldName)
		}
	}
	return &builderFieldConfig{
		exclude: exclude,
	}, nil
}
