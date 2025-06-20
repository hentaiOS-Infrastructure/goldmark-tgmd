package tgmd

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark/util" // Keep for util.BufWriter interface
)

// testBufWriter is a wrapper around bytes.Buffer to satisfy util.BufWriter for testing.
type testBufWriter struct {
	*bytes.Buffer
}

// Buffered implements util.BufWriter.
func (w *testBufWriter) Buffered() int {
	return w.Len()
}

// Available implements util.BufWriter.
// For a bytes.Buffer, available capacity can be considered large or effectively Cap - Len.
// Returning a large number if Cap is 0 or very small, or simply Cap - Len.
func (w *testBufWriter) Available() int {
	if w.Cap() == 0 {
		return 0 // Or a reasonably large number if appropriate for tests
	}
	return w.Cap() - w.Len()
}

// Flush implements util.BufWriter. For bytes.Buffer, it's a no-op.
func (w *testBufWriter) Flush() error {
	return nil
}

// Ensure testBufWriter implements util.BufWriter.
var _ util.BufWriter = (*testBufWriter)(nil)

func TestWriteCustomBytes_RawEscaping(t *testing.T) {
	input := []byte("Test Update: AAAA.000000.001")
	expected := "Test Update: AAAA\\.000000\\.001"

	// Use our testBufWriter
	b := &bytes.Buffer{}
	writer := &testBufWriter{Buffer: b}

	writeCustomBytes(writer, input)

	if writer.String() != expected { // Use writer.String() which delegates to b.String()
		t.Errorf("Output mismatch for writeCustomBytes:\nExpected: %q\nGot:      %q", expected, writer.String())
	}
}
