package code_dsl

import (
	"fmt"
	"strconv"
	"strings"
)

type CodeGenOptionsImpl struct {
	indentLevel    int
	removeComments bool
}

type CodeGenOptions interface {
	hideComments() bool
	getIndent() int
}

func (cgo *CodeGenOptionsImpl) hideComments() bool {
	return cgo.removeComments
}

func (cgo *CodeGenOptionsImpl) getIndent() int {
	return cgo.indentLevel
}

func ParseCodeGenOptions(optionString string) CodeGenOptions {
	cgo := CodeGenOptionsImpl{}
	cgo.indentLevel = 0
	cgo.removeComments = false

	for _, optionStr := range strings.Split(optionString, ",") {
		optionStr = strings.TrimSpace(optionStr)
		if optionStr == "" {
			continue
		}

		kvPair := strings.Split(optionStr, "=")
		optionKey := strings.ToLower(kvPair[0])
		optionValue := kvPair[1]

		switch optionKey {
		case "comments":
			enableComment, err := strconv.ParseBool(optionValue)
			if err != nil {
				fmt.Println("Could not parse option:", err.Error())
			}
			cgo.removeComments = !enableComment
		case "indent":
			indentationLevel, err := strconv.ParseInt(optionValue, 10, 32)
			if err != nil {
				fmt.Println("Could not parse option:", err.Error())
			}
			cgo.indentLevel = int(indentationLevel)
		default:
			fmt.Println("Did not understand option key:", optionKey)
		}
	}

	return &cgo
}

func MakeDefaultCodeGenOptions() CodeGenOptions {
	return &CodeGenOptionsImpl{}
}
