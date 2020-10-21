package main

import (
	"flag"
	"fmt"
	"github.com/vulder/remark_code_injector/internal/html_processor"
	"log"
)

func main() {
	filepathPtr := flag.String("file", "", "Path to HTML document")
	codeRoot := flag.String("code-root", "", "Root folder where code files are stored.")

	flag.Parse()

	if *filepathPtr == "" {
		log.Fatal("User did not provide a filepath to check.")
	}

	fmt.Print(html_processor.FindDependencies(*filepathPtr, *codeRoot))
}
