package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bnixon67/sitestatus"
)

// checkSite checks the status of a site and writes the result to w.
func checkSite(w io.Writer, wg *sync.WaitGroup, url string, insecure, noRedirects bool, timeout time.Duration) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		opts := sitestatus.HTTPClientOptions{
			IgnoreCerts:     insecure,
			IgnoreRedirects: noRedirects,
			Timeout:         timeout,
		}
		result := sitestatus.Check(url, opts)

		fmt.Fprintln(w, url, result)
	}()
}

// readAndCheck reads URLs and checks their status concurrently via goroutines.
func readAndCheck(r io.Reader, w io.Writer, insecure, noRedirects bool, timeout int) {
	dur := time.Duration(timeout) * time.Second

	var wg sync.WaitGroup

	// Use SafeWriter to avoid corrupted output in concurrent writes.
	safeWriter := &SafeWriter{Writer: w}

	// Scan input line by line.
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if sitestatus.IsValidURL(line) {
			checkSite(safeWriter, &wg, line, insecure, noRedirects, dur)
		} else {
			fmt.Fprintln(safeWriter, line, "invalid URL")
		}
	}

	wg.Wait()

	// Display any errors from reading input.
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}

}
