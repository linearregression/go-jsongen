package generator

import (
	"bytes"
	"strconv"

	"github.com/sergei-svistunov/go-jsongen/parser"
)

type method struct {
	recieverType     string
	imports          []string
	body             *bytes.Buffer
	recursiveParsers map[string]*bytes.Buffer
	labelId          uint64
}

func newMethod(typeDescr parser.TypeDescription, forPointer bool) *method {
	res := &method{
		recieverType: typeDescr.GetAlias(),
		imports: []string{
			"fmt",
		},
		body:             &bytes.Buffer{},
		recursiveParsers: make(map[string]*bytes.Buffer),
	}

	if forPointer {
		res.recieverType = "*" + res.recieverType
	}

	res.addUnmarshal(typeDescr)

	return res
}

func (m *method) getText() []byte {
	buf := bytes.Buffer{}
	buf.WriteString("func (this " + m.recieverType + ") FastUnmarshalJSON(data string) error {\n")

	buf.Write(m.body.Bytes())

	buf.WriteString("}\n")
	return buf.Bytes()
}

func (m *method) getNextLabel() string {
	m.labelId++
	return "L" + strconv.FormatUint(m.labelId, 10)
}
