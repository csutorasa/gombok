package main

import (
	"fmt"
	"io"
	"text/template"
)

var getterTemplateString string = `
func (this *{{ .structName }}{{ genericList .genericTypeNames }}) {{ .g }}et{{ capitalize .fieldName }}() {{ .fieldType }} {
	return this.{{ .fieldName }}
}
`
var getterTemplate *template.Template

func init() {
	getterTemplate = loadTemplate("Getter", getterTemplateString)
}

func processGetter(data *typeProcessorData) error {
	commands, found := hasComment(data.typeComments, "Getter")
	if found {
		config, err := parseGetterConfig(commands, data.structName)
		if err != nil {
			return err
		}
		for fieldName, typeName := range data.fields {
			debugLogger.Printf("Generating Getter for %s.%s", data.structName, fieldName)
			fieldCommands, found := hasComment(data.fieldComments[fieldName], "Getter")
			if found {
				config, err := parseGetterConfig(fieldCommands, fmt.Sprintf("%s.%s", data.structName, fieldName))
				if err != nil {
					return err
				}
				addGetter(fieldName, typeName, data, config)
			} else {
				addGetter(fieldName, typeName, data, config)
			}
		}
	} else {
		for fieldName, typeName := range data.fields {
			commands, found := hasComment(data.fieldComments[fieldName], "Getter")
			if found {
				config, err := parseGetterConfig(commands, fmt.Sprintf("%s.%s", data.structName, fieldName))
				if err != nil {
					return err
				}
				debugLogger.Printf("Generating Getter for %s.%s", data.structName, fieldName)
				addGetter(fieldName, typeName, data, config)
			}
		}
	}
	return nil
}

type getterConfig struct {
	private bool
}

func parseGetterConfig(commands []string, structName string) (*getterConfig, error) {
	private := false
	privateSet := false
	for _, command := range commands {
		switch command {
		case "private":
			if privateSet {
				return nil, fmt.Errorf("Invalid command %s on %s Getter", command, structName)
			}
			private = true
			privateSet = true
		case "exported":
			if privateSet {
				return nil, fmt.Errorf("Invalid command %s on %s Getter", command, structName)
			}
			private = false
			privateSet = true
		default:
			return nil, fmt.Errorf("Invalid command %s on %s Getter", command, structName)
		}
	}
	return &getterConfig{
		private: private,
	}, nil
}

func addGetter(fieldName string, fieldType string, data *typeProcessorData, config *getterConfig) {
	g := "G"
	if config.private {
		g = "g"
	}
	data.addCodeWriter(func(wr io.Writer) error {
		return getterTemplate.Execute(wr, map[string]interface{}{
			"structName":       data.structName,
			"fieldName":        fieldName,
			"fieldType":        fieldType,
			"genericTypes":     data.genericTypes,
			"genericTypeNames": data.genericTypeNames,
			"g":                g,
		})
	})
}
