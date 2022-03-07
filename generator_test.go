package main

import "testing"

func TestFileWriterGetImports(t *testing.T) {
	writer := NewFileWriter("mypacakge", "testfile")
	writer.imports["fmt"] = false

	imports := writer.getImports()
	if len(imports) != 1 || imports[0] != "fmt" {
		t.Errorf("fmt should be the one and only import instead of %v", imports)
	}
}

func TestDefaultFilenameGenerator(t *testing.T) {
	filename := getGeneratedFileName("testfile")
	if filename != "testfile_gombok.go" {
		t.Errorf("Filename should be testfile_gombok.go instead of %s", filename)
	}
}
