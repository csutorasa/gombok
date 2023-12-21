package main

import (
	"fmt"
	"io"
	"os"
)

const ignore string = "//go:gombok ignore"

type typeProcessorData struct {
	packageName      string
	structName       string
	fields           map[string]string
	fieldComments    map[string][]string
	fieldNames       []string
	typeComments     []string
	genericTypes     map[string]string
	genericTypeNames []string
	addImport        func(string)
	addCodeWriter    func(codeWriter)
}

type typeProcessor func(data *typeProcessorData) error

type impData struct {
	HasName bool
	Name    string
	Path    string
}

type codeWriter func(w io.Writer) error

type fileWriter struct {
	file        string
	pkg         string
	imports     map[string]bool
	codeWriters []codeWriter
}

func NewFileWriter(pkg, file string) *fileWriter {
	return &fileWriter{
		pkg:         pkg,
		file:        file,
		imports:     map[string]bool{},
		codeWriters: []codeWriter{},
	}
}

func (w *fileWriter) getImports() []string {
	imports := make([]string, len(w.imports))
	i := 0
	for imp := range w.imports {
		imports[i] = imp
		i = i + 1
	}
	return imports
}

func getGeneratedFileName(path string) string {
	return fmt.Sprintf(filenameFormatter, path)
}

func (w *fileWriter) Write(fileImports map[string]*impData) error {
	if len(w.imports) == 0 && len(w.codeWriters) == 0 {
		return nil
	}
	filePath := getGeneratedFileName(w.file)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	return w.WriteTo(f, fileImports)
}

func (w *fileWriter) WriteTo(wr io.Writer, fileImports map[string]*impData) error {
	err := writeHeader(wr, w.pkg)
	if err != nil {
		return err
	}
	err = writeImport(wr, w.imports, fileImports)
	if err != nil {
		return err
	}
	for _, cw := range w.codeWriters {
		err = cw(wr)
		if err != nil {
			return err
		}
	}
	return nil
}
