package tabulon

import (
	"testing"
	"reflect"
//	"fmt"
)

func Test_newReader(t *testing.T) {
	reader := NewCSVReader()
	if len(reader.GetHeader())!=0 {
		t.Error("header should be empty")
	}
	if reader.delimiter != ',' {
		t.Error("default delimiter is not comma")
	}

	row := reader.ParseLine("h1,h2,h3")
	if row != nil {
		t.Error("header not read correctly")
	}
	if len(reader.row)!=0 {
		t.Error("row has old content")
	}

	row = reader.ParseLine("content1,content2,content3")
	if len(reader.GetHeader())!=3 {
		t.Error("header not read correctly")
	}
	if len(row)!=3 {
		t.Error("row has incorrect content")
	}
	expected := []string{"content1", "content2", "content3"}
	if !reflect.DeepEqual(row, expected) {
		t.Error("row has incorrect content")
	}

	row = reader.ParseLine("short1,short2")
	if len(reader.GetHeader())!=3 {
		t.Error("header not read correctly")
	}
	if len(row)!=3 {
		t.Error("row has incorrect content")
	}
	expected = []string{"short1", "short2", ""}
	if !reflect.DeepEqual(row, expected) {
		t.Error("row has incorrect content")
	}

	row = reader.ParseLine("long1,long2,long3,long4")
	if len(reader.GetHeader())!=3 {
		t.Error("header not read correctly")
	}
	if len(row)!=3 {
		t.Error("row has incorrect content")
	}
	expected = []string{"long1", "long2", "long3"}
	if !reflect.DeepEqual(row, expected) {
		t.Error("row has incorrect content")
	}

	row = reader.ParseLine(`"quote1","quote2","quote3","quote4"`)
	if len(reader.GetHeader())!=3 {
		t.Error("header not read correctly")
	}
	if len(row)!=3 {
		t.Error("row has incorrect content")
	}
	expected = []string{"quote1", "quote2", "quote3"}
	if !reflect.DeepEqual(row, expected) {
		t.Error("row has incorrect content")
	}

	row = reader.ParseLine(`"quote1a,quote1b","quote2","quote3a,quote3b"`)
	if len(reader.GetHeader())!=3 {
		t.Error("header not read correctly")
	}
	if len(row)!=3 {
		t.Error("row has incorrect content")
	}
	//fmt.Println(row)
	expected = []string{"quote1a,quote1b", "quote2", "quote3a,quote3b"}
	if !reflect.DeepEqual(row, expected) {
		t.Error("row has incorrect content")
	}
}
