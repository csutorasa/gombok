package main

import (
	_ "embed"
	"fmt"
	"io"
	"text/template"
)

//go:embed templates/builder.tpl
var builderTemplateString string
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
		for fieldName := range data.fields {
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

type builderConfig struct {
}

func parseBuilderConfig(commands []string, structName string) (*builderConfig, error) {
	for _, command := range commands {
		return nil, fmt.Errorf("Invalid command %s on %s Builder", command, structName)
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
				return nil, fmt.Errorf("Invalid command %s on %s.%s Builder", command, structName, fieldName)
			}
			exclude = true
			excludeSet = true
		case "include":
			if excludeSet {
				return nil, fmt.Errorf("Invalid command %s on %s.%s Builder", command, structName, fieldName)
			}
			exclude = false
			excludeSet = true
		default:
			return nil, fmt.Errorf("Invalid command %s on %s.%s Builder", command, structName, fieldName)
		}
	}
	return &builderFieldConfig{
		exclude: exclude,
	}, nil
}
