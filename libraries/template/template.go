package main

import (
	"os"
	"strings"
	"text/template"
)

type Product struct {
	Name     string
	Price    float64
	InStock  bool
	Category string
}

type Store struct {
	Name     string
	Products []Product
}

func main() {
	store := Store{
		Name: "Go Store",
		Products: []Product{
			{Name: "T-Shirt", Price: 499.99, InStock: true, Category: "clothing"},
			{Name: "Mug", Price: 199.50, InStock: true, Category: "kitchen"},
			{Name: "Sticker Pack", Price: 99.00, InStock: false, Category: "accessories"},
		},
	}

	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
	}

	const storeTemplate = `
Store: {{.Name}}

Products:
{{range .Products}}
- {{.Name}} | Rs. {{printf "%.2f" .Price}} | {{upper .Category}}
  {{if .InStock}}Available{{else}}Out of stock{{end}}
{{end}}
`

	tmpl := template.Must(template.New("store").Funcs(funcMap).Parse(storeTemplate))
	tmpl.Execute(os.Stdout, store)
}
