package code_dsl

import (
	"log"
	"strings"
)

//===----------------------------------------------------------------------===//
// DSL
//
// line_num         = { digit };
// range            = { digit }, "-", { digit };
// char_range       = { digit }, "|", { digit };
// char_range_list  = char_range | [ { "," , char_range } ];
// ln_range_list    = range | line_num, [ { "," , range | line_num } ];
// vis_select       = "r", "<", ["h" | "d" | "r"], ln_range_list , ">";
// hl_select        = "r", "{" , ln_range_list , "}";
//
//===----------------------------------------------------------------------===//
// Commands:
//  * "insert_code(filename:" , range | line_num , ")" , vis_select , hl_select
//  * "rev_insert_code(filename:BlockID)" , vis_select , hl_select
//===----------------------------------------------------------------------===//

// Checks if the line contains an code DSL command.
func ContainsDSLCommand(line string) bool {
	if isInsertCode(line) {
		return true
	}
	if isRevInsertCode(line) {
		return true
	}
	return false
}

// Transforms a line by replacing the DSL specific part with the generated
// content.
func TransformLine(line string, codeRoot string) string {
	if isInsertCode(line) {
		return handleInsertCode(line, codeRoot)
	}
	if isRevInsertCode(line) {
		return handleRevInsertCode(line, codeRoot)
	}
	log.Fatal("Transform was called without a transformable line.")
	return line
}

func wrapWithCodeBlock(text string, lang string) string {
	return "```" + lang + "\n" + text + "```"
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

func handleInsertCode(line string, codeRoot string) string {
	ci, err := parseInsertCode(line, codeRoot)
	if err != nil { // In the error case we return the unprocessed line to not destroy the doc.
		print("Could not process insert_code line: ", line)
		return line
	}

	codeLanguageType := ci.progLang

	return wrapWithCodeBlock(ci.renderCodeBlock(), codeLanguageType)
}

//===----------------------------------------------------------------------===//
// rev_insert_code
//
// Examples usage:
//   rev_insert_code(filename.cpp:ExampleID)<4-8,17>{5-6,8}

func isRevInsertCode(line string) bool {
	return strings.HasPrefix(line, "rev_insert_code")
}

func handleRevInsertCode(line string, codeRoot string) string {
	ci, err := parseRevInsertCode(line, codeRoot)
	if err != nil { // In the error case we return the unprocessed line to not destroy the doc.
		print("Could not process rev_insert_code line: ", line)
		return line
	}

	return wrapWithCodeBlock(ci.renderCodeBlock(), ci.progLang)
}
