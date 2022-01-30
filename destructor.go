package main

import (
	_ "embed"
	"fmt"
	"io"
	"text/template"
)

//go:embed templates/destructor.tpl
var destructorTemplateString string
var destructorTemplate *template.Template

func init() {
	destructorTemplate = loadTemplate("Destructor", destructorTemplateString)
}

func processDestructor(data *typeProcessorData) error {
	commands, found := hasComment(data.typeComments, "Destructor")
	if found {
		config, err := parseDestructorConfig(commands, data.structName)
		if err != nil {
			return err
		}
		fieldNames := []string{}
		for fieldName := range data.fields {
			commands, found := hasComment(data.fieldComments[fieldName], "Destructor")
			if found {
				fieldConfig, err := parseDestructorFieldConfig(commands, data.structName, fieldName)
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
		d := "D"
		if config.private {
			d = "d"
		}
		debugLogger.Printf("Generating Destructor for %s", data.structName)
		data.addCodeWriter(func(wr io.Writer) error {
			return destructorTemplate.Execute(wr, map[string]interface{}{
				"structName":       data.structName,
				"fieldNames":       fieldNames,
				"fields":           data.fields,
				"genericTypes":     data.genericTypes,
				"genericTypeNames": data.genericTypeNames,
				"d":                d,
			})
		})
	}
	return nil
}

type destructorConfig struct {
	private bool
}

func parseDestructorConfig(commands []string, structName string) (*destructorConfig, error) {
	private := false
	privateSet := false
	for _, command := range commands {
		switch command {
		case "private":
			if privateSet {
				return nil, fmt.Errorf("Invalid command %s on %s Destructor", command, structName)
			}
			private = true
			privateSet = true
		case "exported":
			if privateSet {
				return nil, fmt.Errorf("Invalid command %s on %s Destructor", command, structName)
			}
			private = false
			privateSet = true
		default:
			return nil, fmt.Errorf("Invalid command %s on %s Destructor", command, structName)
		}
	}
	return &destructorConfig{
		private: private,
	}, nil
}

type destructorFieldConfig struct {
	exclude bool
}

func parseDestructorFieldConfig(commands []string, structName, fieldName string) (*destructorFieldConfig, error) {
	exclude := false
	excludeSet := false
	for _, command := range commands {
		switch command {
		case "exclude":
			if excludeSet {
				return nil, fmt.Errorf("Invalid command %s on %s.%s Destructor", command, structName, fieldName)
			}
			exclude = true
			excludeSet = true
		case "include":
			if excludeSet {
				return nil, fmt.Errorf("Invalid command %s on %s.%s Destructor", command, structName, fieldName)
			}
			exclude = false
			excludeSet = true
		default:
			return nil, fmt.Errorf("Invalid command %s on %s.%s Destructor", command, structName, fieldName)
		}
	}
	return &destructorFieldConfig{
		exclude: exclude,
	}, nil
}
