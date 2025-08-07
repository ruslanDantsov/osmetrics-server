package main

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "запрещает вызов os.Exit в main.main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if isGenerated(pass, file) {
			continue
		}

		if pass.Pkg.Name() != "main" {
			continue
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "main" || fn.Recv != nil {
				continue
			}

			ast.Inspect(fn.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				pkgIdent, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}

				if pkgIdent.Name == "os" && sel.Sel.Name == "Exit" {
					pass.Reportf(call.Pos(), "Deny to use os.Exit in func main in package main")
				}

				return true
			})
		}
	}
	return nil, nil
}

func isGenerated(pass *analysis.Pass, file *ast.File) bool {
	for _, group := range file.Comments {
		for _, comment := range group.List {
			if strings.HasPrefix(comment.Text, "// Code generated") &&
				strings.Contains(comment.Text, "DO NOT EDIT") {
				return true
			}
		}
	}
	return false
}
