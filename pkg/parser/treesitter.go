package parser

import (
	"context"
	"errors"
	"os"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

type Parser interface {
	Parse(context.Context, []byte) (Result, error)
	ParseFile(context.Context, string) (Result, error)
}
type Result struct {
	Package string  `json:"package"`
	File    string  `json:"file"`
	Chunks  []Chunk `json:"chunks"`
}

type Chunk struct {
	StartByte  uint32 `json:"start_byte"`
	EndByte    uint32 `json:"end_byte"`
	StartPoint Point  `json:"start_point"`
	EndPoint   Point  `json:"end_point"`
	Content    []byte `json:"content"`
}

type Import struct {
	Path  []byte
	Start int
	End   int
}

type Point struct {
	Row uint32 `json:"row"`
	Col uint32 `json:"col"`
}
type Function struct {
	Comments   []string
	Name       string
	Block      []byte
	StartPoint Point
	EndPoint   Point
}

// collectPrevSiblingComments collects previous named siblings. If they're comments
// we can consider that the last named sibling is associated with a another named node (structs, interfaces, funcs, etc.)
// It returns the Start Byte of the first comment or the current node if no comments exist.
func collectPrevSiblingComments(node *sitter.Node, input []byte) (sitter.Point, uint32) {
	point := node.StartPoint()
	b := node.StartByte()
	prev := node.PrevNamedSibling()
	for prev != nil && prev.Type() == "comment" {
		b = prev.StartByte()
		point = prev.StartPoint()
		prev = prev.PrevNamedSibling()
	}
	return point, b
}
func ParseFile(ctx context.Context, path string) (*Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	result, err := Parse(ctx, data)
	if err != nil {
		return nil, err
	}
	result.File = path
	return result, nil
}
func Parse(ctx context.Context, input []byte) (*Result, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(golang.GetLanguage())
	tree, err := parser.ParseCtx(ctx, nil, input)
	if err != nil {
		return nil, err
	}
	n := tree.RootNode()
	queryText := `
	(function_declaration) @func
	(type_declaration
		(type_spec
		  (type_identifier))) @type_declaration
	(method_declaration) @method_receiver
	(import_declaration) @imports
	(package_clause (package_identifier)) @package
`
	q, _ := sitter.NewQuery([]byte(queryText), golang.GetLanguage())
	qc := sitter.NewQueryCursor()

	qc.Exec(q, n)
	result := &Result{}
	for {

		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		for _, c := range m.Captures {
			node := c.Node
			if node == nil {
				return nil, errors.New("missing node")
			}
			// 			if node.IsNamed() {
			// 				fmt.Printf(`
			// ------------------------
			// 	Named Node Found:
			// 	node type		: %v
			// 	node prev		: %v
			// 	node next		: %v
			// 	node child_count: %v
			// ------------------------

			// 	`, node.Type(), node.PrevNamedSibling(), node.NextNamedSibling(), node.ChildCount())
			// 			}

			if node.Type() == "package_clause" {
				if node.ChildCount() == 2 {
					result.Package = node.Child(1).Content(input)
				}
			}
			if node.Type() == "function_declaration" ||
				node.Type() == "method_declaration" ||
				node.Type() == "type_declaration" {

				// check to see if the parent is a comment, if it is a comment, iterate over the previous comments, and "attach" it to this node
				startPoint, startComment := collectPrevSiblingComments(node, input)
				end := node.EndByte()

				chunk := Chunk{
					StartByte:  startComment,
					EndByte:    end,
					StartPoint: Point{Row: startPoint.Row, Col: startPoint.Column},
					EndPoint:   Point{Row: node.EndPoint().Row, Col: node.EndPoint().Column},
					Content:    input[startComment:end],
				}
				result.Chunks = append(result.Chunks, chunk)
			}
		}
	}
	return result, nil
}
