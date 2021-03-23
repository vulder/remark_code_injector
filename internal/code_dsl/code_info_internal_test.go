package code_dsl

import (
	"testing"
)

func TestParseCodeGenOptionsEmptyString(t *testing.T) {
	optionString := ""

	options := ParseCodeGenOptions(optionString)

	if options.getIndent() != 0 {
		t.Log("Indent was set to", options.getIndent(), " but expected", 0)
		t.Error("Indent was wrong.")
	}

	if options.hideComments() != false {
		t.Log("Remove comment was set to", options.hideComments(), " but expected", false)
		t.Error("Hide comments was wrong.")
	}
}

func TestParseCodeGenOptionsComments(t *testing.T) {
	optionString := "comments=true"

	options := ParseCodeGenOptions(optionString)

	if options.getIndent() != 0 {
		t.Log("Indent was set to", options.getIndent(), " but expected", 0)
		t.Error("Indent was wrong.")
	}

	if options.hideComments() != false {
		t.Log("Remove comment was set to", options.hideComments(), " but expected", false)
		t.Error("Hide comments was wrong.")
	}
}

func TestParseCodeGenOptionsCommentsUpperCase(t *testing.T) {
	optionString := "comments=True"

	options := ParseCodeGenOptions(optionString)

	if options.getIndent() != 0 {
		t.Log("Indent was set to", options.getIndent(), " but expected", 0)
		t.Error("Indent was wrong.")
	}

	if options.hideComments() != false {
		t.Log("Remove comment was set to", options.hideComments(), " but expected", false)
		t.Error("Hide comments was wrong.")
	}
}

func TestParseCodeGenOptionsPosIndent(t *testing.T) {
	optionString := "indent=2"

	options := ParseCodeGenOptions(optionString)

	if options.getIndent() != 2 {
		t.Log("Indent was set to", options.getIndent(), " but expected", 2)
		t.Error("Indent was wrong.")
	}

	if options.hideComments() != false {
		t.Log("Remove comment was set to", options.hideComments(), " but expected", false)
		t.Error("Hide comments was wrong.")
	}
}

func TestParseCodeGenOptionsNegIndent(t *testing.T) {
	optionString := "indent=-2"

	options := ParseCodeGenOptions(optionString)

	if options.getIndent() != -2 {
		t.Log("Indent was set to", options.getIndent(), " but expected", -2)
		t.Error("Indent was wrong.")
	}

	if options.hideComments() != false {
		t.Log("Remove comment was set to", options.hideComments(), " but expected", false)
		t.Error("Hide comments was wrong.")
	}
}

func TestParseCodeGenOptionsIgnoreFalseOptions(t *testing.T) {
	optionString := "key=value,key2=bar"

	options := ParseCodeGenOptions(optionString)

	if options.getIndent() != 0 {
		t.Log("Indent was set to", options.getIndent(), " but expected", 0)
		t.Error("Indent was wrong.")
	}

	if options.hideComments() != false {
		t.Log("Remove comment was set to", options.hideComments(), " but expected", false)
		t.Error("Hide comments was wrong.")
	}
}

func TestParseCodeGenOptionsMultiple(t *testing.T) {
	optionString := "indent=2,comments=false"

	options := ParseCodeGenOptions(optionString)

	if options.getIndent() != 2 {
		t.Log("Indent was set to", options.getIndent(), " but expected", 2)
		t.Error("Indent was wrong.")
	}

	if options.hideComments() != true {
		t.Log("Remove comment was set to", options.hideComments(), " but expected", true)
		t.Error("Hide comments was wrong.")
	}
}
