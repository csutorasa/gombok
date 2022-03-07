package main

import (
	"testing"
)

func TestCreateDefaultStringer(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").addField("finished", "bool").addComment("Stringer")

	err := processStringer(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t, "fmt")
	testData.validateCodeparts(t, `
func (this *mystruct) String() string {
	return fmt.Sprintf("mypackage.mystruct{finished=%v}", this.finished)
}
`)
}

func TestCreateSingleStringer(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addField("finished", "bool").
		addFieldWithComment("count", "int", "Stringer exclude").
		addComment("Stringer")

	err := processStringer(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t, "fmt")
	testData.validateCodeparts(t, `
func (this *mystruct) String() string {
	return fmt.Sprintf("mypackage.mystruct{finished=%v}", this.finished)
}
`)
}

func TestCreateAllStringers(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addField("finished", "bool").
		addField("count", "int").
		addComment("Stringer")

	err := processStringer(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t, "fmt")
	testData.validateCodeparts(t, `
func (this *mystruct) String() string {
	return fmt.Sprintf("mypackage.mystruct{finished=%v, count=%v}", this.finished, this.count)
}
`)
}
