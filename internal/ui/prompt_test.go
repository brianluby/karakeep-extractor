package ui

import (
	"bytes"
	"testing"
)

func TestPrompt_Ask(t *testing.T) {
	input := "my input\n"
	r := bytes.NewBufferString(input)
	var w bytes.Buffer

	p := NewPrompt(r, &w)
	
	res, err := p.Ask("Question", "default")
	if err != nil {
		t.Fatalf("Ask failed: %v", err)
	}

	if res != "my input" {
		t.Errorf("Expected 'my input', got '%s'", res)
	}
}

func TestPrompt_Ask_Default(t *testing.T) {
	input := "\n"
	r := bytes.NewBufferString(input)
	var w bytes.Buffer

	p := NewPrompt(r, &w)
	
	res, err := p.Ask("Question", "default")
	if err != nil {
		t.Fatalf("Ask failed: %v", err)
	}

	if res != "default" {
		t.Errorf("Expected 'default', got '%s'", res)
	}
}
