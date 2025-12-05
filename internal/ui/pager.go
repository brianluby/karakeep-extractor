package ui

import (
	"io"
	"os"
	"os/exec"
)

// UsePager pipes the given output function to `less` if appropriate.
// If not TTY or less fails, it writes directly to the provided defaultWriter.
func UsePager(defaultWriter io.Writer, renderFunc func(io.Writer) error) error {
	// Basic TTY detection (simplistic)
	// For robustness we might check if os.Stdout is a terminal, but standard lib doesn't make it easy without x/term.
	// We'll assume if we are piping to a file (stdout is not a TTY), we shouldn't use pager.
	
	// Check if defaultWriter is actually os.Stdout (or we can just check os.Stdout directly for TTY)
	// Only attempt pager if we are writing to stdout AND stdout is a TTY.
	
	shouldPage := false
	if defaultWriter == os.Stdout {
		fi, _ := os.Stdout.Stat()
		shouldPage = (fi.Mode() & os.ModeCharDevice) != 0
	}

	if !shouldPage {
		// Not a terminal or not writing to stdout, just write to defaultWriter
		return renderFunc(defaultWriter)
	}

	// Try to start `less`
	cmd := exec.Command("less", "-F", "-R", "-X") // -F: quit if one screen, -R: raw colors, -X: no init/deinit
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	pipe, err := cmd.StdinPipe()
	if err != nil {
		// Failed to pipe, fallback
		return renderFunc(defaultWriter)
	}

	if err := cmd.Start(); err != nil {
		// Failed to start less, fallback
		return renderFunc(defaultWriter)
	}

	// Render to pipe
	err = renderFunc(pipe)
	pipe.Close() // Close pipe to signal EOF to less

	if err != nil {
		return err
	}

	// Wait for less to exit
	return cmd.Wait()
}
