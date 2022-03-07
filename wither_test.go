package main

import (
	"testing"
)

func TestCreateDefaultWither(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").addFieldWithComment("finished", "bool", "Wither")

	err := processWither(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) WithFinished(finished bool) *mystruct {
	return &mystruct{
		finished: finished,
	}
}
`)
}

func TestCreatePrivateWither(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").addFieldWithComment("finished", "bool", "Wither private")

	err := processWither(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) withFinished(finished bool) *mystruct {
	return &mystruct{
		finished: finished,
	}
}
`)
}

func TestCreateExportedWither(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").addFieldWithComment("finished", "bool", "Wither exported")

	err := processWither(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) WithFinished(finished bool) *mystruct {
	return &mystruct{
		finished: finished,
	}
}
`)
}

func TestCreateSingleWither(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addFieldWithComment("finished", "bool", "Wither exported").
		addField("count", "int")

	err := processWither(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) WithFinished(finished bool) *mystruct {
	return &mystruct{
		finished: finished,
		count: this.count,
	}
}
`)
}

func TestCreateAllWithers(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addField("finished", "bool").
		addField("count", "int").
		addComment("Wither")

	err := processWither(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) WithFinished(finished bool) *mystruct {
	return &mystruct{
		finished: finished,
		count: this.count,
	}
}
`, `
func (this *mystruct) WithCount(count int) *mystruct {
	return &mystruct{
		finished: this.finished,
		count: count,
	}
}
`)
}

func TestCreateMixedVisibilityWithers(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addField("finished", "bool").
		addFieldWithComment("count", "int", "Wither private").
		addComment("Wither")

	err := processWither(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
func (this *mystruct) WithFinished(finished bool) *mystruct {
	return &mystruct{
		finished: finished,
		count: this.count,
	}
}
`, `
func (this *mystruct) withCount(count int) *mystruct {
	return &mystruct{
		finished: this.finished,
		count: count,
	}
}
`)
}
