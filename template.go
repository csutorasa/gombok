package main

import (
	"fmt"
	"strings"
	"text/template"
)

var templateFunctions template.FuncMap = template.FuncMap{
	"capitalize":           capitalize,
	"lower":                lower,
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
