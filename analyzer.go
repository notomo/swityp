package swityp

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer :
var Analyzer = &analysis.Analyzer{
	Name:     "swityp",
	Doc:      "Check non-exhaustive switch for `type NewType`",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

var target string

func init() {
	Analyzer.Flags.StringVar(&target, "target", "", "target type")
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	names := []string{}
	{
		nodeFilter := []ast.Node{
			(*ast.GenDecl)(nil),
		}
		inspect.Preorder(nodeFilter, func(n ast.Node) {
			t := n.(*ast.GenDecl)
			if t.Tok != token.VAR {
				return
			}

			for _, spec := range t.Specs {
				vspec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				typ := pass.TypesInfo.Types[vspec.Values[0]].Type.String()
				if typ != target {
					continue
				}
				for _, name := range vspec.Names {
					names = append(names, name.Name)
				}
			}
		})
	}

	{
		nodeFilter := []ast.Node{
			(*ast.SwitchStmt)(nil),
		}
		inspect.Preorder(nodeFilter, func(n ast.Node) {
			t := n.(*ast.SwitchStmt)
			typ := pass.TypesInfo.TypeOf(t.Tag).String()
			if typ != target {
				return
			}

			cases := []ast.Expr{}
			for _, stmt := range t.Body.List {
				cc := stmt.(*ast.CaseClause)
				if cc.List == nil {
					// default:
					return
				}
				cases = append(cases, cc.List...)
			}

			used := map[string]bool{}
			for _, name := range names {
				used[name] = false
			}
			for _, c := range cases {
				id, ok := c.(*ast.Ident)
				if !ok {
					continue
				}
				used[id.Name] = true
			}

			unused := []string{}
			for name, use := range used {
				if use {
					continue
				}
				unused = append(unused, name)
			}
			if len(unused) != 0 {
				msg := fmt.Sprintf("non-exhaustive switch: `%s` not covered", strings.Join(unused, ", "))
				pass.Reportf(t.Pos(), msg)
			}
		})
	}

	return nil, nil
}
