package main

import (
	"bytes"
	"strings"
	"testing"
)

type testTypeProcessorData struct {
	imports           []string
	codeparts         []string
	errors            []string
	typeProcessorData *typeProcessorData
}

func createTestTypeProcessorData(packageName string, structName string, fields map[string]string, fieldComments map[string][]string, fieldNames []string, typeComments []string) *testTypeProcessorData {
	testData := &testTypeProcessorData{
		imports:   []string{},
		codeparts: []string{},
		errors:    []string{},
		typeProcessorData: &typeProcessorData{
			packageName:   packageName,
			structName:    structName,
			fields:        fields,
			fieldComments: fieldComments,
			fieldNames:    fieldNames,
			typeComments:  typeComments,
		},
	}
	testData.typeProcessorData.addImport = func(s string) {
		testData.imports = append(testData.imports, s)
	}
	testData.typeProcessorData.addCodeWriter = func(cw codeWriter) {
		var wr bytes.Buffer
		err := cw(&wr)
		if err != nil {
			testData.errors = append(testData.errors, err.Error())
		} else {
			testData.codeparts = append(testData.codeparts, wr.String())
		}
	}
	return testData
}

func createEmptyTestData(packageName string, structName string) *testTypeProcessorData {
	return createTestTypeProcessorData(
		packageName,
		structName,
		map[string]string{},
		map[string][]string{},
		[]string{},
		[]string{},
	)
}

func (d *testTypeProcessorData) addField(name string, typeName string) *testTypeProcessorData {
	d.typeProcessorData.fields[name] = typeName
	d.typeProcessorData.fieldNames = append(d.typeProcessorData.fieldNames, name)
	return d
}

func (d *testTypeProcessorData) addFieldComment(name string, comments ...string) *testTypeProcessorData {
	curr, ok := d.typeProcessorData.fieldComments[name]
	if ok {
		d.typeProcessorData.fieldComments[name] = append(curr, comments...)
	} else {
		d.typeProcessorData.fieldComments[name] = comments
	}
	return d
}

func (d *testTypeProcessorData) addFieldWithComment(name string, typeName string, comments ...string) *testTypeProcessorData {
	d.addField(name, typeName)
	d.addFieldComment(name, comments...)
	return d
}

func (d *testTypeProcessorData) addComment(comments ...string) *testTypeProcessorData {
	d.typeProcessorData.typeComments = append(d.typeProcessorData.typeComments, comments...)
	return d
}

func (d *testTypeProcessorData) validateErrors(t *testing.T, errors ...string) {
	validate(t, "error", d.errors, errors)
}

func (d *testTypeProcessorData) validateImports(t *testing.T, imports ...string) {
	validate(t, "import", d.imports, imports)
}

func (d *testTypeProcessorData) validateCodeparts(t *testing.T, codeparts ...string) {
	validate(t, "code", d.codeparts, codeparts)
}

func validate(t *testing.T, typeName string, actual []string, expected []string) {
	for _, a := range actual {
		if !contains(a, expected) {
			t.Errorf("Found an unexcpected %s %s", a, typeName)
		}
	}
	for _, e := range expected {
		if !contains(e, actual) {
			t.Errorf("Missing an excpected %s %s", e, typeName)
		}
	}
}

func contains(item string, list []string) bool {
	for _, e := range list {
		if endlineIgnoreEquals(e, item) {
			return true
		}
	}
	return false
}

func endlineIgnoreEquals(a string, b string) bool {
	return strings.ReplaceAll(a, "\r\n", "\n") == strings.ReplaceAll(b, "\r\n", "\n")
}
