package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

/*
FILE INFO, EXISTENCE CHECKS, COPY & RENAME
==========================================

This covers the everyday operations the earlier sections did not:

  os.Stat(name)     -> FileInfo: size, mode/permissions, modtime, IsDir
  existence check   -> errors.Is(err, os.ErrNotExist)
  io.Copy(dst, src) -> stream-copy any reader to any writer (efficient)
  os.Rename(a, b)   -> rename or move a file
  os.Chmod(name, m) -> change permissions

FileInfo is what you get back from Stat; it answers "what is this file like?".
Prefer errors.Is over the older os.IsNotExist for error checks.
*/

func main() {
	dir, err := os.MkdirTemp("", "go-info-copy")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	src := filepath.Join(dir, "source.txt")
	os.WriteFile(src, []byte("copy me!\nline two\n"), 0644)

	// --- 1. os.Stat -> FileInfo ---
	info, err := os.Stat(src)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File info for source.txt:")
	fmt.Println("  name:   ", info.Name())
	fmt.Println("  size:   ", info.Size(), "bytes")
	fmt.Println("  mode:   ", info.Mode())    // e.g. -rw-r--r--
	fmt.Println("  isDir:  ", info.IsDir())
	fmt.Println("  modTime:", info.ModTime().Format("2006-01-02 15:04:05"))

	// --- 2. Check whether a file exists ---
	fmt.Println("\nexistence checks:")
	fmt.Println("  source.txt exists?  ", fileExists(src))
	fmt.Println("  ghost.txt  exists?  ", fileExists(filepath.Join(dir, "ghost.txt")))

	// --- 3. Copy a file by streaming (works for files of any size) ---
	dst := filepath.Join(dir, "copy.txt")
	if err := copyFile(src, dst); err != nil {
		log.Fatal(err)
	}
	fmt.Println("\ncopied source.txt -> copy.txt")

	// --- 4. Rename / move a file ---
	moved := filepath.Join(dir, "renamed.txt")
	if err := os.Rename(dst, moved); err != nil {
		log.Fatal(err)
	}
	fmt.Println("renamed copy.txt -> renamed.txt")

	// --- 5. Change permissions ---
	if err := os.Chmod(moved, 0600); err != nil {
		log.Fatal(err)
	}
	info2, _ := os.Stat(moved)
	fmt.Println("changed mode of renamed.txt to:", info2.Mode())
}

// fileExists reports whether the named file exists, using the modern
// errors.Is check against os.ErrNotExist.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	// Some other error (e.g. permission) — treat as "can't confirm".
	return false
}

// copyFile streams src into dst using io.Copy (constant memory, any size).
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync() // flush to stable storage
}
