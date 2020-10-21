package code_dsl

import (
	"testing"
)

func TestAllsDSLCommands(t *testing.T) {
	goodCommands := [...]string{"insert_code(foo.cpp:1-2)", "rev_insert_code(foo.cpp:ExampleID)"}
	for _, command := range goodCommands {
		if !ContainsDSLCommand(command) {
			t.Errorf("Command %s was not correctly indentified as a DSL command", command)
		}
	}
}

func TestWrongDSLCommand(t *testing.T) {
	if ContainsDSLCommand("some other text") {
		t.Error("Non DSL command was falsly detected as a DSL command.")
	}
}

func TestIsInsertCodeDetectKeyword(t *testing.T) {
	if !isInsertCode("insert_code") {
		t.Error("Insert code keyword not recognized")
	}
}

func TestIsInsertCodeOtherWord(t *testing.T) {
	if isInsertCode("other words") {
		t.Error("Insert code keyword recognized but none present")
	}
}

func TestIsRevInsertCodeDetectKeyword(t *testing.T) {
	if !isRevInsertCode("rev_insert_code") {
		t.Error("Insert code keyword not recognized")
	}
}

func TestIsRevInsertCodeOtherWord(t *testing.T) {
	if isRevInsertCode("other words") {
		t.Error("Insert code keyword recognized but none present")
	}
}
