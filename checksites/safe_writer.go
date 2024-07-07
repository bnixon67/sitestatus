package main

import (
	"io"
	"sync"
)

// SafeWriter ensures thread-safe writes to an underlying writer.
type SafeWriter struct {
	mu     sync.Mutex
	Writer io.Writer
}

// Write implements io.Writer interface for SafeWriter to ensure thread safety.
func (sw *SafeWriter) Write(p []byte) (n int, err error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.Writer.Write(p)
}
