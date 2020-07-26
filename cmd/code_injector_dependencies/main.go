package main

import (
	"flag"
	"fmt"
	"github.com/vulder/remark_code_injector/internal/html_processor"
	"log"
)

func main() {
	filepathPtr := flag.String("file", "", "Path to HTML document")

	flag.Parse()

	if *filepathPtr == "" {
		log.Fatal("User did not provide a filepath to check.")
	}

	fmt.Print(html_processor.FindDependencies(*filepathPtr))
}
