package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sairash/gotots"
)

func main() {
	dir := flag.String("dir", "", "input directory containing Go files")
	output := flag.String("output", "", "output TypeScript file path")
	flag.Parse()

	if *dir == "" || *output == "" {
		fmt.Fprintf(os.Stderr, "Usage: gotots -dir <input_dir> -output <output_file>\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := gotots.New().FromDir(*dir).ToFile(*output).Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s from %s\n", *output, *dir)
}
