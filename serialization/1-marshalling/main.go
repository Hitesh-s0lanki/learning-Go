package main

import (
	"encoding/json"
	"fmt"
	"log"
)

/*
MARSHALLING (Go value -> JSON)
==============================

"Marshalling" means turning a Go value into a JSON []byte.

  json.Marshal(v)              -> compact JSON on one line
  json.MarshalIndent(v, "", "  ") -> pretty-printed, human-readable JSON

Only EXPORTED (capitalized) struct fields are marshalled. Control the output
with struct tags:

  `json:"name"`          -> use "name" as the JSON key
  `json:"phone,omitempty"` -> drop the key entirely if the field is a zero value
  `json:"-"`             -> never marshal this field (e.g. passwords)
  `json:",omitempty"`    -> keep the field name, but omit when empty

Basic Go -> JSON type mapping:
  string        -> string
  int/float     -> number
  bool          -> true/false
  slice/array   -> array
  map/struct    -> object
  nil pointer   -> null
*/

type user struct {
	Name     string   `json:"name"`
	Age      int      `json:"age"`
	Phone    string   `json:"phone,omitempty"` // dropped when empty
	Password string   `json:"-"`               // never leaks into JSON
	IsActive bool     `json:"active"`
	Roles    []string `json:"roles"`
}

func main() {
	u := user{
		Name:     "John Doe",
		Age:      42,
		Password: "s3cr3t", // present in memory, absent from JSON
		IsActive: true,
		Roles:    []string{"admin", "editor"},
		// Phone left empty on purpose to show omitempty.
	}

	// --- 1. Compact marshalling (one line, no spaces) ---
	compact, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("compact:")
	fmt.Println(string(compact))

	// --- 2. Pretty (indented) marshalling ---
	pretty, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nindented:")
	fmt.Println(string(pretty))

	// --- 3. Marshalling a slice of structs -> JSON array ---
	users := []user{
		{Name: "Alice", Age: 30, IsActive: true, Roles: []string{"user"}},
		{Name: "Bob", Age: 25, Roles: []string{}},
	}
	arr, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nslice of users:")
	fmt.Println(string(arr))

	// --- 4. Marshalling a map -> JSON object (keys sorted automatically) ---
	scores := map[string]int{"go": 95, "rust": 90, "python": 88}
	m, err := json.Marshal(scores)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nmap:")
	fmt.Println(string(m))
}
