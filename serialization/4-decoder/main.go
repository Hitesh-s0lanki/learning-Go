package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

/*
STREAMING DECODER (json.Decoder)
================================

json.Unmarshal needs the whole JSON in memory as a []byte. A json.Decoder
READS JSON from an io.Reader (a file, a request body, a socket...) and can:

  - decode one value at a time from a continuous stream
  - reject unknown fields with DisallowUnknownFields (great for strict APIs)
  - stop cleanly at io.EOF when the stream is exhausted

  dec := json.NewDecoder(r)
  dec.DisallowUnknownFields()  // error on keys with no matching struct field
  for {
      var v T
      if err := dec.Decode(&v); err == io.EOF { break }
      ...
  }

In an HTTP handler you'd typically write:
  json.NewDecoder(r.Body).Decode(&v)
*/

type user struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	IsActive bool   `json:"active"`
}

// Three JSON objects back-to-back (whitespace/newlines between them are fine).
var stream = `
{"name":"John Smith","age":42,"active":true}
{"name":"Ada Lovelace","age":36,"active":false}
{"name":"Alan Turing","age":41,"active":true}
`

func main() {
	// --- 1. Decode a continuous stream of values until EOF ---
	dec := json.NewDecoder(strings.NewReader(stream))
	fmt.Println("decoding stream:")
	for {
		var u user
		err := dec.Decode(&u)
		if errors.Is(err, io.EOF) {
			break // no more values
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %+v\n", u)
	}

	// --- 2. Strict decoding: reject unknown fields ---
	strictInput := `{"name":"Grace","age":85,"active":true,"nickname":"Amazing Grace"}`
	strict := json.NewDecoder(strings.NewReader(strictInput))
	strict.DisallowUnknownFields()

	var u user
	if err := strict.Decode(&u); err != nil {
		fmt.Println("\nstrict decode rejected input (as expected):")
		fmt.Println("  ", err)
	} else {
		fmt.Println("\nstrict decode ok:", u)
	}

	// --- 3. Same input decodes fine WITHOUT the strict flag ---
	lenient := json.NewDecoder(strings.NewReader(strictInput))
	if err := lenient.Decode(&u); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nlenient decode ignored the extra field: %+v\n", u)
}
