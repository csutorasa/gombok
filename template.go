package main

import (
	"fmt"
	"strings"
	"text/template"
)

var templateFunctions template.FuncMap = template.FuncMap{
	"capitalize":           capitalize,
	"lower":                lower,
	"genericList":          genericList,
	"genericListWithTypes": genericListWithTypes,
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	if len(s) == 1 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func lower(s string) string {
	if s == "" {
		return ""
	}
	if len(s) == 1 {
		return strings.ToLower(s)
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func genericList(genericTypeNames []string) string {
	if len(genericTypeNames) == 0 {
		return ""
	}
	return fmt.Sprintf("[%s]", strings.Join(genericTypeNames, ", "))
}

func genericListWithTypes(genericTypeNames []string, genericTypes map[string]string) string {
	if len(genericTypeNames) == 0 {
		return ""
	}
	n := []string{}
	for _, name := range genericTypeNames {
		n = append(n, fmt.Sprintf("%s %s", name, genericTypes[name]))
	}
	return fmt.Sprintf("[%s]", strings.Join(n, ", "))
}

func loadTemplate(name, text string) *template.Template {
	t, err := template.New(name).Funcs(templateFunctions).Parse(text)
	if err != nil {
		panic(err)
	}
	return t
}
