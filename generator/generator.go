package generator

import (
	"bytes"
	"github.com/sergei-svistunov/go-jsongen/parser"
	"go/format"
	"os"
)

type Generator struct {
	pkgName string
	imports map[string]struct{}
	methods []*method
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) SetPackageName(name string) {
	g.pkgName = name
}

func (g *Generator) AddTypeDescription(typeDescr parser.TypeDescription) {
	g.methods = append(g.methods, newMethod(typeDescr, true))
}

func (g *Generator) GetText() []byte {
	buf := &bytes.Buffer{}

	buf.WriteString("package " + g.pkgName + "\n\n")

	importsMap := make(map[string]struct{})
	for _, m := range g.methods {
		for _, i := range m.imports {
			importsMap[i] = struct{}{}
		}
	}
	buf.WriteString("import (\n")
	for i := range importsMap {
		buf.WriteString("\"" + i + "\"\n")
	}
	buf.WriteString(")\n\n")

	for _, m := range g.methods {
		buf.Write(m.getText())
		buf.WriteByte('\n')
	}

	res, err := format.Source(buf.Bytes())
	if err != nil {
		println(string(buf.Bytes()))
		panic(err)
	}

	return res
}

func (g *Generator) WriteTo(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(g.GetText())

	return nil
}
