package test

import (
	"encoding/json"
	"testing"
	"github.com/davecgh/go-spew/spew"
)

const inJSONCase2 = `
	{
		"IntField": 123,
		"StrField": "Hello world",
		"RecursiveField": {
			"IntField": 456,
			"StrField": "Hello world 2",
			"RecursiveField": {
				"IntField": 789,
				"StrField": "Hello world 3",
			}
		}
	}`

var inJSOBBytesCase2 = []byte(inJSONCase2)

func TestParsingJSONCase2(t *testing.T) {
	s := &RecursiveStruct{}
	err := s.FastUnmarshalJSON(inJSONCase2)

	if err != nil {
		t.Fatal(err)
	}

	if s.IntField != 123 {
		t.Fatalf("Field `IntField` is '%d' but must be '123'", s.IntField)
	}

	if s.StrField != "Hello world" {
		t.Fatalf("Field `StrField` is '%s' but must be 'Hello world'", s.StrField)
	}

	if s.RecursiveField.StrField != "Hello world 2" {
		t.Fatalf("Field `RecursiveField.StrField` is '%s' but must be 'Hello world 2'", s.RecursiveField.StrField)
	}

	spew.Dump(s)
}

func BenchmarkCase2Std(b *testing.B) {
	for n := 0; n < b.N; n++ {
		res := &RecursiveStruct{}
		json.Unmarshal(inJSOBBytesCase2, &res)
	}
}

func BenchmarkCase2FFJson(b *testing.B) {
	for n := 0; n < b.N; n++ {
		res := &RecursiveStruct{}
		res.UnmarshalFFJSON(inJSOBBytesCase2)
	}
}

func BenchmarkCase2EasyJson(b *testing.B) {
	for n := 0; n < b.N; n++ {
		res := &RecursiveStruct{}
		res.EasyUnmarshalJSON(inJSOBBytesCase2)
	}
}

func BenchmarkCase2JSONGen(b *testing.B) {
	for n := 0; n < b.N; n++ {
		res := &RecursiveStruct{}
		res.FastUnmarshalJSON(inJSON)
	}
}
