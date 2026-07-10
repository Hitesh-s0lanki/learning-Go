package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
)

/*
EMBEDDING FILES: //go:embed
===========================

Go can bake files DIRECTLY INTO the compiled binary at build time using the
`embed` package and a special //go:embed comment. This means you can ship a
single executable with its templates, static assets, config, or migrations
inside — no external files to deploy.

Three forms:
  //go:embed file.txt      var s string   -> embed one file as a string
  //go:embed file.txt      var b []byte    -> embed one file as bytes
  //go:embed public        var f embed.FS  -> embed a whole directory tree

Rules:
  - The //go:embed comment must sit IMMEDIATELY above the var (no blank line).
  - The path is relative to the .go file, and cannot use ".." to escape it.
  - embed.FS is read-only and satisfies the io/fs.FS interface.
*/

// Embed a single file as a string.
//
//go:embed public/message.txt
var message string

// Embed the entire public/ directory as a virtual, read-only file system.
//
//go:embed public
var assets embed.FS

func main() {
	// 1. The single-file embed is available as a plain string.
	fmt.Println("embedded message.txt:")
	fmt.Println(message)

	// 2. Read a specific file out of the embedded FS.
	data, err := assets.ReadFile("public/data.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("embedded data.txt:")
	fmt.Println(string(data))

	// 3. Walk the embedded FS just like a real directory tree.
	fmt.Println("files embedded under public/:")
	err = fs.WalkDir(assets, "public", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fmt.Println("  -", path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
