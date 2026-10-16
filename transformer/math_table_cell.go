package transformer

import (
	cparser "github.com/kovetskiy/mark/v16/parser"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// MathTableCellTransformer wraps the contents of a table cell holding an inline
// formula in a paragraph.
//
// Confluence sets an ac:image inline only inside a paragraph. Its own editor
// always writes a cell as <td><p>...</p></td>, but a GFM cell holds its inline
// content directly, and an inline formula there was lifted out onto a line of
// its own. Other cells are left alone: a macro expanded into one may be block
// markup that a <p> must not wrap.
type MathTableCellTransformer struct{}

// NewMathTableCellTransformer creates a new MathTableCellTransformer instance.
func NewMathTableCellTransformer() *MathTableCellTransformer {
	return &MathTableCellTransformer{}
}

// Transform implements the parser.ASTTransformer interface.
func (t *MathTableCellTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	_ = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		cell, ok := node.(*extast.TableCell)
		if !entering || !ok {
			return ast.WalkContinue, nil
		}

		if hasInlineMath(cell) {
			paragraph := ast.NewParagraph()
			for child := cell.FirstChild(); child != nil; {
				next := child.NextSibling()
				paragraph.AppendChild(paragraph, child)
				child = next
			}
			cell.AppendChild(cell, paragraph)
		}

		return ast.WalkSkipChildren, nil
	})
}

func hasInlineMath(cell ast.Node) bool {
	found := false
	_ = ast.Walk(cell, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if math, ok := node.(*cparser.Math); entering && ok && !math.Display {
			found = true
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	})

	return found
}
