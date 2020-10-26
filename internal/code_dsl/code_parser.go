package code_dsl

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
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

// Checks if a line number is contained in this line number, i.e., is the same.
func (ln LineNumber) Contains(line_num int) bool {
	return ln.value == line_num
}

type LineRange struct {
	start int
	end   int
}

// Checks if a line number is contained in this line range.
func (lr LineRange) Contains(line_num int) bool {
	return line_num >= lr.start && line_num <= lr.end
}

type CodeBlock struct {
	fileRange LineRange
	lines     list.List
}

// Render a CodeBlock as a string
func (cb CodeBlock) render(highlights *Highlights) string {
	strRepr := ""

	line_num := cb.fileRange.start
	for e := cb.lines.Front(); e != nil; e = e.Next() {
		if highlights != nil && highlights.Contains(line_num) {
			strRepr += "*"
		}
		strRepr += e.Value.(string) + "\n"

		line_num++
	}

	return strRepr
}

type Highlights struct {
	highlightBlocks list.List
}

func (hl *Highlights) Init() {
	hl.highlightBlocks.Init()
}

// Adds a new block to the highlights list.
func (hl *Highlights) PushBack(v interface{}) *list.Element {
	return hl.highlightBlocks.PushBack(v)
}

// Checks if the line number is contained in one of the highlight blocks.
func (hl *Highlights) Contains(line_num int) bool {
	for e := hl.highlightBlocks.Front(); e != nil; e = e.Next() {
		switch v := e.Value.(type) {
		case LineNumber:
			if v.Contains(line_num) {
				return true
			}
		case LineRange:
			if v.Contains(line_num) {
				return true
			}
		default:
			panic("Highlight block type not supported.")
		}
	}

	return false
}

// Prints all highlight blocks
func (hl *Highlights) Show() {
	fmt.Printf("Highlights: %v\n", hl.highlightBlocks)
}

type CodeInsertion struct {
	codeBlock  CodeBlock
	progLang   string
	visuals    list.List
	highlights Highlights
}

func (ci CodeInsertion) renderCodeBlock() string {
	return ci.codeBlock.render(&ci.highlights)
}

type insertCodeInfo struct {
	filename  string
	filerange LineRange
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

var codeBlockRgx = regexp.MustCompile(".*code_block\\((?P<BlockID>.*):(?P<filerange>.*)\\).*")

func parseCodeBlockLineRangeFromFile(filepath string, blockID string) (LineRange, error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Could not open Source File: ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		match := codeBlockRgx.FindStringSubmatch(line)
		if match != nil {
			matchResults := make(map[string]string)
			for i, name := range codeBlockRgx.SubexpNames() {
				if i != 0 && name != "" {
					matchResults[name] = match[i]
				}
			}
			if matchResults["BlockID"] == blockID {
				filerange := strings.Split(matchResults["filerange"], "-")
				start, err := strconv.ParseInt(filerange[0], 10, 32)
				if err != nil {
					log.Fatal("Could not parse start of the file Range", err)
				}
				end, err := strconv.ParseInt(filerange[1], 10, 32)
				if err != nil {
					log.Fatal("Could not parse end of the file Range", err)
				}
				return LineRange{int(start), int(end)}, nil
			}
		}
	}

	return LineRange{0, 0}, errors.New("No valid BlockID found in file.")
}

var revInsertCodeRgx = regexp.MustCompile("rev_insert_code\\((?P<filename>.*):(?P<BlockID>.*)\\).*")

func parseRevInsertCodeInfo(line string, codeRoot string) (insertCodeInfo, error) {
	match := revInsertCodeRgx.FindStringSubmatch(line)
	if match == nil {
		panic("Line did not contain correct rev_insert_code pattern.")
	}
	matchResults := make(map[string]string)
	for i, name := range revInsertCodeRgx.SubexpNames() {
		if i != 0 && name != "" {
			matchResults[name] = match[i]
		}
	}

	filename := matchResults["filename"]
	lineRange, err := parseCodeBlockLineRangeFromFile(codeRoot+filename, matchResults["BlockID"])
	return insertCodeInfo{filename, lineRange}, err
}

var highlightCodeRgx = regexp.MustCompile("insert_code\\(.*\\).*?(?P<rel>[r\\{]+)(?P<highlights>.*)\\}")

func parseHighlights(line string, highlights *Highlights, baseCodeRange *LineRange) {
	match := highlightCodeRgx.FindStringSubmatch(line)
	if match == nil { // Return when we did not find any highlights
		return
	}
	matchResults := make(map[string]string)
	for i, name := range highlightCodeRgx.SubexpNames() {
		if i != 0 && name != "" {
			matchResults[name] = match[i]
		}
	}

	handleLinesRelative := matchResults["rel"] == "r{"

	blocks := strings.Split(matchResults["highlights"], ",")
	for _, block := range blocks {
		if strings.Contains(block, "-") {
			block_split := strings.Split(block, "-")
			start, err := strconv.ParseInt(block_split[0], 10, 32)
			if err != nil {
				log.Fatal("Could not parse start of the highlight Range", err)
			}
			end, err := strconv.ParseInt(block_split[1], 10, 32)
			if err != nil {
				log.Fatal("Could not parse end of the highlight Range", err)
			}
			if handleLinesRelative {
				// -1 is relevant because line numbers start a 1 not 0
				start = int64(baseCodeRange.start) + start - 1
				end = int64(baseCodeRange.start) + end - 1
			}
			highlights.PushBack(LineRange{int(start), int(end)})
		} else { // Handle single line number
			line_num, err := strconv.ParseInt(block, 10, 32)
			if err != nil {
				log.Fatal("Could not parse highlight line number", err)
			}
			if handleLinesRelative {
				// -1 is relevant because line numbers start a 1 not 0
				line_num = int64(baseCodeRange.start) + line_num - 1
			}
			highlights.PushBack(LineNumber{int(line_num)})
		}
	}
}

func parseInsertCode(line string, codeRoot string) CodeInsertion {
	ic_info := parserInsertCodeInfo(line)

	ci := CodeInsertion{}
	ci.codeBlock = parseCodeBlock(codeRoot+ic_info.filename, ic_info.filerange.start, ic_info.filerange.end)
	ci.progLang = getProgrammingLanguage(ic_info.filename)
	ci.visuals.Init()
	ci.highlights.Init()
	parseHighlights(line, &ci.highlights, &ic_info.filerange)
	return ci
}

func parseRevInsertCode(line string, codeRoot string) (CodeInsertion, error) {
	ic_info, err := parseRevInsertCodeInfo(line, codeRoot)
	if err != nil {
		return CodeInsertion{}, err
	}

	ci := CodeInsertion{}
	ci.codeBlock = parseCodeBlock(codeRoot+ic_info.filename, ic_info.filerange.start, ic_info.filerange.end)
	ci.progLang = getProgrammingLanguage(ic_info.filename)
	ci.visuals.Init()
	ci.highlights.Init()
	parseHighlights(line, &ci.highlights, &ic_info.filerange)
	return ci, err
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
