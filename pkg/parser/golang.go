package parser

import (
	"errors"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

type Golang struct {
	query string
}

func NewGolangChunker() *Golang {
	return &Golang{
		query: `
		(function_declaration) @func
		(type_declaration
			(type_spec
			  (type_identifier))) @type_declaration
		(method_declaration) @method_receiver
		(import_declaration) @imports
		(package_clause (package_identifier)) @package`,
	}
}
func (v *Golang) Query() string {
	return v.query
}

func (v *Golang) GetLanguage() *sitter.Language {
	return golang.GetLanguage()
}

func (v *Golang) isNodeType(node *sitter.Node) bool {
	if node.Type() == "function_declaration" ||
		node.Type() == "method_declaration" ||
		node.Type() == "type_declaration" {
		return true
	}
	return false
}

func (v *Golang) GetPackage(node *sitter.Node, input []byte) string {
	if node.Type() == "package_clause" {
		if node.ChildCount() == 2 {
			return node.Child(1).Content(input)
		}
	}
	return ""
}

func (v *Golang) ChunkNode(node *sitter.Node, input []byte) (*Chunk, error) {
	if node == nil {
		return nil, errors.New("missing node")
	}
	if v.isNodeType(node) {
		return generateChunk(node, input)
	}
	return nil, errors.New("node type not supported")
}
