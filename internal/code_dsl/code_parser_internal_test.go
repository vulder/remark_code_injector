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

	renderedCode := cb.render(nil)

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

	ci := parseInsertCode("insert_code("+codeFilePath+":1-4)", "")

	renderedCode := ci.renderCodeBlock()
	t.Log(ci.renderCodeBlock())

	if renderedCode != `template <typename T>
T shaveTheYak(T t) {
  return t;
}
` {
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

	ci := parseInsertCode(dsl_string, "")
	ci.highlights.Init()

	parseHighlights(dsl_string, &ci.highlights, nil /*should not be needed*/)

	renderedCode := ci.renderCodeBlock()

	if renderedCode != `*template <typename T>
*T shaveTheYak(T t) {
  return t;
*}
` {
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
	ci := parseInsertCode(dsl_string, "")
	ci.highlights.Init()

	baseCodeRange := LineRange{1, 4}
	parseHighlights(dsl_string, &ci.highlights, &baseCodeRange)

	renderedCode := ci.renderCodeBlock()

	if renderedCode != `template <typename T>
*T shaveTheYak(T t) {
* return t;
}
` {
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
	ci := parseInsertCode(dsl_string, "")
	ci.highlights.Init()

	baseCodeRange := LineRange{1, 4}
	parseHighlights(dsl_string, &ci.highlights, &baseCodeRange)

	renderedCode := ci.renderCodeBlock()

	expectedCode := "*template <typename T>\n"
	expectedCode += "T `shaveTheYak`(T `t`) {\n"
	expectedCode += "* return t;\n"
	expectedCode += "}\n"

	if renderedCode != expectedCode {
		t.Log("renderedCode: ", renderedCode, " expected: ", expectedCode)
		t.Error("Highlights were wrongly generated for `insert_code`.")
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
