package swityp

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

// Analyzer :
var Analyzer = &analysis.Analyzer{
	Name:     "swityp",
	Doc:      "Check non-exhaustive switch for `type NewType`",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func init() {
	Analyzer.Flags.StringVar(&collector.target, "target", "", "target type: {package path}.{type}")
	Analyzer.Flags.Var(&collector.envs, "env", "environment variables for loading the target pakage")
}

func run(pass *analysis.Pass) (interface{}, error) {
	collector.once.Do(func() {
		collector.err = collector.do()
	})
	if collector.err != nil {
		return nil, errors.Wrap(collector.err, "collect target variables")
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.SwitchStmt)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		t := n.(*ast.SwitchStmt)
		typ := pass.TypesInfo.TypeOf(t.Tag).String()
		if typ != collector.target {
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
		for _, name := range collector.names {
			used[name] = false
		}
		for _, c := range cases {
			switch id := c.(type) {
			case *ast.Ident:
				used[id.Name] = true
			case *ast.SelectorExpr:
				used[id.Sel.Name] = true
			}
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

	return nil, nil
}

type targetCollector struct {
	target string
	envs   stringFlags

	// result
	names []string

	once sync.Once
	err  error
}

var collector = &targetCollector{}

func (c *targetCollector) do() error {
	conf := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo,
	}
	if len(c.envs) > 0 {
		conf.Env = append(os.Environ(), c.envs...)
	}

	pkgPaths := strings.Split(c.target, ".")
	targetPkg := strings.Join(pkgPaths[:len(pkgPaths)-1], ".")
	pkgs, err := packages.Load(conf, targetPkg)
	if err != nil {
		return errors.Wrap(err, "load packages")
	}
	if packages.PrintErrors(pkgs) > 0 {
		return errors.New("load target package")
	}

	names := []string{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				names = append(names, c.collectNames(decl, pkg.TypesInfo)...)
			}
		}
	}
	c.names = names

	return nil
}

func (c *targetCollector) collectNames(decl ast.Decl, info *types.Info) []string {
	t, ok := decl.(*ast.GenDecl)
	if !ok {
		return nil
	}
	if t.Tok != token.VAR {
		return nil
	}

	names := []string{}
	for _, spec := range t.Specs {
		vspec, ok := spec.(*ast.ValueSpec)
		if !ok || len(vspec.Values) == 0 {
			continue
		}

		typ := info.Types[vspec.Values[0]].Type.String()
		if typ != c.target {
			continue
		}
		for _, name := range vspec.Names {
			names = append(names, name.Name)
		}
	}

	return names
}

type stringFlags []string

func (flg *stringFlags) String() string {
	return strings.Join(*flg, ",")
}

func (flg *stringFlags) Set(value string) error {
	*flg = append(*flg, value)
	return nil
}
