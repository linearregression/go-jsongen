package test

type RecursiveStruct struct {
	IntField       int
	StrField       string
	RecursiveField *RecursiveStruct
}
