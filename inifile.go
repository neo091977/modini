package main

import (
	"fmt"
	"os"
	"regexp"
)

type iniFile struct {
	sectionOrder []string
	sectionMap   map[string]*iniSection
}

func newFile() iniFile {
	return iniFile{
		sectionOrder: []string{},
		sectionMap:   make(map[string]*iniSection),
	}
}

func (file *iniFile) getSection(name string) *iniSection {
	if section, isFound := file.sectionMap[name]; isFound {
		return section
	}

	section := newSection(name)
	file.sectionOrder = append(file.sectionOrder, name)
	file.sectionMap[name] = section
	return section
}

func (file *iniFile) allSections() []*iniSection {
	sections := []*iniSection{}
	for _, sectionName := range file.sectionOrder {
		section := file.sectionMap[sectionName]
		sections = append(sections, section)
	}
	return sections
}

var emptyRegex = regexp.MustCompile(`^\s*$`)
var sectionRegex = regexp.MustCompile(`^\[([^\[\]]+)\]\s*$`)
var propertyRegex = regexp.MustCompile(`^(\w[^=]*?)(=|\+=|\-=)(.*?)\s*$`)

func (file *iniFile) merge(lines []string, isFile bool) {
	var currentSection *iniSection

	for i, line := range lines {
		if emptyRegex.MatchString(line) {
			continue
		}

		sectionMatch := sectionRegex.FindStringSubmatch(line)
		if len(sectionMatch) > 0 {
			section := sectionMatch[1]
			currentSection = file.getSection(section)
			continue
		}

		propertyMatch := propertyRegex.FindStringSubmatch(line)
		if len(propertyMatch) > 0 {
			property := propertyMatch[1]
			operator := propertyMatch[2]
			value := propertyMatch[3]
			if isFile {
				if operator == "=" {
					currentSection.getProperty(property).addValue(value)
					continue
				}
			} else {
				switch operator {
				case "=":
					currentSection.getProperty(property).setValue(value)
					continue
				case "+=":
					currentSection.getProperty(property).addValue(value)
					continue
				case "-=":
					currentSection.getProperty(property).removeValue(value)
					continue
				}
			}
		}

		fmt.Fprintf(os.Stderr, "skipping line #%v: %v\n", i+1, line)
	}
}

func (file *iniFile) render() []string {
	result := []string{}
	for i, section := range file.allSections() {
		if i > 0 {
			result = append(result, "")
		}

		result = append(result, "["+section.name+"]")

		for _, property := range section.allProperties() {
			for _, value := range property.values {
				result = append(result, property.name+"="+value)
			}
		}
	}
	return result
}
