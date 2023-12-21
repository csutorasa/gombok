package main

import (
	_ "embed"
	"fmt"
	"io"
	"text/template"
)

//go:embed templates/setter.go.tmpl
var setterTemplateString string
var setterTemplate *template.Template

//go:embed templates/chainsetter.go.tmpl
var chainSetterTemplateString string
var chainSetterTemplate *template.Template

func init() {
	setterTemplate = loadTemplate("Setter", setterTemplateString)
	chainSetterTemplate = loadTemplate("ChainSetter", chainSetterTemplateString)
}

func processSetter(data *typeProcessorData) error {
	commands, found := hasComment(data.typeComments, "Setter")
	if found {
		config, err := parseSetterConfig(commands, data.structName)
		if err != nil {
			return err
		}
		for _, fieldName := range data.fieldNames {
			typeName := data.fields[fieldName]
			debugLogger.Printf("Generating Setter for %s.%s", data.structName, fieldName)
			fieldCommands, found := hasComment(data.fieldComments[fieldName], "Setter")
			if found {
				config, err := parseSetterConfig(fieldCommands, fmt.Sprintf("%s.%s", data.structName, fieldName))
				if err != nil {
					return err
				}
				addSetter(fieldName, typeName, data, config)
			} else {
				addSetter(fieldName, typeName, data, config)
			}
		}
	} else {
		for _, fieldName := range data.fieldNames {
			typeName := data.fields[fieldName]
			commands, found := hasComment(data.fieldComments[fieldName], "Setter")
			if found {
				config, err := parseSetterConfig(commands, fmt.Sprintf("%s.%s", data.structName, fieldName))
				if err != nil {
					return err
				}
				debugLogger.Printf("Generating Setter for %s.%s", data.structName, fieldName)
				addSetter(fieldName, typeName, data, config)
			}
		}
	}
	return nil
}

type setterConfig struct {
	private bool
	chain   bool
}

func parseSetterConfig(commands []string, structName string) (*setterConfig, error) {
	private := false
	privateSet := false
	chain := false
	chainSet := false
	for _, command := range commands {
		switch command {
		case "private":
			if privateSet {
				return nil, fmt.Errorf("invalid command %s on %s Setter", command, structName)
			}
			private = true
			privateSet = true
		case "exported":
			if privateSet {
				return nil, fmt.Errorf("invalid command %s on %s Setter", command, structName)
			}
			private = false
			privateSet = true
		case "chained":
			if chainSet {
				return nil, fmt.Errorf("invalid command %s on %s Setter", command, structName)
			}
			chain = true
			chainSet = true
		default:
			return nil, fmt.Errorf("invalid command %s on %s Setter", command, structName)
		}
	}
	return &setterConfig{
		private: private,
		chain:   chain,
	}, nil
}

func addSetter(fieldName, fieldType string, data *typeProcessorData, config *setterConfig) {
	s := "S"
	if config.private {
		s = "s"
	}
	data.addCodeWriter(func(wr io.Writer) error {
		var t *template.Template = nil
		if config.chain {
			t = chainSetterTemplate
		} else {
			t = setterTemplate
		}
		return t.Execute(wr, map[string]interface{}{
			"structName":       data.structName,
			"fieldName":        fieldName,
			"fieldType":        fieldType,
			"genericTypes":     data.genericTypes,
			"genericTypeNames": data.genericTypeNames,
			"s":                s,
		})
	})
}
