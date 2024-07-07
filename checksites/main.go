// Copyright 2024 Bill Nixon. All rights reserved.
// Use of this source code is governed by the license found in the LICENSE file.

package main

import (
	"flag"
	"os"
)

func main() {
	insecure := flag.Bool("insecure", false, "Ignore certificates")
	noRedirects := flag.Bool("noRedirects", false, "Ignore redirects")
	timeout := flag.Int("timeout", 15, "Timeout in seconds for each URL")
	flag.Parse()

	readAndCheck(os.Stdin, os.Stdout, *insecure, *noRedirects, *timeout)
}
