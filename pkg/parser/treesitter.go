package parser

import (
	"context"
	"errors"
	"log/slog"
	"os"

	sitter "github.com/smacker/go-tree-sitter"
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

	prev = node.PrevNamedSibling()
	for prev != nil {
		slog.Debug(prev.Type())
		prev = prev.PrevNamedSibling()
	}
	return point, b
}
func ParseFile(ctx context.Context, chunker Chunker, path string) (*Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	result, err := Parse(ctx, chunker, data)
	if err != nil {
		return nil, err
	}
	result.File = path
	return result, nil
}
func Parse(ctx context.Context, chunker Chunker, input []byte) (*Result, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(chunker.GetLanguage())
	tree, err := parser.ParseCtx(ctx, nil, input)
	if err != nil {
		return nil, err
	}
	n := tree.RootNode()
	q, _ := sitter.NewQuery([]byte(chunker.Query()), chunker.GetLanguage())
	qc := sitter.NewQueryCursor()

	qc.Exec(q, n)
	result := &Result{}
	for {

		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		m = qc.FilterPredicates(m, input)
		slog.Debug("new match", slog.Attr{Key: "match", Value: slog.AnyValue(m)})
		for _, c := range m.Captures {
			node := c.Node
			if node == nil {
				slog.Warn("missing node")
				return nil, errors.New("missing node")
			}

			pkg := chunker.GetPackage(node, input)
			if pkg != "" {
				result.Package = pkg
			}

			chunk, err := chunker.ChunkNode(node, input)
			if err == nil {
				result.Chunks = append(result.Chunks, *chunk)
			}
		}
	}
	return result, nil
}
