package main

import (
	_ "embed"
	"fmt"
	"io"
	"text/template"
)

//go:embed templates/wither.go.tmpl
var witherTemplateString string
var witherTemplate *template.Template

func init() {
	witherTemplate = loadTemplate("Wither", witherTemplateString)
}

func processWither(data *typeProcessorData) error {
	commands, found := hasComment(data.typeComments, "Wither")
	if found {
		config, err := parseWitherConfig(commands, data.structName)
		if err != nil {
			return err
		}
		for _, fieldName := range data.fieldNames {
			typeName := data.fields[fieldName]
			debugLogger.Printf("Generating Wither for %s.%s", data.structName, fieldName)
			fieldCommands, found := hasComment(data.fieldComments[fieldName], "Wither")
			if found {
				config, err := parseWitherConfig(fieldCommands, fmt.Sprintf("%s.%s", data.structName, fieldName))
				if err != nil {
					return err
				}
				addWither(fieldName, typeName, data, config)
			} else {
				addWither(fieldName, typeName, data, config)
			}
		}
	} else {
		for _, fieldName := range data.fieldNames {
			typeName := data.fields[fieldName]
			commands, found := hasComment(data.fieldComments[fieldName], "Wither")
			if found {
				config, err := parseWitherConfig(commands, fmt.Sprintf("%s.%s", data.structName, fieldName))
				if err != nil {
					return err
				}
				debugLogger.Printf("Generating Wither for %s.%s", data.structName, fieldName)
				addWither(fieldName, typeName, data, config)
			}
		}
	}
	return nil
}

type witherConfig struct {
	private bool
}

func parseWitherConfig(commands []string, structName string) (*witherConfig, error) {
	private := false
	privateSet := false
	for _, command := range commands {
		switch command {
		case "private":
			if privateSet {
				return nil, fmt.Errorf("Invalid command %s on %s Wither", command, structName)
			}
			private = true
			privateSet = true
		case "exported":
			if privateSet {
				return nil, fmt.Errorf("Invalid command %s on %s Wither", command, structName)
			}
			private = false
			privateSet = true
		default:
			return nil, fmt.Errorf("Invalid command %s on %s Wither", command, structName)
		}
	}
	return &witherConfig{
		private: private,
	}, nil
}

func addWither(fieldName, fieldType string, data *typeProcessorData, config *witherConfig) {
	w := "W"
	if config.private {
		w = "w"
	}
	data.addCodeWriter(func(wr io.Writer) error {
		return witherTemplate.Execute(wr, map[string]interface{}{
			"structName":       data.structName,
			"fieldName":        fieldName,
			"fieldType":        fieldType,
			"fieldNames":       data.fieldNames,
			"fields":           data.fields,
			"genericTypes":     data.genericTypes,
			"genericTypeNames": data.genericTypeNames,
			"w":                w,
		})
	})
}
