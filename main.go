package main

import (
	"github.com/sergei-svistunov/go-jsongen/parser"
	"github.com/sergei-svistunov/go-jsongen/generator"
)

func main() {
	pkg, err := parser.ParseDir("/home/svistunov/goprojects/src/github.com/sergei-svistunov/go-jsongen/test")
	if err != nil {
		panic(err)
	}

	//spew.Dump(pkg.GetTypeDescription("EasyStruct"))

	gen := generator.NewGenerator()
	gen.SetPackageName(pkg.GetName())

	//td, err := pkg.GetTypeDescription("EasyStruct")
	td, err := pkg.GetTypeDescription("RecursiveStruct")
	if err != nil {
		panic(err)
	}
	gen.AddTypeDescription(td)

	//println(string(gen.GetText()))

	gen.WriteTo("/home/svistunov/goprojects/src/github.com/sergei-svistunov/go-jsongen/test/case2_jsongen.go")
}
