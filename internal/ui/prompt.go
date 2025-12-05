package ui

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"syscall"

	"golang.org/x/term"
)

type Prompt struct {
	reader *bufio.Reader
	writer io.Writer
}

func NewPrompt(r io.Reader, w io.Writer) *Prompt {
	return &Prompt{
		reader: bufio.NewReader(r),
		writer: w,
	}
}

// Ask prompts the user for input.
func (p *Prompt) Ask(label string, defaultValue string) (string, error) {
	prompt := label
	if defaultValue != "" {
		prompt = fmt.Sprintf("%s [%s]", label, defaultValue)
	}
	fmt.Fprintf(p.writer, "%s: ", prompt)

	input, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue, nil
	}
	return input, nil
}

// AskSecret prompts for hidden input.
func (p *Prompt) AskSecret(label string) (string, error) {
	fmt.Fprintf(p.writer, "%s: ", label)
	
	// term.ReadPassword reads directly from stdin typically.
	// For testability, we might need an abstraction, but for this helper we assume stdin/tty.
	// NOTE: This might fail in non-interactive tests.
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(p.writer) // Newline after input
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}

// AskConfirm prompts for a Yes/No confirmation.
func (p *Prompt) AskConfirm(label string) (bool, error) {
	fmt.Fprintf(p.writer, "%s [y/N]: ", label)
	input, err := p.reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes", nil
}

