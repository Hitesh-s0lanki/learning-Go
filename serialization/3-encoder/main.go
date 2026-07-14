package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

/*
STREAMING ENCODER (json.Encoder)
================================

json.Marshal builds the whole JSON in memory and hands you a []byte.
A json.Encoder instead WRITES JSON straight to an io.Writer (a file, a network
connection, os.Stdout, a buffer...). This is the right tool for:

  - large data you don't want to hold entirely in memory
  - HTTP handlers (json.NewEncoder(w).Encode(v))
  - writing many JSON values back-to-back (newline-delimited JSON)

  enc := json.NewEncoder(w)
  enc.SetIndent("", "  ")   // optional pretty printing
  enc.SetEscapeHTML(false)  // keep <, >, & literal instead of < etc.
  enc.Encode(v)             // writes v + a trailing newline

Each Encode call appends a newline, which makes it perfect for streaming
one-object-per-line logs (NDJSON).
*/

type user struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Homepage string `json:"homepage"`
	IsActive bool   `json:"active"`
}

func main() {
	users := []user{
		{Name: "John Smith", Age: 45, Homepage: "https://a.example.com?x=1&y=2", IsActive: true},
		{Name: "Jane Roe", Age: 30, Homepage: "https://b.example.com", IsActive: false},
	}

	// --- 1. Encode straight to stdout ---
	fmt.Println("encode to stdout:")
	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(users[0]); err != nil { // note: trailing newline is added
		log.Fatal(err)
	}

	// --- 2. Pretty-printed encoding + keep HTML characters literal ---
	fmt.Println("indented, no HTML escaping:")
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false) // '&' stays '&' instead of becoming &
	if err := enc.Encode(users[0]); err != nil {
		log.Fatal(err)
	}

	// --- 3. Stream many values as newline-delimited JSON (NDJSON) into a buffer ---
	var buf strings.Builder
	stream := json.NewEncoder(&buf)
	for _, u := range users {
		if err := stream.Encode(u); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Print("\nNDJSON stream (one object per line):\n")
	fmt.Print(buf.String())
}
