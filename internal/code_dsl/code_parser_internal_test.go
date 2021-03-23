package code_dsl

import (
	"github.com/Flaque/filet"
	"testing"
)

func TestInsertAtInbetween(t *testing.T) {
	baseString := "123456789"
	expectedString := "1234@56789"

	insertedString := insertAt(baseString, 4, "@")

	if insertedString != expectedString {
		t.Log("insertedLine: ", insertedString, " but expected ", expectedString)
		t.Error("Text was inserted wronly into line.")
	}
}
func TestInsertAtStart(t *testing.T) {
	baseString := "123456789"
	expectedString := "@123456789"

	insertedString := insertAt(baseString, 0, "@")

	if insertedString != expectedString {
		t.Log("insertedString: ", insertedString, " but expected ", expectedString)
		t.Error("Text was inserted wronly into line.")
	}
}
func TestInsertAtEnd(t *testing.T) {
	baseString := "123456789"
	expectedString := "123456789@"

	insertedString := insertAt(baseString, 10, "@")

	if insertedString != expectedString {
		t.Log("insertedLine: ", insertedString, " but expected ", expectedString)
		t.Error("Text was inserted wronly into line.")
	}
}

func TestAdaptIndentZero(t *testing.T) {
	baseString := "  FooBar"
	expectedString := "  FooBar"
	trimValue := 0

	trimmedString := adaptIndent(baseString, trimValue)

	if trimmedString != expectedString {
		t.Log("trimmedString:", "'"+trimmedString+"'", " but expected ", "'"+expectedString+"'")
		t.Error("Line indent was not correctly adapted.")
	}
}

func TestAdaptIndentAddLessThanPresent(t *testing.T) {
	baseString := "  FooBar"
	expectedString := "      FooBar"
	trimValue := 4

	trimmedString := adaptIndent(baseString, trimValue)

	if trimmedString != expectedString {
		t.Log("trimmedString:", "'"+trimmedString+"'", " but expected ", "'"+expectedString+"'")
		t.Error("Line indent was not correctly adapted.")
	}
}

func TestAdaptIndentMoreThanPresent(t *testing.T) {
	baseString := "  FooBar"
	expectedString := "FooBar"
	trimValue := -4

	trimmedString := adaptIndent(baseString, trimValue)

	if trimmedString != expectedString {
		t.Log("trimmedString:", "'"+trimmedString+"'", " but expected ", "'"+expectedString+"'")
		t.Error("Line indent was not correctly adapted.")
	}
}

func TestAdaptIndentRemoveLessThanPresent(t *testing.T) {
	baseString := "      FooBar"
	expectedString := "  FooBar"
	trimValue := -4

	trimmedString := adaptIndent(baseString, trimValue)

	if trimmedString != expectedString {
		t.Log("trimmedString:", "'"+trimmedString+"'", " but expected ", "'"+expectedString+"'")
		t.Error("Line indent was not correctly adapted.")
	}
}

//===----------------------------------------------------------------------===//
// Highlights

func TestRenderOneSubrange(t *testing.T) {
	baseLine := "this is a line"
	hl := Highlights{}
	hl.Init()
	hl.PushBack(CharRange{LineNumber{1}, 1, 4})

	renderedLine := hl.RenderSubrange(baseLine, 1)

	expectedLine := "`this` is a line"
	if renderedLine != expectedLine {
		t.Log("renderedLine: ", renderedLine, " but expected ", expectedLine)
		t.Error("Line was wrongly redered.")
	}
}

func TestRenderMultipleSubranges(t *testing.T) {
	baseLine := "this is a line"
	hl := Highlights{}
	hl.Init()
	hl.PushBack(CharRange{LineNumber{1}, 1, 4})
	hl.PushBack(CharRange{LineNumber{1}, 11, 14})

	renderedLine := hl.RenderSubrange(baseLine, 1)

	expectedLine := "`this` is a `line`"
	if renderedLine != expectedLine {
		t.Log("renderedLine: ", renderedLine, " but expected ", expectedLine)
		t.Error("Line was wrongly redered.")
	}
}

func TestRenderMultipleMissorderedSubranges(t *testing.T) {
	baseLine := "this is a line"
	hl := Highlights{}
	hl.Init()
	// Ordered in reverse on purpose
	hl.PushBack(CharRange{LineNumber{1}, 11, 14})
	hl.PushBack(CharRange{LineNumber{1}, 1, 4})

	renderedLine := hl.RenderSubrange(baseLine, 1)

	expectedLine := "`this` is a `line`"
	if renderedLine != expectedLine {
		t.Log("renderedLine: ", renderedLine, " but expected ", expectedLine)
		t.Error("Line was wrongly redered.")
	}
}

func TestParseCodeBlock(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/foo.cpp"
	filet.File(t, codeFilePath, `template <typename T>
T shaveTheYak(T t) {
  return t;
}

int main(/* int argc, char *argv[] */) {
  return shaveTheYak(42);
}
`)

	cb := parseCodeBlock(codeFilePath, 2, 4)

	if cb.fileRange.start != 2 {
		t.Error("CodeBlock starts at the wrong line.")
	}

	if cb.fileRange.end != 4 {
		t.Error("CodeBlock ends at the wrong line.")
	}

	renderedCode := cb.render(nil, nil, "cpp", MakeDefaultCodeGenOptions())

	if renderedCode != `T shaveTheYak(T t) {
  return t;
}
` {
		t.Log("renderedCode: ", renderedCode)
		t.Error("Parsed Code was wrongly redered.")
	}
}

//===----------------------------------------------------------------------===//
// insert_code

func TestRenderCodeBlock(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/foo.cpp"
	filet.File(t, codeFilePath, `template <typename T>
T shaveTheYak(T t) {
  return t;
}

int main(/* int argc, char *argv[] */) {
  return shaveTheYak(42);
}
`)

	ci, err := parseInsertCode("insert_code("+codeFilePath+":1-4)", "")

	renderedCode := ci.renderCodeBlock()
	t.Log(ci.renderCodeBlock())

	if renderedCode != `template <typename T>
T shaveTheYak(T t) {
  return t;
}
` || err != nil {
		t.Log("renderedCode: ", renderedCode)
		t.Error("Code was wrongly generated for `insert_code`.")
	}
}

func TestRenderHighlights(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/foo.cpp"
	filet.File(t, codeFilePath, `template <typename T>
T shaveTheYak(T t) {
  return t;
}

int main(/* int argc, char *argv[] */) {
  return shaveTheYak(42);
}
`)
	dsl_string := "insert_code(" + codeFilePath + ":1-4){4,1-2}"

	ci, err := parseInsertCode(dsl_string, "")
	ci.highlights.Init()

	parseHighlights(dsl_string, &ci.highlights, nil /*should not be needed*/)

	renderedCode := ci.renderCodeBlock()

	if renderedCode != `*template <typename T>
*T shaveTheYak(T t) {
  return t;
*}
` || err != nil {
		t.Log("renderedCode: ", renderedCode)
		t.Error("Code was wrongly generated for `insert_code` with highlights.")
	}
}

func TestRenderHighlightsRelative(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/foo.cpp"
	filet.File(t, codeFilePath, `template <typename T>
T shaveTheYak(T t) {
  return t;
}

int main(/* int argc, char *argv[] */) {
  return shaveTheYak(42);
}
`)
	dsl_string := "insert_code(" + codeFilePath + ":1-4)r{2-3}"
	ci, err := parseInsertCode(dsl_string, "")
	ci.highlights.Init()

	baseCodeRange := LineRange{1, 4}
	parseHighlights(dsl_string, &ci.highlights, &baseCodeRange)

	renderedCode := ci.renderCodeBlock()

	if renderedCode != `template <typename T>
*T shaveTheYak(T t) {
* return t;
}
` || err != nil {
		t.Log("renderedCode: ", renderedCode)
		t.Error("Code was wrongly generated for `insert_code` with highlights.")
	}
}

func TestRenderHighlightsRelativeCharRange(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/foo.cpp"
	filet.File(t, codeFilePath, `template <typename T>
T shaveTheYak(T t) {
  return t;
}

int main(/* int argc, char *argv[] */) {
  return shaveTheYak(42);
}
`)
	dsl_string := "insert_code(" + codeFilePath + ":1-4)r{1,2:{3-13|17-17},3}"
	ci, err := parseInsertCode(dsl_string, "")
	ci.highlights.Init()

	baseCodeRange := LineRange{1, 4}
	parseHighlights(dsl_string, &ci.highlights, &baseCodeRange)

	renderedCode := ci.renderCodeBlock()

	expectedCode := "*template <typename T>\n"
	expectedCode += "T `shaveTheYak`(T `t`) {\n"
	expectedCode += "* return t;\n"
	expectedCode += "}\n"

	if renderedCode != expectedCode || err != nil {
		t.Log("renderedCode: ", renderedCode, " expected: ", expectedCode)
		t.Error("Highlights were wrongly generated for `insert_code`.")
	}
}

func TestRenderDotReplacementsLines(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/foo.cpp"
	filet.File(t, codeFilePath, `template <typename T>
T shaveTheYak(T t) {
  return t;
}

int main(int argc, char *argv[]) {
  return shaveTheYak(42);
}
`)
	dsl_string := "insert_code(" + codeFilePath + ":1-8)<d2-3,d7,d6:{9-31}>"

	ci, err := parseInsertCode(dsl_string, "")
	ci.highlights.Init()

	parseHighlights(dsl_string, &ci.highlights, nil /*should not be needed*/)

	renderedCode := ci.renderCodeBlock()

	if renderedCode != `template <typename T>
  // ...
}

int main(/* ... */) {
  // ...
}
` || err != nil {
		t.Log("renderedCode: ", renderedCode)
		t.Error("Code was wrongly generated for `insert_code` with dot replacements.")
	}
}

func TestRenderHideLines(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/foo.cpp"
	filet.File(t, codeFilePath, `template <typename T>
T shaveTheYak(T t) {
  return t;
}

int main(int argc, char *argv[]) {
  return shaveTheYak(42);
}
`)
	dsl_string := "insert_code(" + codeFilePath + ":1-8)<h2-3,h7,h6:{9-31}>"

	ci, err := parseInsertCode(dsl_string, "")
	ci.highlights.Init()

	parseHighlights(dsl_string, &ci.highlights, nil /*should not be needed*/)

	renderedCode := ci.renderCodeBlock()

	if renderedCode != `template <typename T>


}

int main(/**/) {

}
` || err != nil {
		t.Log("renderedCode: ", renderedCode)
		t.Error("Code was wrongly generated for `insert_code` with dot replacements.")
	}
}

func TestRenderWithPosIndent(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/foo.cpp"
	filet.File(t, codeFilePath, `template <typename T>
T shaveTheYak(T t) {
  return t;
}
`)

	dsl_string := "insert_code(" + codeFilePath + ":1-4)[indent=2]"

	ci, err := parseInsertCode(dsl_string, "")
	ci.highlights.Init()

	optionsStr := ""
	optionsStr, dsl_string = consumeOptionsString(dsl_string)
	ci.options = ParseCodeGenOptions(optionsStr)

	parseHighlights(dsl_string, &ci.highlights, nil /*should not be needed*/)

	renderedCode := ci.renderCodeBlock()

	expectedCode := `  template <typename T>
  T shaveTheYak(T t) {
    return t;
  }
`

	if renderedCode != expectedCode || err != nil {
		t.Logf("renderedCode:\n%sbut expected\n%s", renderedCode, expectedCode)
		t.Error("Code was wrongly generated for `insert_code` with dot replacements.")
	}
}

func TestRenderWithNegIndent(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/foo.cpp"
	filet.File(t, codeFilePath, `template <typename T>
T shaveTheYak(T t) {
  return t;
}
`)

	dsl_string := "insert_code(" + codeFilePath + ":1-4)[indent=-2]"

	ci, err := parseInsertCode(dsl_string, "")
	ci.highlights.Init()

	optionsStr := ""
	optionsStr, dsl_string = consumeOptionsString(dsl_string)
	ci.options = ParseCodeGenOptions(optionsStr)

	parseHighlights(dsl_string, &ci.highlights, nil /*should not be needed*/)

	renderedCode := ci.renderCodeBlock()

	expectedCode := `template <typename T>
T shaveTheYak(T t) {
return t;
}
`

	if renderedCode != expectedCode || err != nil {
		t.Logf("renderedCode:\n%sbut expected\n%s", renderedCode, expectedCode)
		t.Error("Code was wrongly generated for `insert_code` with dot replacements.")
	}
}

func TestRenderWithCommentThatShouldBeRemoved(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/foo.cpp"
	filet.File(t, codeFilePath, `template <typename T>
T shaveTheYak(T t) {
  // this is a comment
  return t;
}
`)

	dsl_string := "insert_code(" + codeFilePath + ":1-5)[comments=false]"

	ci, err := parseInsertCode(dsl_string, "")
	ci.highlights.Init()

	optionsStr := ""
	optionsStr, dsl_string = consumeOptionsString(dsl_string)
	ci.options = ParseCodeGenOptions(optionsStr)

	parseHighlights(dsl_string, &ci.highlights, nil /*should not be needed*/)

	renderedCode := ci.renderCodeBlock()

	expectedCode := `template <typename T>
T shaveTheYak(T t) {
  return t;
}
`

	if renderedCode != expectedCode || err != nil {
		t.Logf("renderedCode:\n%sbut expected\n%s", renderedCode, expectedCode)
		t.Error("Code was wrongly generated for `insert_code` with dot replacements.")
	}
}

//===----------------------------------------------------------------------===//
// rev_insert_code

func TestParseRevInsertCode(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/foo.cpp"
	filet.File(t, codeFilePath, `// code_block(FooID:1-4)
template <typename T>
T shaveTheYak(T t) {
  return t;
}

int main(/* int argc, char *argv[] */) {
  return shaveTheYak(42);
}
`)

	ci, err := parseRevInsertCode("rev_insert_code("+codeFilePath+":FooID)", "")

	renderedCode := ci.renderCodeBlock()

	if renderedCode != `template <typename T>
T shaveTheYak(T t) {
  return t;
}
` || err != nil {
		t.Log("renderedCode: ", renderedCode)
		t.Error("Code was wrongly generated for `insert_code`.")
	}
}

func TestParseRevInsertCodeTwoIDs(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/foo.cpp"
	filet.File(t, codeFilePath, `// code_block(BarID:1-2)
void barFunc {
}

// code_block(FooID:1-4)
template <typename T>
T shaveTheYak(T t) {
  return t;
}

int main(/* int argc, char *argv[] */) {
  return shaveTheYak(42);
}
`)

	ci, err := parseRevInsertCode("rev_insert_code("+codeFilePath+":FooID)", "")

	renderedCode := ci.renderCodeBlock()

	if renderedCode != `template <typename T>
T shaveTheYak(T t) {
  return t;
}
` || err != nil {
		t.Log("renderedCode: ", renderedCode)
		t.Error("Code was wrongly generated for `insert_code`.")
	}
}

func TestMissingParseRevInsertCode(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/foo.cpp"
	filet.File(t, codeFilePath, `// no code block here
template <typename T>
T shaveTheYak(T t) {
  return t;
}

int main(/* int argc, char *argv[] */) {
  return shaveTheYak(42);
}
`)

	_, err := parseRevInsertCode("rev_insert_code("+codeFilePath+":FooID)", "")

	if err == nil {
		t.Error("Code was wrongly generated for `insert_code`.")
	}
}

func TestRenderDotReplacementsLinesRevInsert(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/bazz.cpp"
	filet.File(t, codeFilePath, `// code_block(BazzID:1-8)
template <typename T>
T shaveTheYak(T t) {
  return t;
}

int main(int argc, char *argv[]) {
  return shaveTheYak(42);
}
`)
	dsl_string := "rev_insert_code(" + codeFilePath + ":BazzID)r<d2-3,d7,d6:{9-31}>"

	ci, err := parseRevInsertCode(dsl_string, "")

	renderedCode := ci.renderCodeBlock()

	if renderedCode != `template <typename T>
  // ...
}

int main(/* ... */) {
  // ...
}
` || err != nil {
		t.Log("renderedCode: ", renderedCode)
		t.Error("Code was wrongly generated for `rev_insert_code` with dot replacements.")
	}
}

func TestRenderRemoveLinesRevInsert(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/bazz.cpp"
	filet.File(t, codeFilePath, `// code_block(BazzID:1-8)
template <typename T>
T shaveTheYak(T t) {
  return t;
}

int main(int argc, char *argv[]) {
  return shaveTheYak(42);
}
`)
	dsl_string := "rev_insert_code(" + codeFilePath + ":BazzID)r<r2-3,r7,d6:{9-31}>"

	ci, err := parseRevInsertCode(dsl_string, "")

	renderedCode := ci.renderCodeBlock()

	if renderedCode != `template <typename T>
}

int main(/* ... */) {
}
` || err != nil {
		t.Log("renderedCode: ", renderedCode)
		t.Error("Code was wrongly generated for `rev_insert_code` with dot replacements.")
	}
}

func TestRenderDotReplacementsLinesAndHighlightRevInsert(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/bazz.cpp"
	filet.File(t, codeFilePath, `// code_block(BazzID:1-9)
template <typename T>
T shaveTheYak(T t) {
	t + 1;
  return t;
}

int main(int argc, char *argv[]) {
  return shaveTheYak(42);
}
`)
	dsl_string := "rev_insert_code(" + codeFilePath + ":BazzID)r<r3-4,r8,d7:{9-31}>r{1-2,7}"

	ci, err := parseRevInsertCode(dsl_string, "")

	renderedCode := ci.renderCodeBlock()

	if renderedCode != `*template <typename T>
*T shaveTheYak(T t) {
}

*int main(/* ... */) {
}
` || err != nil {
		t.Log("renderedCode: ", renderedCode)
		t.Error("Code was wrongly generated for `rev_insert_code` with dot replacements.")
	}
}

func TestRevInsertRenderWithCommentThatShouldBeRemoved(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/bazz.cpp"
	filet.File(t, codeFilePath, `// code_block(BazzID:1-9)
template <typename T>
T shaveTheYak(T t) {
  // comment
  return t;
}
`)
	dsl_string := "rev_insert_code(" + codeFilePath + ":BazzID)[comments=false]"

	ci, err := parseRevInsertCode(dsl_string, "")
	ci.highlights.Init()

	optionsStr := ""
	optionsStr, dsl_string = consumeOptionsString(dsl_string)
	ci.options = ParseCodeGenOptions(optionsStr)

	parseHighlights(dsl_string, &ci.highlights, nil /*should not be needed*/)

	renderedCode := ci.renderCodeBlock()

	expectedCode := `template <typename T>
T shaveTheYak(T t) {
  return t;
}
`

	if renderedCode != expectedCode || err != nil {
		t.Logf("renderedCode:\n%sbut expected\n%s", renderedCode, expectedCode)
		t.Error("Code was wrongly generated for `insert_code` with dot replacements.")
	}
}

//===----------------------------------------------------------------------===//
// Programming language tests

func TestGetProgrammingLanguage(t *testing.T) {
	if getProgrammingLanguage("foo.cpp") != "cpp" {
		t.Error("Wrong language detected from filename")
	}

	if getProgrammingLanguage("foo.py") != "python" {
		t.Error("Wrong language detected from filename")
	}
}
