package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

/*
READING & WRITING FILES
=======================

Go's file I/O lives mainly in the `os`, `bufio`, and `io` packages.

The simplest, whole-file helpers (great for small files):
  os.WriteFile(name, data, perm)  -> write a whole []byte to a file
  os.ReadFile(name)               -> read a whole file into a []byte

For more control (streaming, appending, large files) you open a *os.File:
  os.Create(name)   -> create/truncate for writing
  os.Open(name)     -> open for reading
  os.OpenFile(...)  -> full control via flags (append, create, read-write...)

Always `defer file.Close()` after opening.

File permissions use Unix octal, e.g. 0644 = owner read/write, others read.
*/

func main() {
	// Work inside a temp dir so the example is self-contained and repeatable.
	dir, err := os.MkdirTemp("", "go-file-rw")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up everything when done

	path := dir + "/languages.txt"

	// --- 1. Write a whole file at once ---
	data := "Go is simple, fast, and fun.\n"
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote file:", path)

	// --- 2. Read a whole file at once ---
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("read back: %q\n", string(content))

	// --- 3. Append to a file using OpenFile + flags ---
	// O_APPEND: write at end, O_CREATE: create if missing, O_WRONLY: write-only.
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	for _, lang := range []string{"C", "Rust", "Python"} {
		if _, err := f.WriteString("- " + lang + "\n"); err != nil {
			f.Close()
			log.Fatal(err)
		}
	}
	f.Close()
	fmt.Println("appended 3 lines")

	// --- 4. Buffered writing with bufio.Writer (efficient for many small writes) ---
	bufPath := dir + "/buffered.txt"
	out, err := os.Create(bufPath)
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(out)
	for i := 1; i <= 5; i++ {
		fmt.Fprintf(w, "line %d\n", i)
	}
	w.Flush() // IMPORTANT: flush the buffer to actually write to disk
	out.Close()
	fmt.Println("wrote buffered file:", bufPath)

	// --- 5. Read a file line by line with bufio.Scanner (best for large files) ---
	in, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	fmt.Println("reading line by line:")
	scanner := bufio.NewScanner(in)
	lineNum := 1
	for scanner.Scan() {
		fmt.Printf("  %d: %s\n", lineNum, scanner.Text())
		lineNum++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
