package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	inFile    = kingpin.Flag("input", "Path to input ini file.").Required().Short('i').String()
	outFile   = kingpin.Flag("output", "Path to output ini file.").Default("").Short('o').String()
	delimiter = kingpin.Flag("delimit", "Split the mod string on this value.").Default(";").Short('d').String()
	modify    = kingpin.Flag("modify", "Modifications to make. ex: [section1];prop1=value1;prop2=value2;[section2];prop1=value3.").Default("").Short('m').String()
)

var version string

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version(version).Author("Greg Todd")
	kingpin.Parse()

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
