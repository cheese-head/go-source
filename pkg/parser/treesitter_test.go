package parser_test

import (
	"context"
	"encoding/json"
	"go-source/pkg/parser"
	"log"
	"os"
	"testing"
)

// Name struct
type Name struct {
	N string
}

// implements hello
func (v *Name) Hello() {

}

/*
some comments
*/
// hello orl
type Namer interface {
	Hello()
}

// just a sample comment

// hello
// okay
// TestParser is responsible for testing the function
func TestParser(t *testing.T) {

	body, err := os.ReadFile("treesitter_test.go")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	result, err := parser.Parse(context.Background(), body)
	if err != nil {
		log.Fatalf("unable to parse file: %v", err)
	}
	data, err := json.Marshal(result)
	if err != nil {
		log.Fatal("unable to marshal result")
	}
	log.Default().Println(string(data))
}

// sample
func sample(name string) string {
	return ""
}

func nocomment() {

}
