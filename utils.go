package pedantic

import (
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"
)

// TODO(Quasilyte): maybe it's OK to make a new go-toolsmith
// package that helps to reuse some code between several
// go-lintpack packages like pedantic and go-critic/checkers.
//
// For now, I'm doing a copy/paste of some util functions
// from the go-critic/checkers package.

// findNode applies pred for root and all it's childs until it returns true.
// Matched node is returned.
// If none of the nodes matched predicate, nil is returned.
func findNode(root ast.Node, pred func(ast.Node) bool) ast.Node {
	var found ast.Node
	astutil.Apply(root, nil, func(cur *astutil.Cursor) bool {
		if pred(cur.Node()) {
			found = cur.Node()
			return false
		}
		return true
	})
	return found
}

// containsNode reports whether `findNode(root, pred)!=nil`.
func containsNode(root ast.Node, pred func(ast.Node) bool) bool {
	return findNode(root, pred) != nil
}
