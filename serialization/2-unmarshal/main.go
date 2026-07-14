package main

import (
	"encoding/json"
	"fmt"
	"log"
)

/*
UNMARSHALLING (JSON -> Go value)
================================

"Unmarshalling" parses a JSON []byte into a Go value. You pass a POINTER so
Go can fill it in:

  json.Unmarshal(data, &v)

Matching is case-insensitive and driven by struct tags. JSON keys with no
matching field are silently ignored; struct fields with no matching key keep
their zero value.

When you don't know the shape ahead of time, unmarshal into:
  map[string]interface{}  -> a JSON object
  []interface{}           -> a JSON array
  interface{}             -> anything

Remember how JSON types arrive in an interface{}:
  number -> float64   string -> string   bool -> bool
  object -> map[string]interface{}   array -> []interface{}   null -> nil
*/

type user struct {
	Name     string  `json:"name"`
	Age      int     `json:"age"`
	Phone    string  `json:"phone"`
	IsActive bool    `json:"active"`
	Profile  profile `json:"profile"` // nested struct
}

type profile struct {
	URL string `json:"url"`
}

var payload = `{
  "name": "Jane",
  "age": 20,
  "phone": "123-456-789",
  "active": true,
  "profile": { "url": "https://jane.example.com" },
  "extra": "ignored because there is no matching field"
}`

func main() {
	// --- 1. Unmarshal into a known struct ---
	var u user
	if err := json.Unmarshal([]byte(payload), &u); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("struct: %+v\n", u)
	fmt.Println("nested URL:", u.Profile.URL)

	// --- 2. Unmarshal into a map when the shape is unknown ---
	var generic map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &generic); err != nil {
		log.Fatal(err)
	}
	// Note: numbers come back as float64, so cast when you need an int.
	age := generic["age"].(float64)
	fmt.Printf("\ngeneric age (float64): %v -> int: %d\n", age, int(age))

	// --- 3. Unmarshal a JSON array into a slice ---
	list := `[{"name":"Ada","age":36},{"name":"Linus","age":54}]`
	var users []user
	if err := json.Unmarshal([]byte(list), &users); err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nparsed array:")
	for _, x := range users {
		fmt.Printf("  %s is %d\n", x.Name, x.Age)
	}

	// --- 4. Round-trip: parse, then re-marshal indented ---
	pretty, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nre-marshalled:")
	fmt.Println(string(pretty))
}
