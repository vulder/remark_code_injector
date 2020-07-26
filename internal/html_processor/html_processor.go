package html_processor

import (
	"bufio"
	"log"
	"os"

	"github.com/vulder/remark_code_injector/internal/code_dsl"
)

// Generates a new version of the HTML document, replacing all DSL annotations
// with the generated content.
func ProcessHTMLDocument(inputFilepath string, outputFilepath string, codeRoot string) {
	file, err := os.Open(inputFilepath)
	if err != nil {
		log.Fatal("Could not open HTML document: ", err)
	}
	defer file.Close()
	outputFile, err := os.Create(outputFilepath)

	if err != nil {
		log.Fatal("Could not create output file: ", err)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(file)
	sep := ""
	for scanner.Scan() {
		outputFile.WriteString(sep + handleHTMLLine(scanner.Text(), codeRoot))
		sep = "\n"
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Could not scan HTML docuemnt: ", err)
	}
}

// Computes a ';' separated string of all files used in DSL commands.
func FindDependencies(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Could not open HTML document: ", err)
	}
	defer file.Close()

	dependencies := ""
	scanner := bufio.NewScanner(file)
	sep := ""
	for scanner.Scan() {
		dep := code_dsl.GetFileDependency(scanner.Text())
		if dep != "" {
			dependencies += sep + dep
			sep = ";"
		}
	}

	return dependencies
}

func handleHTMLLine(line string, codeRoot string) string {
	if code_dsl.ContainsDSLCommand(line) {
		return code_dsl.TransformLine(line, codeRoot)
	}
	return line
}
