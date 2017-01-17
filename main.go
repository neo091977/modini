package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

var version string

func main() {
	inFile := flag.String("input", "", "(required) Path to input ini file.")
	outFile := flag.String("output", "", "(optional) Path to output ini file.")
	delimiter := flag.String("delimit", ";", "(optional) Split the modify string on this value.")
	modify := flag.String("modify", "", "(optional) Modifications to make. ex: [section1];prop1=value1;prop2=value2;[section2];prop1=value3.")
	showVersion := flag.Bool("version", false, "output the version")

	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	if *inFile == "" {
		log.Fatal("argument --input is required")
	}

	inputLines := readFile(*inFile)
	modifyLines := strings.Split(*modify, *delimiter)

	file := newFile()
	file.merge(inputLines, true)
	file.merge(modifyLines, false)

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
		log.Fatal(err)
	}
	return strings.Split(string(content), "\n")
}

func writeFile(filename string, contents string) {
	err := ioutil.WriteFile(filename, []byte(contents), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
