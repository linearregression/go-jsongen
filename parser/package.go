package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type Package struct {
	name      string
	typeDecls map[string]ast.Expr
}

func ParseDir(dir string) (*Package, error) {
	fset := token.NewFileSet()

	packages, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return nil, err
	}

	if len(packages) > 1 {
		return nil, fmt.Errorf("Found %d packages", len(packages))
	}

	var pkgName string
	for pkgName, _ = range packages {
		break
	}

	typeDecls := make(map[string]ast.Expr)
	for _, pkgAST := range packages {
		for _, fAST := range pkgAST.Files {
			ast.Inspect(fAST, func(n ast.Node) bool {
				if typeSpec, ok := n.(*ast.Ident); ok && typeSpec.Obj != nil {
					typeDecls[typeSpec.Name] = typeSpec
				}
				return true
			})
		}
	}

	return &Package{
		name:      pkgName,
		typeDecls: typeDecls,
	}, nil
}

func (p *Package) GetName() string {
	return p.name
}

func (p *Package) GetTypeDescription(name string) (TypeDescription, error) {
	typeDecl, exists := p.typeDecls[name]
	if !exists {
		return nil, fmt.Errorf("Type '%s' wasn't found in package '%s'", name, p.name)
	}

	cache := make(map[string]TypeDescription)

	t := astToRpcTypeRecursive(typeDecl, cache)
	t.SetAlias(name)

	return fixDefer(t, cache), nil
}
