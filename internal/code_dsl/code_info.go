package code_dsl

// Returns the relative filepath of dependent file, i.e., the path to a file
// that is needed in a DSL command.
func GetFileDependency(line string, codeRoot string) string {
	if isInsertCode(line) {
		return parserInsertCodeInfo(line).filename
	}
	if isRevInsertCode(line) {
		ci, err := parseRevInsertCodeInfo(line, codeRoot)
		if err == nil {
			return ci.filename
		}
	}
	return ""
}
