package parser

import (
	"fmt"
	"os"

	"github.com/go-enry/go-enry/v2"
	sitter "github.com/smacker/go-tree-sitter"
)

// DetectLanguage uses enry for language detection.
// TODO: provide support for detecting multiple languages inside of a single file (e.g. JS/HTML, HEEX support, etc.)
func DetectLanguageFromFile(filepath string) (Chunker, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	language := enry.GetLanguage(filepath, data)

	switch language {
	case "Go":
		return NewGolangChunker(), nil
	case "Python":
		return NewPythonChunker(), nil
	default:
		return nil, fmt.Errorf("unable to get chunker for %v", language)
	}
}

// DetectLanguage uses enry for language detection.
// TODO: provide support for detecting multiple languages inside of a single file (e.g. JS/HTML, HEEX support, etc.)
func DetectLanguage(filepath string, data []byte) (Chunker, error) {
	language := enry.GetLanguage(filepath, data)
	switch language {
	case "Go":
		return NewGolangChunker(), nil
	case "Python":
		return NewPythonChunker(), nil
	default:
		return nil, fmt.Errorf("unable to get chunker for %v", language)
	}
}

type Chunker interface {
	Query() string
	ChunkNode(*sitter.Node, []byte) (*Chunk, error)
	GetPackage(*sitter.Node, []byte) string
	GetLanguage() *sitter.Language
}

func generateChunk(node *sitter.Node, input []byte) (*Chunk, error) {
	startPoint, startComment := collectPrevSiblingComments(node, input)
	end := node.EndByte()
	chunk := &Chunk{
		StartByte:  startComment,
		EndByte:    end,
		StartPoint: Point{Row: startPoint.Row, Col: startPoint.Column},
		EndPoint:   Point{Row: node.EndPoint().Row, Col: node.EndPoint().Column},
		Content:    input[startComment:end],
	}
	return chunk, nil
}
