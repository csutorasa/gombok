package main

import (
	"fmt"
	"io"
	"text/template"
)

var constructorTemplateString string = `{{ $fields := .fields }}
func New{{ .structName }}{{ genericListWithTypes .genericTypeNames .genericTypes }}({{ range $i, $e := .fieldNames  }}{{ if $i }}{{", "}}{{ end }}{{ $e }} {{ index $fields $e }}{{ end }}) *{{ .structName }}{{ genericList .genericTypeNames }} {
	return &{{ .structName }}{{ genericList .genericTypeNames }}{
{{ range .fieldNames }}		{{ . }}: {{ . }},
{{ end }}	}
}
`
var constructorTemplate *template.Template

func init() {
	constructorTemplate = loadTemplate("Constructor", constructorTemplateString)
}

func processConstructor(data *typeProcessorData) error {
	commands, found := hasComment(data.typeComments, "Constructor")
	if found {
		_, err := parseConstructorConfig(commands, data.structName)
		if err != nil {
			return err
		}
		fieldNames := []string{}
		for fieldName := range data.fields {
			commands, found := hasComment(data.fieldComments[fieldName], "Constructor")
			if found {
				fieldConfig, err := parseConstructorFieldConfig(commands, data.structName, fieldName)
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
		debugLogger.Printf("Generating Constructor for %s", data.structName)
		data.addCodeWriter(func(wr io.Writer) error {
			return constructorTemplate.Execute(wr, map[string]interface{}{
				"structName":       data.structName,
				"fieldNames":       fieldNames,
				"fields":           data.fields,
				"genericTypes":     data.genericTypes,
				"genericTypeNames": data.genericTypeNames,
			})
		})
	}
	return nil
}

type constructorConfig struct {
}

func parseConstructorConfig(commands []string, structName string) (*constructorConfig, error) {
	for _, command := range commands {
		return nil, fmt.Errorf("Invalid command %s on %s Constructor", command, structName)
	}
	return &constructorConfig{}, nil
}

type constructorFieldConfig struct {
	exclude bool
}

func parseConstructorFieldConfig(commands []string, structName, fieldName string) (*constructorFieldConfig, error) {
	exclude := false
	excludeSet := false
	for _, command := range commands {
		switch command {
		case "exclude":
			if excludeSet {
				return nil, fmt.Errorf("Invalid command %s on %s.%s Constructor", command, structName, fieldName)
			}
			exclude = true
			excludeSet = true
		case "include":
			if excludeSet {
				return nil, fmt.Errorf("Invalid command %s on %s.%s Constructor", command, structName, fieldName)
			}
			exclude = false
			excludeSet = true
		default:
			return nil, fmt.Errorf("Invalid command %s on %s.%s Constructor", command, structName, fieldName)
		}
	}
	return &constructorFieldConfig{
		exclude: exclude,
	}, nil
}
