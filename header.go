package main

import (
	_ "embed"
	"io"
	"text/template"
)

//go:embed templates/header.go.tmpl
var headerTemplateString string
var headerTemplate *template.Template

//go:embed templates/import.go.tmpl
var importTemplateString string
var importTemplate *template.Template

//go:embed templates/imports.go.tmpl
var importsTemplateString string
var importsTemplate *template.Template

const header string = "// Autogenerated file by gombok"

func init() {
	headerTemplate = loadTemplate("Header", headerTemplateString)
	importTemplate = loadTemplate("Import", importTemplateString)
	importsTemplate = loadTemplate("Imports", importsTemplateString)
}

func writeHeader(wr io.Writer, packageName string) error {
	return headerTemplate.Execute(wr, map[string]interface{}{
		"packageName": packageName,
		"header":      header,
		"ignore":      ignore,
	})
}

func writeImport(wr io.Writer, imports map[string]bool, fileImports map[string]*impData) error {
	if len(imports) == 1 {
		for imp := range imports {
			fileImp := fileImports[imp]
			if imp == "fmt" && fileImp == nil {
				fileImp = &impData{
					HasName: false,
					Name:    "fmt",
					Path:    "fmt",
				}
			}
			return importTemplate.Execute(wr, map[string]interface{}{
				"HasName": fileImp.HasName,
				"Name":    fileImp.Name,
				"Path":    fileImp.Path,
			})
		}
	} else if len(imports) > 1 {
		imps := map[string]*impData{}
		for imp := range imports {
			fileImp := fileImports[imp]
			if imp == "fmt" && fileImp == nil {
				imps["fmt"] = &impData{
					HasName: false,
					Name:    "fmt",
					Path:    "fmt",
				}
			} else {
				imps[imp] = fileImp
			}
		}
		return importsTemplate.Execute(wr, map[string]interface{}{
			"imports": imps,
		})
	}
	return nil
}
