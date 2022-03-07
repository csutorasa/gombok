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