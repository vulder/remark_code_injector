package main

import (
	"flag"
	"github.com/vulder/remark_code_injector/internal/html_processor"
	"strings"
)

func main() {
	inputFilepathPtr := flag.String("in", "index_raw.html", "Input file")
	outputFilepathPtr := flag.String("out", "nil", "Output file")
	codeRoot := flag.String("code-root", "", "Root folder where code files are stored.")

	flag.Parse()
	outputFilepath := *outputFilepathPtr
	if outputFilepath == "nil" {
		// If the user did not provided an output filename, try to infer the name
		// from the input file.
		outputFilepath = getDefaultOutputFile(*inputFilepathPtr)
	}

	html_processor.ProcessHTMLDocument(*inputFilepathPtr, outputFilepath, *codeRoot)
}

func getDefaultOutputFile(inputFile string) string {
	if strings.Contains(inputFile, "_raw") {
		return strings.ReplaceAll(inputFile, "_raw", "")
	}
	return "index.html"
}
