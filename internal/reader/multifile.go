package reader

import (
	"fmt"
	"io"
)

// MultiReader reads from multiple file paths in sequence,
// tagging each RawLine with its source file name.
type MultiReader struct {
	paths []string
}

// NewMulti creates a MultiReader for the provided file paths.
func NewMulti(paths []string) *MultiReader {
	return &MultiReader{paths: paths}
}

// Lines emits all lines from each file in order.
// Each RawLine's SourceName is set to the originating file path.
// Line numbers reset per file.
func (mr *MultiReader) Lines() (<-chan RawLine, <-chan error) {
	ch := make(chan RawLine, 64)
	errCh := make(chan error, len(mr.paths))

	go func() {
		defer close(ch)
		defer close(errCh)

		for _, path := range mr.paths {
			lr, err := New(path)
			if err != nil {
				errCh <- fmt.Errorf("open %s: %w", path, err)
				continue
			}

			for line := range lr.Lines() {
				line.SourceName = path
				ch <- line
			}

			if cerr := lr.Close(); cerr != nil && cerr != io.EOF {
				errCh <- fmt.Errorf("close %s: %w", path, cerr)
			}
		}
	}()

	return ch, errCh
}
