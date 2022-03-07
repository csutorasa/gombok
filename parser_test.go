package main

import "testing"

func TestHasComment(t *testing.T) {
	commands, found := hasComment([]string{"ignore"}, "ignore")
	if len(commands) != 0 || !found {
		t.Errorf("ignore should be found without commands")
	}
}

func TestHasNoComment(t *testing.T) {
	commands, found := hasComment([]string{"random:ignore"}, "ignore")
	if len(commands) != 0 || found {
		t.Errorf("ignore should not be found without commands")
	}
}
