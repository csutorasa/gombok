package main

import (
	"testing"
)

func TestCapitalize(t *testing.T) {
	result := capitalize("testName")
	if result != "TestName" {
		t.Errorf("Result should be TestName instead of %s", result)
	}
	result = capitalize(result)
	if result != "TestName" {
		t.Errorf("Result should be TestName instead of %s", result)
	}
}

func TestLower(t *testing.T) {
	result := lower("TestName")
	if result != "testName" {
		t.Errorf("Result should be testName instead of %s", result)
	}
	result = lower(result)
	if result != "testName" {
		t.Errorf("Result should be testName instead of %s", result)
	}
}

func TestGenericList(t *testing.T) {
	result := genericList([]string{"S", "T"})
	if result != "[S, T]" {
		t.Errorf("Result should be [S, T] instead of %s", result)
	}
}

func TestGenericListWithTypes(t *testing.T) {
	result := genericListWithTypes([]string{"S", "T"}, map[string]string{"S": "bool", "T": "int"})
	if result != "[S bool, T int]" {
		t.Errorf("Result should be [S bool, T int] instead of %s", result)
	}
}
