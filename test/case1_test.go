package test

import (
	"encoding/json"
	"testing"
)

const inJSON = `
	{
		"IntField": 123,
		"StrField": "Hello world",
	}`

var inJSOBBytes = []byte(inJSON)

func TestParsingJSON(t *testing.T) {
	s := &EasyStruct{}
	err := s.FastUnmarshalJSON(inJSON)

	if err != nil {
		t.Fatal(err)
	}

	if s.IntField != 123 {
		t.Fatalf("Field `IntField` is '%d' but must be '123'", s.IntField)
	}

	if s.StrField != "Hello world" {
		t.Fatalf("Field `StrField` is '%s' but must be 'Hello world'", s.StrField)
	}
}

func BenchmarkStd(b *testing.B) {
	for n := 0; n < b.N; n++ {
		res := &EasyStruct{}
		json.Unmarshal(inJSOBBytes, &res)
	}
}

func BenchmarkFFJson(b *testing.B) {
	for n := 0; n < b.N; n++ {
		res := &EasyStruct{}
		res.UnmarshalFFJSON(inJSOBBytes)
	}
}

func BenchmarkEasyJson(b *testing.B) {
	for n := 0; n < b.N; n++ {
		res := &EasyStruct{}
		res.EasyUnmarshalJSON(inJSOBBytes)
	}
}

func BenchmarkJSONGen(b *testing.B) {
	for n := 0; n < b.N; n++ {
		res := &EasyStruct{}
		res.FastUnmarshalJSON(inJSON)
	}
}
