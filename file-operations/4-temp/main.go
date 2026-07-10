package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

/*
TEMPORARY FILES & DIRECTORIES
=============================

Temp files/dirs are perfect for scratch data that should not stick around:
intermediate downloads, test fixtures, generated output, etc. The OS puts them
in a system temp location (os.TempDir(), e.g. /tmp on Linux/macOS).

  os.CreateTemp(dir, pattern)  -> create a new, uniquely-named temp FILE
  os.MkdirTemp(dir, pattern)   -> create a new, uniquely-named temp DIR

Passing "" as dir uses the default system temp directory.
A "*" in the pattern is replaced with the random part; otherwise it is a suffix.

GOLDEN RULE: always clean up with defer os.Remove / os.RemoveAll.
*/

func main() {
	// --- Temp file ---
	// Pattern "logs-*.txt" -> file named like logs-1234567890.txt
	tmpFile, err := os.CreateTemp("", "logs-*.txt")
	if err != nil {
		log.Fatal(err)
	}
	// Clean up the temp file no matter how we exit.
	defer func() {
		fmt.Println("removing temp file:", tmpFile.Name())
		os.Remove(tmpFile.Name())
	}()

	fmt.Println("created temp file:", tmpFile.Name())
	if _, err := tmpFile.WriteString("some scratch log data\n"); err != nil {
		tmpFile.Close()
		log.Fatal(err)
	}
	tmpFile.Close()

	// --- Temp directory ---
	tmpDir, err := os.MkdirTemp("", "my-app-*")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		fmt.Println("removing temp dir: ", tmpDir)
		os.RemoveAll(tmpDir)
	}()

	fmt.Println("created temp dir: ", tmpDir)

	// Put a file inside the temp dir.
	inside := filepath.Join(tmpDir, "cache.bin")
	if err := os.WriteFile(inside, []byte("cached bytes"), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote file inside temp dir:", inside)

	fmt.Println("system temp dir is:", os.TempDir())
}
