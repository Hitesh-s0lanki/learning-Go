package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

/*
DIRECTORIES: create, list, walk, remove
========================================

  os.Mkdir(name, perm)       -> create ONE directory (parent must exist)
  os.MkdirAll(name, perm)    -> create a directory AND all missing parents
  os.ReadDir(name)           -> list the entries in a directory (non-recursive)
  filepath.WalkDir(root, fn) -> walk a whole tree recursively
  os.Remove(name)            -> remove ONE empty file/dir
  os.RemoveAll(name)         -> remove a path and everything under it (like rm -rf)

Directory permission 0755 = owner rwx, group/others r-x (standard for dirs).
*/

func main() {
	base, err := os.MkdirTemp("", "go-dirs")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(base) // clean up the whole tree at the end

	// --- Create a nested directory tree in one call ---
	nested := filepath.Join(base, "project", "static", "images")
	if err := os.MkdirAll(nested, 0755); err != nil {
		log.Fatal(err)
	}
	fmt.Println("created:", nested)

	// Drop a few files in so there is something to list/walk.
	os.WriteFile(filepath.Join(base, "project", "README.md"), []byte("# demo"), 0644)
	os.WriteFile(filepath.Join(nested, "logo.png"), []byte("fake-png"), 0644)
	os.WriteFile(filepath.Join(nested, "banner.jpg"), []byte("fake-jpg"), 0644)

	// --- List entries in a single directory (non-recursive) ---
	fmt.Println("\nos.ReadDir of project/static/images:")
	entries, err := os.ReadDir(nested)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range entries {
		fmt.Printf("  %s (dir=%v)\n", e.Name(), e.IsDir())
	}

	// --- Walk the whole tree recursively ---
	fmt.Println("\nfilepath.WalkDir of the whole tree:")
	err = filepath.WalkDir(base, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err // propagate errors (e.g. permission denied)
		}
		// Show a path relative to base so the output is readable.
		rel, _ := filepath.Rel(base, path)
		kind := "file"
		if d.IsDir() {
			kind = "dir "
		}
		fmt.Printf("  [%s] %s\n", kind, rel)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// --- Remove everything (RemoveAll is also handled by the defer above) ---
	fmt.Println("\ncleaning up with os.RemoveAll")
}
