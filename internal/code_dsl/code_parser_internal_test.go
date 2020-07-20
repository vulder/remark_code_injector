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

	renderedCode := cb.render()

	if renderedCode != `T shaveTheYak(T t) {
  return t;
}
` {
		t.Log("renderedCode: ", renderedCode)
		t.Error("Parsed Code was wrongly redered.")
	}
}

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

	ci := parseInsertCode("insert_code(" + codeFilePath + ":1-4)")

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

func TestGetProgrammingLanguage(t *testing.T) {
	if getProgrammingLanguage("foo.cpp") != "cpp" {
		t.Error("Wrong language detected from filename")
	}

	if getProgrammingLanguage("foo.py") != "python" {
		t.Error("Wrong language detected from filename")
	}
}
