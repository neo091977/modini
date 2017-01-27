package main

import (
	"reflect"
	"testing"
)

func TestRenderOrder(t *testing.T) {
	input := sliceToChannel([]string{
		"[My.Section]",
		"Prop=value",
		"Prop=this is a test",
		"",
		"[Other.Section]",
		"Thing2=b",
		"Thing1=a",
		"Strange[0]=one",
		"",
		"[My.Section]",
		"Other=yes",
		"a=b",
	})

	expected := []string{
		"[My.Section]",
		"Prop=value",
		"Prop=this is a test",
		"Other=yes",
		"a=b",
		"",
		"[Other.Section]",
		"Thing2=b",
		"Thing1=a",
		"Strange[0]=one",
	}

	file := newFile()
	file.merge(input, true)
	actual := file.render()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(`Output not as expected. Actual: %v`, actual)
	}
}

func TestSkipsGarbage(t *testing.T) {
	input := sliceToChannel([]string{
		"garbage1",
		"[My.Section]",
		"Prop=value",
		"garbage2",
		"Prop=this is a test",
		"garbage3",
		"[Other.Section]",
		"Thing2=b",
		"Thing1=a",
		"garbage4",
		"=a",
		"[begin",
		"end]",
	})

	expected := []string{
		"[My.Section]",
		"Prop=value",
		"Prop=this is a test",
		"",
		"[Other.Section]",
		"Thing2=b",
		"Thing1=a",
	}

	file := newFile()
	file.merge(input, true)
	actual := file.render()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(`Output not as expected. Actual: %v`, actual)
	}
}

func TestMergesCorrectly(t *testing.T) {
	input1 := sliceToChannel([]string{
		"[My.Section]",
		"Prop=value",
		"Prop=this is a test",
		"",
		"[Other.Section]",
		"Thing2=b",
		"Thing1=a",
		"",
		"[My.Section]",
		"Other=yes",
		"Other=no",
	})

	input2 := sliceToChannel([]string{
		"[New.Section]",
		"Window=yes",
		"[My.Section]",
		"Prop+=another",
		"Prop-=this is a test",
		"Other=maybe",
		"[Other.Section]",
		"Thing2=c",
		"Thing3=d",
		"",
	})

	expected := []string{
		"[My.Section]",
		"Prop=value",
		"Prop=another",
		"Other=maybe",
		"",
		"[Other.Section]",
		"Thing2=c",
		"Thing1=a",
		"Thing3=d",
		"",
		"[New.Section]",
		"Window=yes",
	}

	file := newFile()
	file.merge(input1, true)
	file.merge(input2, false)
	actual := file.render()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(`Output not as expected. Actual: %v`, actual)
	}
}

func TestNoDuplicatesFromInput(t *testing.T) {
	input := sliceToChannel([]string{
		"[My.Section]",
		"Prop=value",
		"Prop=value",
	})

	expected := []string{
		"[My.Section]",
		"Prop=value",
	}

	file := newFile()
	file.merge(input, true)
	actual := file.render()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(`Output not as expected. Actual: %v`, actual)
	}
}

func TestNoDuplicatesFromModify(t *testing.T) {
	input1 := sliceToChannel([]string{
		"[My.Section]",
		"Prop=value",
		"Prop=b",
	})

	input2 := sliceToChannel([]string{
		"[My.Section]",
		"Prop+=value",
	})

	expected := []string{
		"[My.Section]",
		"Prop=value",
		"Prop=b",
	}

	file := newFile()
	file.merge(input1, true)
	file.merge(input2, false)
	actual := file.render()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(`Output not as expected. Actual: %v`, actual)
	}
}

func sliceToChannel(list []string) <-chan string {
	channel := make(chan string)
	go sliceToChannelHelper(channel, list)
	return channel
}

func sliceToChannelHelper(channel chan<- string, lines []string) {
	for _, line := range lines {
		channel <- line
	}
	close(channel)
}
