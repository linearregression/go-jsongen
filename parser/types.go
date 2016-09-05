package parser

import (
	"fmt"
	"go/ast"
	"strings"
)

type (
	TypeDescription interface {
		rpcType()
		SetAlias(string)
		GetAlias() string
		SetRecursive(bool)
		IsRecursive() bool
	}

	TypePointer struct {
		Alias     string
		Value     TypeDescription
		Recursive bool
	}

	TypeUint struct {
		Alias string
		Size  uint8
	}

	TypeInt struct {
		Alias string
		Size  uint8
	}

	TypeFloat struct {
		Alias string
		Size  uint8
	}

	TypeString struct {
		Alias string
	}

	TypeBool struct {
		Alias string
	}

	TypeArray struct {
		Alias     string
		Len       int
		ItemType  TypeDescription
		Recursive bool
	}

	TypeMap struct {
		Alias     string
		KeyType   TypeDescription
		ValueType TypeDescription
		Recursive bool
	}

	TypeStruct struct {
		Alias     string
		Fields    []StructField
		Recursive bool
	}

	TypeExternal struct {
		Alias   string
		Package string
		Name    string
	}

	typeDeferredAlias struct {
		link string
	}

	StructField struct {
		Alias string
		Name  string
		Type  TypeDescription
		Tag   string
	}
)

func (*TypePointer) rpcType()       {}
func (*TypeUint) rpcType()          {}
func (*TypeInt) rpcType()           {}
func (*TypeFloat) rpcType()         {}
func (*TypeBool) rpcType()          {}
func (*TypeString) rpcType()        {}
func (*TypeArray) rpcType()         {}
func (*TypeMap) rpcType()           {}
func (*TypeStruct) rpcType()        {}
func (*TypeExternal) rpcType()      {}
func (*typeDeferredAlias) rpcType() {}

func (t *TypePointer) SetAlias(a string)   { t.Alias = a }
func (t *TypeUint) SetAlias(a string)      { t.Alias = a }
func (t *TypeInt) SetAlias(a string)       { t.Alias = a }
func (t *TypeFloat) SetAlias(a string)     { t.Alias = a }
func (t *TypeBool) SetAlias(a string)      { t.Alias = a }
func (t *TypeString) SetAlias(a string)    { t.Alias = a }
func (t *TypeArray) SetAlias(a string)     { t.Alias = a }
func (t *TypeMap) SetAlias(a string)       { t.Alias = a }
func (t *TypeStruct) SetAlias(a string)    { t.Alias = a }
func (t *TypeExternal) SetAlias(a string)  { t.Alias = a }
func (*typeDeferredAlias) SetAlias(string) {}

func (t *TypePointer) GetAlias() string       { return t.Alias }
func (t *TypeUint) GetAlias() string          { return t.Alias }
func (t *TypeInt) GetAlias() string           { return t.Alias }
func (t *TypeFloat) GetAlias() string         { return t.Alias }
func (t *TypeBool) GetAlias() string          { return t.Alias }
func (t *TypeString) GetAlias() string        { return t.Alias }
func (t *TypeArray) GetAlias() string         { return t.Alias }
func (t *TypeMap) GetAlias() string           { return t.Alias }
func (t *TypeStruct) GetAlias() string        { return t.Alias }
func (t *TypeExternal) GetAlias() string      { return t.Alias }
func (t *typeDeferredAlias) GetAlias() string { return t.link }

func (t *TypePointer) SetRecursive(r bool)   { t.Recursive = r }
func (t *TypeUint) SetRecursive(bool)        {}
func (t *TypeInt) SetRecursive(bool)         {}
func (t *TypeFloat) SetRecursive(bool)       {}
func (t *TypeBool) SetRecursive(bool)        {}
func (t *TypeString) SetRecursive(bool)      {}
func (t *TypeArray) SetRecursive(r bool)     { t.Recursive = r }
func (t *TypeMap) SetRecursive(r bool)       { t.Recursive = r }
func (t *TypeStruct) SetRecursive(r bool)    { t.Recursive = r }
func (t *TypeExternal) SetRecursive(bool)    {}
func (*typeDeferredAlias) SetRecursive(bool) {}

func (t *TypePointer) IsRecursive() bool       { return t.Recursive }
func (t *TypeUint) IsRecursive() bool          { return false }
func (t *TypeInt) IsRecursive() bool           { return false }
func (t *TypeFloat) IsRecursive() bool         { return false }
func (t *TypeBool) IsRecursive() bool          { return false }
func (t *TypeString) IsRecursive() bool        { return false }
func (t *TypeArray) IsRecursive() bool         { return t.Recursive }
func (t *TypeMap) IsRecursive() bool           { return t.Recursive }
func (t *TypeStruct) IsRecursive() bool        { return t.Recursive }
func (t *TypeExternal) IsRecursive() bool      { return false }
func (t *typeDeferredAlias) IsRecursive() bool { return true }

func astToRpcTypeRecursive(expr ast.Expr, aliasCache map[string]TypeDescription) TypeDescription {
	switch expr := expr.(type) {
	case *ast.Ident:
		switch expr.Name {
		case "int":
			return &TypeInt{
				Size: 0,
			}
		case "int64":
			return &TypeInt{
				Size: 64,
			}
		case "int32":
			return &TypeInt{
				Size: 32,
			}
		case "int16":
			return &TypeInt{
				Size: 16,
			}
		case "int8":
			return &TypeInt{
				Size: 8,
			}
		case "uint":
			return &TypeUint{
				Size: 0,
			}
		case "uint64":
			return &TypeUint{
				Size: 64,
			}
		case "uint32":
			return &TypeUint{
				Size: 32,
			}
		case "uint16":
			return &TypeUint{
				Size: 16,
			}
		case "uint8":
			return &TypeUint{
				Size: 8,
			}
		case "float32":
			return &TypeFloat{
				Size: 32,
			}
		case "float64":
			return &TypeFloat{
				Size: 64,
			}
		case "bool":
			return &TypeBool{}
		case "string":
			return &TypeString{}
		default:
			if expr.Obj != nil {
				if t, exists := aliasCache[expr.Name]; exists {
					return t
				} else {
					aliasCache[expr.Name] = &typeDeferredAlias{
						link: expr.Name,
					}
					res := astToRpcTypeRecursive(expr.Obj.Decl.(*ast.TypeSpec).Type, aliasCache)
					res.SetAlias(expr.Name)
					aliasCache[expr.Name] = res

					return res
				}
			} else {
				fmt.Printf("%#v\n", expr)
				panic(expr)
			}
		}
	case *ast.StarExpr:
		return &TypePointer{
			Value: astToRpcTypeRecursive(expr.X, aliasCache),
		}
	case *ast.StructType:
		res := &TypeStruct{
			Fields: make([]StructField, expr.Fields.NumFields()),
		}
		for i, astField := range expr.Fields.List {
			res.Fields[i].Name = astField.Names[0].Name
			res.Fields[i].Type = astToRpcTypeRecursive(astField.Type, aliasCache)
			if astField.Tag != nil {
				res.Fields[i].Tag = strings.Trim(astField.Tag.Value, "`")
			}
		}
		return res
	case *ast.ArrayType:
		arrLen := 0
		if expr.Len != nil {
			fmt.Printf("%#v\n", expr.Len)
			panic("Not implemented") // ToDo: fix it
		}
		return &TypeArray{
			ItemType: astToRpcTypeRecursive(expr.Elt, aliasCache),
			Len:      arrLen,
		}
	case *ast.MapType:
		return &TypeMap{
			KeyType:   astToRpcTypeRecursive(expr.Key, aliasCache),
			ValueType: astToRpcTypeRecursive(expr.Value, aliasCache),
		}
	case *ast.SelectorExpr:
		return &TypeExternal{
			Package: expr.X.(*ast.Ident).Name,
			Name:    expr.Sel.Name,
		}
	default:
		fmt.Printf("%#v\n", expr)
		panic(expr)
	}
}

func fixDefer(t TypeDescription, cache map[string]TypeDescription) TypeDescription {
	switch t := t.(type) {
	case *typeDeferredAlias:
		res := cache[t.link]
		res.SetRecursive(true)
		return res
	case *TypePointer:
		t.Value = fixDefer(t.Value, cache)
	case *TypeStruct:
		for i, _ := range t.Fields {
			t.Fields[i].Type = fixDefer(t.Fields[i].Type, cache)
		}
		return t
	case *TypeMap:
		t.KeyType = fixDefer(t.KeyType, cache)
		t.ValueType = fixDefer(t.ValueType, cache)
		return t
	case *TypeArray:
		t.ItemType = fixDefer(t.ItemType, cache)
		return t
	}

	return t
}
