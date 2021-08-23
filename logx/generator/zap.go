package main

import (
	"bytes"
	"fmt"
	"go/ast"
	gobuild "go/build"
	"go/format"
	goparser "go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/win5do/go-lib/errx"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	dir, err := getImportPkg("go.uber.org/zap")
	if err != nil {
		return errx.WithStackOnce(err)
	}
	log.Printf("dir: %+v", dir)

	pkg, err := parseDir(token.NewFileSet(), dir, "zap")
	if err != nil {
		return errx.WithStackOnce(err)
	}

	funcs, err := walkAst(pkg)
	if err != nil {
		return errx.WithStackOnce(err)
	}

	err = writeGoFile(os.Stdout, funcs)
	if err != nil {
		return errx.WithStackOnce(err)
	}

	return nil
}

func getImportPkg(pkg string) (string, error) {
	p, err := gobuild.Import(pkg, "", gobuild.FindOnly)
	if err != nil {
		return "", err
	}

	return p.Dir, err

}

func parseDir(tset *token.FileSet, dir, pkgName string) (*ast.Package, error) {
	pkgMap, err := goparser.ParseDir(
		tset,
		dir,
		func(info os.FileInfo) bool {
			// skip go-test
			return !strings.Contains(info.Name(), "_test.go")
		},
		goparser.Mode(0), // no comment
	)
	if err != nil {
		return nil, errx.WithStackOnce(err)
	}

	pkg, ok := pkgMap[pkgName]
	if !ok {
		err := errors.New("not found")
		return nil, errx.WithStackOnce(err)
	}

	return pkg, nil
}

type visitor struct {
	funcs []*ast.FuncDecl
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Recv == nil ||
			!n.Name.IsExported() ||
			len(n.Recv.List) != 1 {
			return nil
		}
		t, ok := n.Recv.List[0].Type.(*ast.StarExpr)
		if !ok {
			return nil
		}

		if t.X.(*ast.Ident).String() != "SugaredLogger" {
			return nil
		}

		log.Printf("func name: %s", n.Name.String())

		v.funcs = append(v.funcs, n)
	}
	return v
}

func walkAst(node ast.Node) ([]ast.Decl, error) {
	v := &visitor{}
	ast.Walk(v, node)

	log.Printf("funcs len: %d", len(v.funcs))

	var decls []ast.Decl
	for _, v := range v.funcs {
		decls = append(decls, rewriteFunc(v))
	}

	return decls, nil
}

func rewriteFunc(fn *ast.FuncDecl) *ast.FuncDecl {
	fn.Recv = nil

	fnName := fn.Name.String()

	// handle return: *Logger -> *zap.Logger
	if fn.Type.Results != nil {
		for _, result := range fn.Type.Results.List {
			star, ok := result.Type.(*ast.StarExpr)
			if !ok {
				continue
			}
			ident, ok := star.X.(*ast.Ident)
			if !ok || !strings.Contains(ident.Name, "Logger") {
				continue
			}

			exprStr := fmt.Sprintf(`*zap.%s`, ident.Name)
			expr, err := goparser.ParseExpr(exprStr)
			if err != nil {
				panic(err)
			}
			result.Type = expr
		}
	}

	// handle body: _globalS.Info(args...)
	var args []string
	for _, field := range fn.Type.Params.List {
		for _, id := range field.Names {
			idStr := id.String()
			_, ok := field.Type.(*ast.Ellipsis)
			if ok {
				// Ellipsis args
				idStr += "..."
			}
			args = append(args, idStr)
		}
	}

	exprStr := fmt.Sprintf(`getSugarLogger().%s(%s)`, fnName, strings.Join(args, ","))
	expr, err := goparser.ParseExpr(exprStr)
	if err != nil {
		panic(err)
	}

	var body []ast.Stmt
	if fn.Type.Results != nil {
		body = []ast.Stmt{
			&ast.ReturnStmt{
				// Return:
				Results: []ast.Expr{expr},
			},
		}
	} else {
		body = []ast.Stmt{
			&ast.ExprStmt{
				X: expr,
			},
		}
	}

	fn.Body.List = body

	return fn
}

func astToGo(dst *bytes.Buffer, node interface{}) error {
	addNewline := func() {
		err := dst.WriteByte('\n') // add newline
		if err != nil {
			log.Panicln(err)
		}
	}

	addNewline()

	err := format.Node(dst, token.NewFileSet(), node)
	if err != nil {
		return err
	}

	addNewline()

	return nil
}

// Output Go code
func writeGoFile(wr io.Writer, funcs []ast.Decl) error {
	header := `// Code generated by log-gen. DO NOT EDIT.
package logx

import (
	"go.uber.org/zap"
)
`
	buffer := bytes.NewBufferString(header)

	for _, fn := range funcs {
		err := astToGo(buffer, fn)
		if err != nil {
			return errx.WithStackOnce(err)
		}
	}

	_, err := wr.Write(buffer.Bytes())
	return err
}
