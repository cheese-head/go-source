package parser

import (
	"errors"
	"log/slog"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"
)

type Python struct {
	query string
}

func NewPythonChunker() *Python {
	return &Python{
		query: `
		(class_definition .
			name: (identifier)
			body: (block 
					(function_definition) )
		) @definition.class
		(function_definition) @definition.function

	`,
	}
}
func (v *Python) Query() string {
	return v.query
}

func (v *Python) GetLanguage() *sitter.Language {
	return python.GetLanguage()
}

func (v *Python) isNodeType(node *sitter.Node) bool {

	if node.Type() == "function_definition" {
		parent := node.Parent()
		for parent != nil {
			if parent.Type() == "class_definition" {
				return false
			}
			parent = parent.Parent()
		}
		return true
	}
	if node.Type() == "class_definition" {
		return true
	}
	return false
}

func (v *Python) GetPackage(node *sitter.Node, input []byte) string {
	return ""
}

func (v *Python) ChunkNode(node *sitter.Node, input []byte) (*Chunk, error) {
	slog.Debug("chunking now")
	if node == nil {
		return nil, errors.New("missing node")
	}

	if v.isNodeType(node) {
		return generateChunk(node, input)
	}
	return nil, errors.New("node type not supported")
}
