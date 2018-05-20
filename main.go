package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "USAGE: circle-art <jpg-file>")
		os.Exit(1)
	}

	input := os.Args[1]

	sg := NewSVGGrid()
	//sg.RenderGrid(&CircularGradient{})
	ic, err := NewImageContent(input)
	if err != nil {
		panic(err)
	}
	outputPrefix := filepath.Base(input)
	outputPrefix = strings.TrimSuffix(outputPrefix, filepath.Ext(outputPrefix))
	sg.RenderGrid(ic, outputPrefix)
}
