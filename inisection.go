package main

type iniSection struct {
	name          string
	propertyOrder []string
	propertyMap   map[string]*iniProperty
}

func newSection(name string) *iniSection {
	return &iniSection{
		name:          name,
		propertyOrder: []string{},
		propertyMap:   make(map[string]*iniProperty),
	}
}

func (section *iniSection) getProperty(name string) *iniProperty {
	if property, isFound := section.propertyMap[name]; isFound {
		return property
	}

	property := newProperty(name)
	section.propertyOrder = append(section.propertyOrder, name)
	section.propertyMap[name] = property
	return property
}

func (section *iniSection) allProperties() []*iniProperty {
	properties := []*iniProperty{}
	for _, propertyName := range section.propertyOrder {
		property := section.propertyMap[propertyName]
		properties = append(properties, property)
	}
	return properties
}
