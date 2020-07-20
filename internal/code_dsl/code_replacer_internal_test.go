package code_dsl

import (
	"testing"
)

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
