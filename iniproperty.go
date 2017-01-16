package main

type iniProperty struct {
	name   string
	values []string
}

func newProperty(name string) *iniProperty {
	return &iniProperty{
		name:   name,
		values: []string{},
	}
}

func (property *iniProperty) setValue(value string) {
	property.values = []string{value}
}

func (property *iniProperty) addValue(value string) {
	property.values = append(property.values, value)
}

func (property *iniProperty) removeValue(value string) {
	for i := len(property.values) - 1; i >= 0; i-- {
		if property.values[i] == value {
			property.values = append(property.values[:i], property.values[i+1:]...)
		}
	}
}
