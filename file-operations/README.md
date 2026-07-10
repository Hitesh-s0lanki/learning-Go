# File Operations in Go

A hands-on tour of working with files, paths, directories, temp storage, and
embedded assets in Go. Each folder is a standalone `package main` you can run
on its own.

## How to run

```bash
cd file-operations/1-read-write
go run .
```

Examples write only to the system temp directory (or clean up after
themselves), so running them leaves your repo clean.

## Topics (learn them in order)

| # | Folder | Concept |
|---|--------|---------|
| 1 | [1-read-write](1-read-write/main.go) | `os.ReadFile`/`WriteFile`, `os.Create`/`Open`/`OpenFile` flags, `bufio.Writer`, `bufio.Scanner` |
| 2 | [2-paths](2-paths/main.go) | `path/filepath`: `Join`, `Base`, `Dir`, `Ext`, `Clean`, `Split`, `Abs`, `Rel` |
| 3 | [3-directories](3-directories/main.go) | `Mkdir`/`MkdirAll`, `os.ReadDir`, `filepath.WalkDir`, `Remove`/`RemoveAll` |
| 4 | [4-temp](4-temp/main.go) | `os.CreateTemp`, `os.MkdirTemp`, cleanup with `defer` |
| 5 | [5-embed](5-embed/main.go) | `//go:embed` — bake files into the binary (`string`, `embed.FS`) |
| 6 | [6-info-copy](6-info-copy/main.go) | `os.Stat`/`FileInfo`, existence checks, `io.Copy`, `os.Rename`, `os.Chmod` |

## Key packages

- **`os`** — the core: open/create/read/write files, stat, remove, temp,
  permissions.
- **`io`** — generic streaming (`io.Copy`, readers/writers) that works with
  files, network connections, buffers — anything.
- **`bufio`** — buffered reading (`Scanner` for line-by-line) and writing
  (`Writer`) for efficiency.
- **`path/filepath`** — build and inspect OS-correct file paths. Never
  concatenate paths by hand.
- **`embed`** — compile files into the binary at build time.

## Key habits

- Always `defer file.Close()` right after a successful open/create.
- Always `Flush()` a `bufio.Writer` (or its data may never hit disk).
- Build paths with `filepath.Join`, not string `+`, so they work on every OS.
- Check "does it exist?" with `errors.Is(err, os.ErrNotExist)`.
- Clean up temp files/dirs with `defer os.Remove` / `os.RemoveAll`.
- Permissions are Unix octal: `0644` for files, `0755` for directories.

## What this adds beyond the basics

Compared to a minimal walkthrough, these examples also cover the operations you
hit constantly in real code but are easy to miss:

- **Buffered writing** with `bufio.Writer` (and why you must `Flush`).
- **Copying files** by streaming with `io.Copy` (constant memory, any size).
- **File metadata** via `os.Stat` → `FileInfo` (size, mode, modtime, isDir).
- **Existence checks** the modern way with `errors.Is(err, os.ErrNotExist)`.
- **Listing and walking** directories (`os.ReadDir`, `filepath.WalkDir`).
- **Rename/move** and **chmod**.
