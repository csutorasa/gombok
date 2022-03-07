package main

import (
	"testing"
)

func TestCreateDefaultConstructor(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").addField("finished", "bool").addComment("Constructor")

	err := processConstructor(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func NewMystruct(finished bool) *mystruct {
	return &mystruct{
		finished: finished,
	}
}
`)
}

func TestCreateSingleConstructor(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addField("finished", "bool").
		addFieldWithComment("count", "int", "Constructor exclude").
		addComment("Constructor")

	err := processConstructor(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func NewMystruct(finished bool) *mystruct {
	return &mystruct{
		finished: finished,
	}
}
`)
}

func TestCreateAllConstructors(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addField("finished", "bool").
		addField("count", "int").
		addComment("Constructor")

	err := processConstructor(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func NewMystruct(finished bool, count int) *mystruct {
	return &mystruct{
		finished: finished,
		count: count,
	}
}
`)
}
