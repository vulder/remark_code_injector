package code_dsl

import (
	"github.com/Flaque/filet"
	"testing"
)

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
*  return t;
}
` {
		t.Log("renderedCode: ", renderedCode)
		t.Error("Code was wrongly generated for `insert_code` with highlights.")
	}
}

//===----------------------------------------------------------------------===//
// rev_insert_code

func TestParseRevInsertCode(t *testing.T) {
	defer filet.CleanUp(t)
	codeFilePath := filet.TmpDir(t, "") + "/foo.cpp"
	filet.File(t, codeFilePath, `// code_block(FooID:2-5)
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
	filet.File(t, codeFilePath, `// code_block(BarID:2-3)
void barFunc {
}

// code_block(FooID:6-9)
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
