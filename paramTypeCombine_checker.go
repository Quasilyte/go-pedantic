package pedantic

import (
	"go/ast"

	"github.com/go-lintpack/lintpack"
	"github.com/go-lintpack/lintpack/astwalk"
	"github.com/go-toolsmith/astequal"
)

func init() {
	var info lintpack.CheckerInfo
	info.Name = "paramTypeCombine"
	info.Tags = []string{"style"}
	info.Summary = "Detects if function parameters could be combined by type and suggest the way to do it"
	info.Before = `func foo(a, b int, c, d int, e, f int, g int) {}`
	info.After = `func foo(a, b, c, d, e, f, g int) {}`

	lintpack.AddChecker(&info, func(ctx *lintpack.CheckerContext) lintpack.FileWalker {
		return astwalk.WalkerForFuncDecl(&paramTypeCombineChecker{ctx: ctx})
	})
}

type paramTypeCombineChecker struct {
	astwalk.WalkHandler
	ctx *lintpack.CheckerContext
}

func (c *paramTypeCombineChecker) EnterFunc(*ast.FuncDecl) bool {
	return true
}

func (c *paramTypeCombineChecker) VisitFuncDecl(decl *ast.FuncDecl) {
	typ := c.optimizeFuncType(decl.Type)
	if !astequal.Expr(typ, decl.Type) {
		c.warn(decl.Type, typ)
	}
}

func (c *paramTypeCombineChecker) optimizeFuncType(f *ast.FuncType) *ast.FuncType {
	return &ast.FuncType{
		Params:  c.optimizeParams(f.Params),
		Results: c.optimizeParams(f.Results),
	}
}
func (c *paramTypeCombineChecker) optimizeParams(params *ast.FieldList) *ast.FieldList {
	// To avoid false positives, skip unnamed param lists.
	//
	// We're using a property that Go only permits unnamed params
	// for the whole list, so it's enough to check whether any of
	// ast.Field have empty name list.
	skip := params == nil ||
		len(params.List) < 2 ||
		len(params.List[0].Names) == 0
	if skip {
		return params
	}

	list := []*ast.Field{}
	names := make([]*ast.Ident, len(params.List[0].Names))
	copy(names, params.List[0].Names)
	list = append(list, &ast.Field{
		Names: names,
		Type:  params.List[0].Type,
	})
	for i, p := range params.List[1:] {
		names = make([]*ast.Ident, len(p.Names))
		copy(names, p.Names)
		if astequal.Expr(p.Type, params.List[i].Type) {
			list[len(list)-1].Names =
				append(list[len(list)-1].Names, names...)
		} else {
			list = append(list, &ast.Field{
				Names: names,
				Type:  params.List[i+1].Type,
			})
		}
	}
	return &ast.FieldList{
		List: list,
	}
}

func (c *paramTypeCombineChecker) warn(f1, f2 *ast.FuncType) {
	c.ctx.Warn(f1, "%s could be replaced with %s", f1, f2)
}
