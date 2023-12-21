package main

import (
	"fmt"
	"io"
	"text/template"
)

var stringerTemplateString string = `
func (this *{{ .structName }}) String() string {
	return fmt.Sprintf("{{ .packageName }}.{{ .structName }}{{"{"}}{{ range $i, $e := .fieldNames  }}{{ if $i }}{{", "}}{{ end }}{{ $e }}=%v{{ end }}{{"}"}}", {{ range $i, $e := .fieldNames  }}{{ if $i }}{{", "}}{{ end }}this.{{ $e }}{{ end }})
}
`
var stringerTemplate *template.Template

func init() {
	stringerTemplate = loadTemplate("Stringer", stringerTemplateString)
}

func processStringer(data *typeProcessorData) error {
	commands, found := hasComment(data.typeComments, "Stringer")
	if found {
		_, err := parseStringerConfig(commands, data.structName)
		if err != nil {
			return err
		}
		fieldNames := []string{}
		for _, fieldName := range data.fieldNames {
			commands, found := hasComment(data.fieldComments[fieldName], "Stringer")
			if found {
				fieldConfig, err := parseStringerFieldConfig(commands, data.structName, fieldName)
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
		data.addImport("fmt")
		debugLogger.Printf("Generating Stringer for %s", data.structName)
		data.addCodeWriter(func(wr io.Writer) error {
			return stringerTemplate.Execute(wr, map[string]interface{}{
				"packageName": data.packageName,
				"structName":  data.structName,
				"fieldNames":  fieldNames,
			})
		})
	}
	return nil
}

type stringerConfig struct {
}

func parseStringerConfig(commands []string, structName string) (*stringerConfig, error) {
	for _, command := range commands {
		return nil, fmt.Errorf("invalid command %s on %s Stringer", command, structName)
	}
	return &stringerConfig{}, nil
}

type stringerFieldConfig struct {
	exclude bool
}

func parseStringerFieldConfig(commands []string, structName, fieldName string) (*stringerFieldConfig, error) {
	exclude := false
	excludeSet := false
	for _, command := range commands {
		switch command {
		case "exclude":
			if excludeSet {
				return nil, fmt.Errorf("invalid command %s on %s.%s Stringer", command, structName, fieldName)
			}
			exclude = true
			excludeSet = true
		case "include":
			if excludeSet {
				return nil, fmt.Errorf("invalid command %s on %s.%s Stringer", command, structName, fieldName)
			}
			exclude = false
			excludeSet = true
		default:
			return nil, fmt.Errorf("invalid command %s on %s.%s Stringer", command, structName, fieldName)
		}
	}
	return &stringerFieldConfig{
		exclude: exclude,
	}, nil
}
