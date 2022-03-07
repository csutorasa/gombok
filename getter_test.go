package main

import (
	"testing"
)

func TestCreateDefaultGetter(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").addFieldWithComment("finished", "bool", "Getter")

	err := processGetter(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) GetFinished() bool {
	return this.finished
}
`)
}

func TestCreatePrivateGetter(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").addFieldWithComment("finished", "bool", "Getter private")

	err := processGetter(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) getFinished() bool {
	return this.finished
}
`)
}

func TestCreateExportedGetter(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").addFieldWithComment("finished", "bool", "Getter exported")

	err := processGetter(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) GetFinished() bool {
	return this.finished
}
`)
}

func TestCreateSingleGetter(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addFieldWithComment("finished", "bool", "Getter exported").
		addField("count", "int")

	err := processGetter(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) GetFinished() bool {
	return this.finished
}
`)
}

func TestCreateAllGetters(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addField("finished", "bool").
		addField("count", "int").
		addComment("Getter")

	err := processGetter(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) GetFinished() bool {
	return this.finished
}
`, `
func (this *mystruct) GetCount() int {
	return this.count
}
`)
}

func TestCreateMixedVisibilityGetters(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addField("finished", "bool").
		addFieldWithComment("count", "int", "Getter private").
		addComment("Getter")

	err := processGetter(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) GetFinished() bool {
	return this.finished
}
`, `
func (this *mystruct) getCount() int {
	return this.count
}
`)
}
