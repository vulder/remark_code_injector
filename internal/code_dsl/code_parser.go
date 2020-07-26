package code_dsl

import (
	"bufio"
	"container/list"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type LineNumber struct {
	value int
}

type LineRange struct {
	start int
	end   int
}

type CodeBlock struct {
	fileRange LineRange
	lines     list.List
}

// Render a CodeBlock as a string
func (cb CodeBlock) render() string {
	strRepr := ""
	for e := cb.lines.Front(); e != nil; e = e.Next() {
		strRepr += e.Value.(string) + "\n"
	}

	return strRepr
}

type CodeInsertion struct {
	codeBlock  CodeBlock
	progLang   string
	visuals    list.List
	highlights list.List
}

type insertCodeInfo struct {
	filename  string
	filerange LineRange
}

func (ci CodeInsertion) renderCodeBlock() string {
	return ci.codeBlock.render()
}

var insertCodeRgx = regexp.MustCompile("insert_code\\((?P<filename>.*):(?P<filerange>.*).*\\).*")

func parserInsertCodeInfo(line string) insertCodeInfo {
	match := insertCodeRgx.FindStringSubmatch(line)
	matchResults := make(map[string]string)
	for i, name := range insertCodeRgx.SubexpNames() {
		if i != 0 && name != "" {
			matchResults[name] = match[i]
		}
	}

	filename := matchResults["filename"]
	filerange := strings.Split(matchResults["filerange"], "-")
	start, err := strconv.ParseInt(filerange[0], 10, 32)
	if err != nil {
		log.Fatal("Could not parse start of the file Range", err)
	}
	end, err := strconv.ParseInt(filerange[1], 10, 32)
	if err != nil {
		log.Fatal("Could not parse end of the file Range", err)
	}

	return insertCodeInfo{filename, LineRange{int(start), int(end)}}
}

func parseInsertCode(line string, codeRoot string) CodeInsertion {
	ic_info := parserInsertCodeInfo(line)

	ci := CodeInsertion{}
	ci.codeBlock = parseCodeBlock(codeRoot+ic_info.filename, ic_info.filerange.start, ic_info.filerange.end)
	ci.progLang = getProgrammingLanguage(ic_info.filename)
	ci.visuals.Init()
	ci.highlights.Init()
	return ci
}

func parseCodeBlock(filepath string, start int, end int) CodeBlock {
	cb := CodeBlock{}
	cb.fileRange = LineRange{start, end}
	cb.lines.Init()

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Could not open Source File: ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		if lineNumber >= start && lineNumber <= end {
			cb.lines.PushBack(scanner.Text())
		}
		lineNumber++
	}

	return cb
}

func getProgrammingLanguage(filename string) string {
	filetype := getFiletype(filename)
	filetype = strings.ReplaceAll(filetype, ".", "")
	switch {
	case filetype == "py":
		return "python"
	default:
		return filetype
	}
}

func getFiletype(path string) string {
	return filepath.Ext(path)
}
