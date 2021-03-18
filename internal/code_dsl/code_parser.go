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
	"sort"
	"strconv"
	"strings"
)

type LineNumber struct {
	value int
}

// Checks if a line number is contained in this line number, i.e., is the same.
func (ln LineNumber) Contains(lineNum int) bool {
	return ln.value == lineNum
}

type LineRange struct {
	start int
	end   int
}

// Checks if a line number is contained in this line range.
func (lr LineRange) Contains(lineNum int) bool {
	return lineNum >= lr.start && lineNum <= lr.end
}

type CharRange struct {
	lineNum LineNumber
	start   int
	end     int
}

func (cr CharRange) Contains(lineNum int) bool {
	return lineNum >= cr.lineNum.value && lineNum <= cr.lineNum.value
}

type CodeBlock struct {
	fileRange LineRange
	lines     list.List
}

// Render a CodeBlock as a string
func (cb CodeBlock) render(highlights *Highlights, visuals *VisualModifications, language string) string {
	strRepr := ""

	lineNum := cb.fileRange.start
	for e := cb.lines.Front(); e != nil; e = e.Next() {
		line := e.Value.(string)

		if highlights != nil && highlights.Contains(lineNum) {
			if highlights.HasSubrange(lineNum) {
				line = highlights.RenderSubrange(line, lineNum)
			} else {
				strRepr += "*"
				line = strings.TrimPrefix(line, " ")
			}
		}

		if visuals != nil {
			vLine, useLine := visuals.ModifyLine(line, lineNum, language)
			if !useLine { // Skip line
				lineNum++
				continue
			}
			line = vLine
		}

		strRepr += line + "\n"

		lineNum++
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
func (hl *Highlights) Contains(lineNum int) bool {
	for e := hl.highlightBlocks.Front(); e != nil; e = e.Next() {
		switch v := e.Value.(type) {
		case CharRange:
			if v.Contains(lineNum) {
				return true
			}
		case LineNumber:
			if v.Contains(lineNum) {
				return true
			}
		case LineRange:
			if v.Contains(lineNum) {
				return true
			}
		default:
			panic("Highlight block type not supported.")
		}
	}

	return false
}

func (hl *Highlights) HasSubrange(lineNum int) bool {
	for e := hl.highlightBlocks.Front(); e != nil; e = e.Next() {
		switch v := e.Value.(type) {
		case CharRange:
			if v.Contains(lineNum) {
				return true
			}
		case LineNumber:
		case LineRange:
		default:
			panic("Highlight block type not supported.")
		}
	}

	return false
}

type VisualModificationType int

const (
	ReplaceWithDots VisualModificationType = iota
	Hide                                   = iota
	Remove                                 = iota
)

type VisualModification struct {
	lineRangeSpecifier interface{}
	modeType           VisualModificationType
}

type VisualModifications struct {
	modifications list.List
}

func (vm *VisualModifications) Init() {
	vm.modifications.Init()
}

func (vm *VisualModifications) PushBack(v VisualModification) *list.Element {
	return vm.modifications.PushBack(v)
}

func (vm *VisualModifications) ModifyLine(line string, lineNum int, language string) (string, bool) {
	for e := vm.modifications.Front(); e != nil; e = e.Next() {
		vm := e.Value.(VisualModification)
		switch lrs := vm.lineRangeSpecifier.(type) {
		case CharRange:
			if lrs.Contains(lineNum) {
				placeHolder := ""
				if vm.modeType == ReplaceWithDots {
					placeHolder = " ... "
				} else if vm.modeType == Hide {
					placeHolder = ""
				}
				lineRunes := []rune(line)
				if lrs.end > len(line) {
					panic("Specified visual range ends after line.")
				}
				return string(lineRunes[:lrs.start]) + makeMultilineComment(placeHolder, language) + string(lineRunes[lrs.end:]), true
			}
		case LineNumber:
			if lrs.Contains(lineNum) {
				if vm.modeType == ReplaceWithDots {
					return strings.Repeat(" ", getIndent(line)) + makeComment(" ...", language), true
				} else if vm.modeType == Hide {
					return "", true
				} else if vm.modeType == Remove {
					return "", false
				} else {
					panic("Unsupported visual modifier")
				}
			}
		case LineRange:
			if lrs.Contains(lineNum) {
				if vm.modeType == ReplaceWithDots {
					if lrs.end != lineNum {
						return "", false
					}
					return strings.Repeat(" ", getIndent(line)) + makeComment(" ...", language), true
				} else if vm.modeType == Hide {
					return "", true
				} else if vm.modeType == Remove {
					return "", false
				} else {
					panic("Unsupported visual modifier")
				}
			}
		default:
			panic("VisualModification lineRangeSpecifier type not supported.")
		}
	}

	return line, true
}

// TODO: refactor to own util file
func insertAt(baseStr string, pos int, text string) string {
	updatedString := ""
	if pos <= len(baseStr) {
		updatedString += baseStr[:pos]
	} else {
		updatedString += baseStr
	}

	updatedString += text

	if pos <= len(baseStr) {
		updatedString += baseStr[pos:]
	}
	return updatedString
}

type ReverseCharRange []CharRange

func (r ReverseCharRange) Len() int           { return len(r) }
func (r ReverseCharRange) Less(i, j int) bool { return r[i].start > r[j].start }
func (r ReverseCharRange) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

// TODO: maybe we should put highlights into it's own file
func (hl *Highlights) RenderSubrange(line string, lineNum int) string {
	charRanges := []CharRange{}
	for e := hl.highlightBlocks.Front(); e != nil; e = e.Next() {
		switch v := e.Value.(type) {
		case CharRange:
			if v.Contains(lineNum) {
				charRanges = append(charRanges, v)
			}
		case LineNumber:
		case LineRange:
		default:
			panic("Highlight block type not supported.")
		}
	}

	sort.Sort(ReverseCharRange(charRanges))

	for _, charRange := range charRanges {
		line = insertAt(line, charRange.end, "`")
		line = insertAt(line, charRange.start-1, "`")
	}

	return line
}

// Prints all highlight blocks
func (hl *Highlights) Show() {
	fmt.Printf("Highlights: %v\n", hl.highlightBlocks)
}

type CodeInsertion struct {
	codeBlock  CodeBlock
	progLang   string
	visuals    VisualModifications
	highlights Highlights
}

func (ci CodeInsertion) renderCodeBlock() string {
	return ci.codeBlock.render(&ci.highlights, &ci.visuals, ci.progLang)
}

type insertCodeInfo struct {
	filename  string
	filerange LineRange
}

var insertCodeRgx = regexp.MustCompile("insert_code\\((?P<filename>.*):(?P<filerange>.*).*\\).*")

func parserInsertCodeInfo(line string) (insertCodeInfo, error) {
	match := insertCodeRgx.FindStringSubmatch(line)
	if match == nil {
		panic("Line did not contain correct insert_code pattern.")
	}
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

	return insertCodeInfo{filename, LineRange{int(start), int(end)}}, nil
}

var codeBlockRgx = regexp.MustCompile(".*code_block\\((?P<BlockID>.*):(?P<filerange>.*)\\).*")

func parseCodeBlockLineRangeFromFile(filepath string, blockID string) (LineRange, error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Could not open Source File: ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
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
				return LineRange{lineNumber + int(start), lineNumber + int(end)}, nil
			}
		}

		lineNumber++
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

type appendCharRange func(CharRange)

func parseCharRanges(block string, addCharRange appendCharRange, baseCodeRange *LineRange, handleLinesRelative bool) {
	lineCharRange := strings.Split(block, ":")
	lineNum, err := strconv.ParseInt(lineCharRange[0], 10, 32)
	charRanges := lineCharRange[1]
	if err != nil {
		log.Fatal("Could not parse line number", err)
	}

	if handleLinesRelative {
		// -1 is relevant because line numbers start a 1 not 0
		lineNum = int64(baseCodeRange.start) + lineNum - 1
	}

	charRanges = strings.TrimLeft(charRanges, "{")
	charRanges = strings.TrimRight(charRanges, "}")

	splitCharRanges := strings.Split(charRanges, "|")
	for _, charRange := range splitCharRanges {
		splitCharRange := strings.Split(charRange, "-")

		charStart, err := strconv.ParseInt(splitCharRange[0], 10, 32)
		if err != nil {
			log.Fatal("Could not parse char range start", err)
		}

		charEnd, err := strconv.ParseInt(splitCharRange[1], 10, 32)
		if err != nil {
			log.Fatal("Could not parse char range end", err)
		}

		addCharRange(CharRange{LineNumber{int(lineNum)}, int(charStart), int(charEnd)})
	}
}

func parseCharRangesHighlights(hlBlock string, highlights *Highlights, baseCodeRange *LineRange, handleLinesRelative bool) {
	addHighlight := func(cr CharRange) {
		highlights.PushBack(cr)
	}
	parseCharRanges(hlBlock, addHighlight, baseCodeRange, handleLinesRelative)
}

func parseCharRangesVisuals(vlBlock string, visuals *VisualModifications, baseCodeRange *LineRange, handleLinesRelative bool, vmt VisualModificationType) {
	addVisual := func(cr CharRange) {
		visuals.PushBack(VisualModification{cr, vmt})
	}
	parseCharRanges(vlBlock, addVisual, baseCodeRange, handleLinesRelative)
}

type appendLineRange func(LineRange)

func parseLineRange(block string, addLineRange appendLineRange, baseCodeRange *LineRange, handleLinesRelative bool) {
	block_split := strings.Split(block, "-")
	start, err := strconv.ParseInt(block_split[0], 10, 32)
	if err != nil {
		log.Fatal("Could not parse start of the Range ", err)
	}
	end, err := strconv.ParseInt(block_split[1], 10, 32)
	if err != nil {
		log.Fatal("Could not parse end of the Range ", err)
	}
	if handleLinesRelative {
		// -1 is relevant because line numbers start a 1 not 0
		start = int64(baseCodeRange.start) + start - 1
		end = int64(baseCodeRange.start) + end - 1
	}
	addLineRange(LineRange{int(start), int(end)})
}

func parseLineRangeHightlights(block string, highlights *Highlights, baseCodeRange *LineRange, handleLinesRelative bool) {
	addHighlight := func(cr LineRange) {
		highlights.PushBack(cr)
	}
	parseLineRange(block, addHighlight, baseCodeRange, handleLinesRelative)
}

func parseLineRangeVisuals(block string, visuals *VisualModifications, baseCodeRange *LineRange, handleLinesRelative bool, vmt VisualModificationType) {
	addVisual := func(cr LineRange) {
		visuals.PushBack(VisualModification{cr, vmt})
	}
	parseLineRange(block, addVisual, baseCodeRange, handleLinesRelative)
}

type appendLineNumber func(LineNumber)

func parseLineNumber(block string, addLineNumber appendLineNumber, baseCodeRange *LineRange, handleLinesRelative bool) {
	lineNum, err := strconv.ParseInt(block, 10, 32)
	if err != nil {
		log.Fatal("Could not parse line number", err)
	}
	if handleLinesRelative {
		// -1 is relevant because line numbers start a 1 not 0
		lineNum = int64(baseCodeRange.start) + lineNum - 1
	}
	addLineNumber(LineNumber{int(lineNum)})
}

func parseLineNumberHightlights(block string, highlights *Highlights, baseCodeRange *LineRange, handleLinesRelative bool) {
	addHighlight := func(cr LineNumber) {
		highlights.PushBack(cr)
	}
	parseLineNumber(block, addHighlight, baseCodeRange, handleLinesRelative)
}

func parseLineNumberVisuals(block string, visuals *VisualModifications, baseCodeRange *LineRange, handleLinesRelative bool, vmt VisualModificationType) {
	addVisual := func(cr LineNumber) {
		visuals.PushBack(VisualModification{cr, vmt})
	}
	parseLineNumber(block, addVisual, baseCodeRange, handleLinesRelative)
}

// Works for rev_insert_code and insert_code
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

	if strings.HasPrefix(matchResults["highlights"], "<") ||
		strings.HasPrefix(matchResults["highlights"], "r<") { // Skip visual matches with sub range constructs
		return
	}

	handleLinesRelative := matchResults["rel"] == "r{"

	blocks := strings.Split(matchResults["highlights"], ",")
	for _, block := range blocks {
		if strings.Contains(block, ":") { // Got and inline hl block
			parseCharRangesHighlights(block, highlights, baseCodeRange, handleLinesRelative)
		} else if strings.Contains(block, "-") {
			parseLineRangeHightlights(block, highlights, baseCodeRange, handleLinesRelative)
		} else { // Handle single line number
			parseLineNumberHightlights(block, highlights, baseCodeRange, handleLinesRelative)
		}
	}
}

// Works for rev_insert_code and insert_code
var visualCodeRgx = regexp.MustCompile("insert_code\\(.*\\).*?(?P<mod>[r]*\\<)(?P<visuals>.*)\\>")

func parseVisuals(line string, visuals *VisualModifications, baseCodeRange *LineRange) {
	match := visualCodeRgx.FindStringSubmatch(line)
	if match == nil { // Return when we did not find any visuals
		return
	}
	matchResults := make(map[string]string)
	for i, name := range visualCodeRgx.SubexpNames() {
		if i != 0 && name != "" {
			matchResults[name] = match[i]
		}
	}

	handleLinesRelative := strings.Contains(matchResults["mod"], "r")

	blocks := strings.Split(matchResults["visuals"], ",")
	for _, block := range blocks {
		replaceWithDots := strings.HasPrefix(block, "d")
		hideLines := strings.HasPrefix(block, "h")
		removeLines := strings.HasPrefix(block, "r")

		if !replaceWithDots && !hideLines && !removeLines {
			log.Println("No visual modification type set, defaulting to hidding the lines.")
			hideLines = true
		}
		getVisualModType := func() VisualModificationType {
			if replaceWithDots {
				return ReplaceWithDots
			} else if removeLines {
				return Remove
			} else {
				return Hide
			}
		}

		block = strings.TrimLeft(block, "hdr")

		if strings.Contains(block, ":") { // Got and inline hl block
			parseCharRangesVisuals(block, visuals, baseCodeRange, handleLinesRelative, getVisualModType())
		} else if strings.Contains(block, "-") {
			parseLineRangeVisuals(block, visuals, baseCodeRange, handleLinesRelative, getVisualModType())
		} else { // Handle single line number
			parseLineNumberVisuals(block, visuals, baseCodeRange, handleLinesRelative, getVisualModType())
		}
	}
}

func parseInsertCode(line string, codeRoot string) (CodeInsertion, error) {
	icInfo, err := parserInsertCodeInfo(line)
	if err != nil {
		return CodeInsertion{}, err
	}

	ci := CodeInsertion{}
	ci.codeBlock = parseCodeBlock(codeRoot+icInfo.filename, icInfo.filerange.start, icInfo.filerange.end)
	ci.progLang = getProgrammingLanguage(icInfo.filename)
	ci.visuals.Init()
	ci.highlights.Init()
	parseHighlights(line, &ci.highlights, &icInfo.filerange)
	parseVisuals(line, &ci.visuals, &icInfo.filerange)
	return ci, nil
}

func parseRevInsertCode(line string, codeRoot string) (CodeInsertion, error) {
	icInfo, err := parseRevInsertCodeInfo(line, codeRoot)
	if err != nil {
		return CodeInsertion{}, err
	}

	ci := CodeInsertion{}
	ci.codeBlock = parseCodeBlock(codeRoot+icInfo.filename, icInfo.filerange.start, icInfo.filerange.end)
	ci.progLang = getProgrammingLanguage(icInfo.filename)
	ci.visuals.Init()
	ci.highlights.Init()
	parseHighlights(line, &ci.highlights, &icInfo.filerange)
	parseVisuals(line, &ci.visuals, &icInfo.filerange)
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

// Returns the number of spaces used to indent a line
func getIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

func makeComment(line string, language string) string {
	return "//" + line
}

func makeMultilineComment(line string, language string) string {
	return "/*" + line + "*/"
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
