package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

/*
BASE64 ENCODING (encoding/base64)
=================================

Base64 turns arbitrary BINARY data into an ASCII-safe text string so it can
travel through systems that only handle text (URLs, JSON, email, HTTP headers).
It is NOT encryption — anyone can decode it. It's just a safe transport format.

Go gives you a few ready-made encodings:

  base64.StdEncoding     -> standard alphabet, '+' and '/', padded with '='
  base64.URLEncoding     -> URL/filename safe: '-' and '_' instead of '+' '/'
  base64.RawStdEncoding  -> standard alphabet, NO '=' padding
  base64.RawURLEncoding  -> URL-safe, no padding

Core calls:
  enc.EncodeToString(b) -> string
  enc.DecodeString(s)   -> ([]byte, error)
*/

func main() {
	data := []byte("Welcome to the wonderful world of Go! <3")

	// --- 1. Standard encoding (padded with '=') ---
	std := base64.StdEncoding.EncodeToString(data)
	fmt.Println("StdEncoding:    ", std)

	// --- 2. URL-safe encoding (safe inside query strings / filenames) ---
	url := base64.URLEncoding.EncodeToString(data)
	fmt.Println("URLEncoding:    ", url)

	// --- 3. Raw (unpadded) encoding ---
	raw := base64.RawStdEncoding.EncodeToString(data)
	fmt.Println("RawStdEncoding: ", raw)

	// --- 4. Decode back and verify the round-trip ---
	decoded, err := base64.StdEncoding.DecodeString(std)
	if err != nil {
		log.Fatal(err)
	}
	if string(decoded) != string(data) {
		log.Fatal("round-trip mismatch!")
	}
	fmt.Printf("\ndecoded round-trip: %q (matches original: %v)\n", decoded, string(decoded) == string(data))

	// --- 5. Base64 also works on raw binary bytes, not just text ---
	rawBytes := []byte{0xDE, 0xAD, 0xBE, 0xEF, 0xCA, 0xFE}
	encoded := base64.StdEncoding.EncodeToString(rawBytes)
	fmt.Printf("\nraw bytes %v -> %q\n", rawBytes, encoded)

	back, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("decoded bytes: % X\n", back)
}
