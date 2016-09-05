package generator

import (
	"reflect"

	"bytes"

	"github.com/sergei-svistunov/go-jsongen/parser"
)

func (m *method) addUnmarshal(typeDescr parser.TypeDescription) {
	m.body.WriteString(
		`	pos := 0
	dataLen := len(data)
	if dataLen == 0 {
		return fmt.Errorf("Empty JSON")
	}

	var res `+ typeDescr.GetAlias() + `
`,
	)

	body := &bytes.Buffer{}
	m.addUnmarshalType("res", body, typeDescr)

	for typeName, _ := range m.recursiveParsers {
		m.body.WriteString("var parser_" + typeName + " func(*" + typeName + ") error\n")
	}

	for typeName, body := range m.recursiveParsers {
		m.body.WriteString("parser_" + typeName + " = func(ptr *" + typeName + ") error {\n")
		m.body.WriteString("var res " + typeName + "\n")
		m.body.Write(body.Bytes())
		m.body.WriteString("*ptr = res\nreturn nil\n}\n")
	}

	m.body.Write(body.Bytes())

	m.body.WriteString("*this = res\n")

	m.body.WriteString(`return nil`)
}

func (m *method) addUnmarshalType(variable string, body *bytes.Buffer, typeDescr parser.TypeDescription) {
	if typeDescr.IsRecursive() {
		body.WriteString("parser_" + typeDescr.GetAlias() + "(&" + variable + ")\n")

		if _, exists := m.recursiveParsers[typeDescr.GetAlias()]; exists {
			return
		} else {
			body = &bytes.Buffer{}
			m.recursiveParsers[typeDescr.GetAlias()] = body
		}
	}

	switch t := typeDescr.(type) {
	case *parser.TypePointer:
		m.addUnmarshalPointer(variable, body, t)
	case *parser.TypeInt:
		m.addUnmarshalInt(variable, body, t)
	case *parser.TypeString:
		m.addUnmarshalString(variable, body)
	case *parser.TypeStruct:
		m.addUnmarshalStruct(variable, body, t)
	default:
		println(reflect.TypeOf(t).String())
		panic("Not implemented") // ToDo
	}

}

func (m *method) addUnmarshalSkipSpaces(body *bytes.Buffer) {
	label := m.getNextLabel()
	body.WriteString(label + `:
	    for {
			if pos >= dataLen {
				break ` + label + `
			}

			switch data[pos] {
			case ' ', '\t', '\r', '\n':
				pos++
				continue
			default:
				break ` + label + `
			}
		}
		`)
}

func (m *method) addUnmarshalPointer(variable string, body *bytes.Buffer, typeDescr *parser.TypePointer) {
	body.WriteString(`var tmp	` + typeDescr.Value.GetAlias() + "\n")

	m.addUnmarshalType("tmp", body, typeDescr.Value)

	body.WriteString(variable + " = &tmp\n")
}

func (m *method) addUnmarshalInt(variable string, body *bytes.Buffer, typeDescr *parser.TypeInt) {
	m.addUnmarshalSkipSpaces(body)

	body.WriteString(`
		if pos >= dataLen {
			return fmt.Errorf("Waited digit, got EOF")
		}
		r := pos
		for data[r] >= '0' && data[r] <= '9' {
			r++
			if pos >= dataLen {
				return fmt.Errorf("Waited digit, got EOF")
			}
		}
		if pos == r {
			return fmt.Errorf("Waited digit, got %c", data[pos])
		}
		` + variable + ` = 0
		for ; pos < r; pos++ {
			` + variable + ` *= 10
			` + variable + ` += int(data[pos] - '0')
		}
	`)
}

func (m *method) addUnmarshalStruct(variable string, body *bytes.Buffer, typeDescr *parser.TypeStruct) {
	m.addUnmarshalSkipSpaces(body)

	body.WriteString(`
	if pos >= dataLen {
		return fmt.Errorf("Waited {, got EOF")
	}

	if data[pos] != '{' {
		return fmt.Errorf("Waited {, got %c", data[pos])
	}
	pos++

	`)

	label := m.getNextLabel()
	body.WriteString(label + `:
		for {
	`)
	m.addUnmarshalSkipSpaces(body)
	body.WriteString(`
			if pos >= dataLen {
				return fmt.Errorf("Invalid JSON")
			}

			switch data[pos] {
			case ',':
				pos++
				continue
			case '}':
				pos++
				break ` + label + `
			}

			var fieldName string
	`)

	m.addUnmarshalString("fieldName", body)

	m.addUnmarshalSkipSpaces(body)

	body.WriteString(`
			if pos >= dataLen {
				return fmt.Errorf("Waited :, got EOF")
			}

			if data[pos] != ':' {
				return fmt.Errorf("Waited :, got %c", data[pos])
			}

			pos++
	`)

	body.WriteString(`
			switch fieldName {
	`)

	for _, f := range typeDescr.Fields {
		fname := reflect.StructTag(f.Tag).Get("json")
		if fname == "" {
			fname = f.Name

			body.WriteString(`case "` + fname + `":`)
			m.addUnmarshalType(variable+"."+f.Name, body, f.Type)
		}
	}

	body.WriteString(`
			}
	`)

	body.WriteString(`
		}
	`)
}

func (m *method) addUnmarshalString(variable string, body *bytes.Buffer) {
	m.addUnmarshalSkipSpaces(body)
	body.WriteString(`
		if pos >= dataLen {
			return fmt.Errorf("Waited \", got EOF")
		}

		if data[pos] != '"' {
			return fmt.Errorf("Waited \", got %c", data[pos])
		}
		pos++
		r := pos
		for data[r] != '"' {
			r++
			if pos >= dataLen {
				return fmt.Errorf("Invalid JSON")
			}
		}
		` + variable + ` = data[pos:r]
		pos = r + 1
	`)
}
