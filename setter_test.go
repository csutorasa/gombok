package main

import (
	"testing"
)

func TestCreateDefaultSetter(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").addFieldWithComment("finished", "bool", "Setter")

	err := processSetter(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) SetFinished(finished bool) {
	this.finished = finished
}
`)
}

func TestCreatePrivateSetter(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").addFieldWithComment("finished", "bool", "Setter private")

	err := processSetter(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) setFinished(finished bool) {
	this.finished = finished
}
`)
}

func TestCreateExportedSetter(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").addFieldWithComment("finished", "bool", "Setter exported")

	err := processSetter(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) SetFinished(finished bool) {
	this.finished = finished
}
`)
}

func TestCreateSingleSetter(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addFieldWithComment("finished", "bool", "Setter exported").
		addField("count", "int")

	err := processSetter(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) SetFinished(finished bool) {
	this.finished = finished
}
`)
}

func TestCreateAllSetters(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addField("finished", "bool").
		addField("count", "int").
		addComment("Setter")

	err := processSetter(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) SetFinished(finished bool) {
	this.finished = finished
}
`, `
func (this *mystruct) SetCount(count int) {
	this.count = count
}
`)
}

func TestCreateMixedVisibilitySetters(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addField("finished", "bool").
		addFieldWithComment("count", "int", "Setter private").
		addComment("Setter")

	err := processSetter(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) SetFinished(finished bool) {
	this.finished = finished
}
`, `
func (this *mystruct) setCount(count int) {
	this.count = count
}
`)
}
