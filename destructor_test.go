package main

import (
	"testing"
)

func TestCreateDefaultDestructor(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").addField("finished", "bool").addComment("Destructor")

	err := processDestructor(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) Destruct() (bool) {
	return this.finished
}
`)
}

func TestCreatePrivateDestructor(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").addField("finished", "bool").addComment("Destructor private")

	err := processDestructor(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) destruct() (bool) {
	return this.finished
}
`)
}

func TestCreateExportedDestructor(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").addField("finished", "bool").addComment("Destructor exported")

	err := processDestructor(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) Destruct() (bool) {
	return this.finished
}
`)
}

func TestCreateSingleDestructor(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addField("finished", "bool").
		addFieldWithComment("count", "int", "Destructor exclude").
		addComment("Destructor")

	err := processDestructor(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) Destruct() (bool) {
	return this.finished
}
`)
}

func TestCreateAllDestructors(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addField("finished", "bool").
		addField("count", "int").
		addComment("Destructor")

	err := processDestructor(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) Destruct() (bool, int) {
	return this.finished, this.count
}
`)
}
