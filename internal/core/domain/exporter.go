package domain

import "io"

// Exporter defines the interface for exporting ranked repositories.
type Exporter interface {
	Export(repos []ExtractedRepo, w io.Writer) error
}
