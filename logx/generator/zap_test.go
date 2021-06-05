package main

import (
	"go/ast"
	goparser "go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetImportPkg(t *testing.T) {
	r, err := getImportPkg("go.uber.org/zap")
	require.NoError(t, err)
	t.Log(r)
}

func TestParseFunc(t *testing.T) {
	src := `
package test

import (
	"go.uber.org/zap"
)

func bar(a, b, c int) *zap.SugaredLogger {
	x := 1
	x = s.Info(a, b, c)
	return x
}
`
	fset := token.NewFileSet()
	r, err := goparser.ParseFile(fset, "foo.go", src, 0)
	require.NoError(t, err)
	err = ast.Print(fset, r)
	require.NoError(t, err)
	body := r.Decls[1].(*ast.FuncDecl).Body.List[0]
	err = ast.Print(fset, body)
	require.NoError(t, err)
}

func TestParseExpr(t *testing.T) {
	fset := token.NewFileSet()
	r, err := goparser.ParseExpr(`zap.S().Info(a, b, c)`)
	require.NoError(t, err)
	err = ast.Print(fset, r)
	require.NoError(t, err)
}

func TestWalkAst(t *testing.T) {
	dir, err := getImportPkg("go.uber.org/zap")
	require.NoError(t, err)

	tset := token.NewFileSet()
	pkg, err := parseDir(tset, dir, "zap")
	require.NoError(t, err)

	v := &visitor{}
	ast.Walk(v, pkg)
	err = ast.Print(tset, v.funcs)
	require.NoError(t, err)
}
