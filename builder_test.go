package main

import (
	"testing"
)

func TestCreateDefaultBuilder(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").addField("finished", "bool").addComment("Builder")

	err := processBuilder(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
type mystructBuilder struct {
	data mystruct
}

func (this *mystructBuilder) Finished(finished bool) *mystructBuilder {
	this.data.finished = finished
	return this
}

func (this *mystructBuilder) Build() mystruct {
	return this.data
}

func NewMystructBuilder() *mystructBuilder {
	return &mystructBuilder{
		data: mystruct{},
	}
}
`)
}

func TestCreateSingleBuilder(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addField("finished", "bool").
		addFieldWithComment("count", "int", "Builder exclude").
		addComment("Builder")

	err := processBuilder(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
type mystructBuilder struct {
	data mystruct
}

func (this *mystructBuilder) Finished(finished bool) *mystructBuilder {
	this.data.finished = finished
	return this
}

func (this *mystructBuilder) Build() mystruct {
	return this.data
}

func NewMystructBuilder() *mystructBuilder {
	return &mystructBuilder{
		data: mystruct{},
	}
}
`)
}

func TestCreateAllBuilders(t *testing.T) {
	testData := createEmptyTestData("mypackage", "mystruct").
		addField("finished", "bool").
		addField("count", "int").
		addComment("Builder")

	err := processBuilder(testData.typeProcessorData)
	if err != nil {
		t.Errorf("Should not fail with %s", err.Error())
	}
	testData.validateErrors(t)
	testData.validateImports(t)
	testData.validateCodeparts(t, `
type mystructBuilder struct {
	data mystruct
}

func (this *mystructBuilder) Finished(finished bool) *mystructBuilder {
	this.data.finished = finished
	return this
}

func (this *mystructBuilder) Count(count int) *mystructBuilder {
	this.data.count = count
	return this
}

func (this *mystructBuilder) Build() mystruct {
	return this.data
}

func NewMystructBuilder() *mystructBuilder {
	return &mystructBuilder{
		data: mystruct{},
	}
}
`)
}
