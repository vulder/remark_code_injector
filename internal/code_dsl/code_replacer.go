package code_dsl

import (
	"log"
	"strings"
)

//===----------------------------------------------------------------------===//
// DSL
//
// line_num      = { digit };
// range         = { digit }, "-", { digit };
// ln_range_list = range | line_num, [ { "," , range | line_num } ];
// vis_select    = "<", ln_range_list , ">";
// hl_select     = "{" , ln_range_list , "}";
//
//===----------------------------------------------------------------------===//
// Commands:
//  * "insert_code(filename:" , range | line_num , ")" , vis_select , hl_select
//===----------------------------------------------------------------------===//

// Checks if the line contains an code DSL command.
func ContainsDSLCommand(line string) bool {
	if isInsertCode(line) {
		return true
	}
	return false
}

// Transforms a line by replacing the DSL specific part with the generated
// content.
func TransformLine(line string) string {
	if isInsertCode(line) {
		return handleInsertCode(line)
	}
	log.Fatal("Transform was called without a transformable line.")
	return line
}

func wrapWithCodeBlock(text string, lang string) string {
	return "```" + lang + "\n" + text + "\n```"
}

//

//===----------------------------------------------------------------------===//
// insert_code
//
// Examples usage:
//   insert_code(filename.cpp:4-17)<4-8,17>{5-6,8}

func isInsertCode(line string) bool {
	return strings.HasPrefix(line, "insert_code")
}

func handleInsertCode(line string) string {
	ci := parseInsertCode(line)

	code_language_type := ci.progLang

	return wrapWithCodeBlock(ci.renderCodeBlock(), code_language_type)
}
