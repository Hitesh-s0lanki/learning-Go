package main

import (
	"fmt"
	"path/filepath"
)

/*
FILE PATHS: path/filepath
=========================

NEVER build paths by hand with string concatenation ("dir" + "/" + "file").
The separator differs across OSes ("/" on Linux/macOS, "\" on Windows).
The path/filepath package builds and inspects paths correctly for the current OS.

Key functions:
  filepath.Join(a, b, ...)  -> join parts with the OS separator, cleaned
  filepath.Base(p)          -> last element (the file name)
  filepath.Dir(p)           -> everything but the last element (the directory)
  filepath.Ext(p)           -> file extension, including the dot
  filepath.Clean(p)         -> simplify "." and ".." and duplicate separators
  filepath.Split(p)         -> split into (dir, file)
  filepath.Abs(p)           -> absolute path (resolves against cwd)
  filepath.IsAbs(p)         -> is this an absolute path?
  filepath.Rel(base, targ)  -> relative path from base to target
*/

func main() {
	// Join builds a valid path for the current OS.
	p := filepath.Join("config", "env", "app.yaml")
	fmt.Println("Join:      ", p)

	fmt.Println("Base:      ", filepath.Base(p)) // app.yaml
	fmt.Println("Dir:       ", filepath.Dir(p))  // config/env
	fmt.Println("Ext:       ", filepath.Ext(p))  // .yaml

	// Split returns dir and file separately.
	dir, file := filepath.Split(p)
	fmt.Printf("Split:      dir=%q file=%q\n", dir, file)

	// Clean simplifies messy paths with . and ..
	dirty := "./users/./dir/../other_dir/./file.txt"
	fmt.Println("Clean:     ", filepath.Clean(dirty)) // users/other_dir/file.txt

	// IsAbs / Abs
	fmt.Println("IsAbs:     ", filepath.IsAbs(p))
	abs, _ := filepath.Abs(p)
	fmt.Println("Abs:       ", abs)

	// Rel: how to get from base to target.
	rel, _ := filepath.Rel("config", "config/env/app.yaml")
	fmt.Println("Rel:       ", rel) // env/app.yaml
}
