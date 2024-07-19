package exitcall

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "exitcall",
	Doc:      "checks for a direct call os.Exit in the main package's main function.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	i := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Fast path: if the package is not main and if the package doesn't import "os"
	// skip the traversal.
	if pass.Pkg.Name() != "main" && !imports(pass.Pkg, "os") {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	i.WithStack(nodeFilter, func(node ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return true
		}

		call := node.(*ast.CallExpr)
		if !isCallExit(pass.TypesInfo, call) {
			return true
		}

		if isCallDirectInMain(stack) {
			pass.ReportRangef(node, "direct call os.Exit in the main package's main function")
			return false
		}

		return true
	})

	return nil, nil
}

func isCallExit(info *types.Info, expr *ast.CallExpr) bool {
	fun, ok := expr.Fun.(*ast.SelectorExpr)
	if !ok || fun.Sel.Name != "Exit" {
		return false
	}

	typ := info.Types[fun.X].Type
	if typ == nil {
		id, ok := fun.X.(*ast.Ident)
		return ok && id.Name == "os" // function in os package
	}

	return false
}

func isCallDirectInMain(stack []ast.Node) bool {
	for i := len(stack) - 2; i >= 0; i-- {
		switch parentNode := stack[i].(type) {
		case *ast.FuncDecl:
			return parentNode.Name.Name != "main" // call os.Exit inside from function main or not
		case *ast.CallExpr:
			return false // call os.Exit inside from some other function
		case *ast.FuncLit:
			return false // call os.Exit inside from anonymous function
		}
	}

	return false
}

// Imports returns true if path is imported by pkg.
func imports(pkg *types.Package, path string) bool {
	for _, imp := range pkg.Imports() {
		if imp.Path() == path {
			return true
		}
	}
	return false
}
