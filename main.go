package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var version string

func main() {
	inFile := flag.String("input", "", "Path to input ini file.")
	outFile := flag.String("output", "", "Path to output ini file.")
	delimiter := flag.String("delimit", ";", "Split the modify string on this value.")
	modify := flag.String("modify", "", "Modifications to make. ex: [section1];prop1=value1;prop2=value2;[section2];prop1=value3.")
	showVersion := flag.Bool("version", false, "output the version")

	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	if *inFile == "" && *modify == "" {
		fmt.Fprintf(os.Stderr, "arguments -input and/or -modify are required\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	file := newFile()

	if *inFile != "" {
		inputLines := readFile(*inFile)
		file.merge(inputLines, true)
	}

	if *modify != "" {
		modifyLines := strings.Split(*modify, *delimiter)
		file.merge(modifyLines, false)
	}

	output := strings.Join(file.render(), "\n")

	if *outFile == "" {
		fmt.Println(output)
	} else {
		writeFile(*outFile, output+"\n")
	}
}

func readFile(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(2)
	}
	return strings.Split(string(content), "\n")
}

func writeFile(filename string, contents string) {
	err := ioutil.WriteFile(filename, []byte(contents), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing file: %v\n", err)
		os.Exit(3)
	}
}
