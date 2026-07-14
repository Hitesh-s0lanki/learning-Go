# Serialization in Go (JSON & Base64)

A hands-on tour of turning Go values into portable data and back again:
marshalling/unmarshalling JSON, streaming with encoders/decoders, and encoding
binary data as Base64. Each folder is a standalone `package main` you can run
on its own.

## How to run

```bash
cd serialization/1-marshalling
go run .
```

Every example prints to stdout and touches nothing on disk, so running them
leaves your repo clean.

## Topics (learn them in order)

| # | Folder | Concept |
|---|--------|---------|
| 1 | [1-marshalling](1-marshalling/main.go) | `json.Marshal` / `MarshalIndent`, struct tags, `omitempty`, `-`, slices & maps |
| 2 | [2-unmarshal](2-unmarshal/main.go) | `json.Unmarshal` into structs, `map[string]interface{}`, nested & unknown fields |
| 3 | [3-encoder](3-encoder/main.go) | `json.Encoder`: streaming to any `io.Writer`, `SetIndent`, `SetEscapeHTML`, NDJSON |
| 4 | [4-decoder](4-decoder/main.go) | `json.Decoder`: streaming from any `io.Reader`, `DisallowUnknownFields`, EOF loop |
| 5 | [5-base64](5-base64/main.go) | `encoding/base64`: `Std`/`URL`/`Raw` encodings, encode/decode round-trips |

## Marshal vs. Encoder — which one?

|  | In-memory (`Marshal`/`Unmarshal`) | Streaming (`Encoder`/`Decoder`) |
|--|--|--|
| Works on | `[]byte` | `io.Reader` / `io.Writer` |
| Best for | small values, quick round-trips | large data, files, HTTP bodies, many values |
| Memory | holds the whole value | processes as it goes |
| Typical use | `json.Marshal(v)` | `json.NewEncoder(w).Encode(v)` |

In real HTTP handlers you almost always reach for the streaming form:
`json.NewDecoder(r.Body).Decode(&v)` and `json.NewEncoder(w).Encode(resp)`.

## Struct tag cheat sheet

```go
type User struct {
    Name     string `json:"name"`            // rename the key
    Age      int    `json:"age,omitempty"`   // drop when zero value
    Password string `json:"-"`               // never serialize
    Note     string `json:",omitempty"`      // keep name "Note", omit when empty
    internal string                          // unexported -> always ignored
}
```

## Key habits

- Only **exported** (capitalized) fields are marshalled — unexported fields are
  invisible to `encoding/json`.
- Always pass a **pointer** to `Unmarshal`/`Decode` (`&v`) so Go can fill it in.
- Always check the returned `error` — malformed JSON is common in the wild.
- JSON numbers decode to **`float64`** when the target is `interface{}`; cast to
  `int` when you need one.
- Use `omitempty` to keep payloads small, and `json:"-"` to keep secrets out.
- Reach for the **streaming** encoder/decoder for files, network I/O, and large
  or repeated values.

## Key packages

- **`encoding/json`** — marshal/unmarshal and the streaming `Encoder`/`Decoder`.
- **`encoding/base64`** — encode binary data as ASCII-safe text (`Std`, `URL`,
  and `Raw` variants). Note: Base64 is *transport encoding*, **not** encryption.
- **`io`** — the `Reader`/`Writer` interfaces the streaming APIs are built on.
