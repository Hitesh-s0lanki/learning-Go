# String Formatting in Go

Go uses the `fmt` package for string formatting.

The most common formatting functions are:

```go
fmt.Printf("Name: %s", name)
fmt.Sprintf("Name: %s", name)
fmt.Fprintf(writer, "Name: %s", name)
```

## Printf vs Sprintf

`Printf` prints formatted output directly.

```go
name := "Hemant"
age := 22

fmt.Printf("Name: %s, Age: %d\n", name, age)
```

Output:

```text
Name: Hemant, Age: 22
```

`Sprintf` returns the formatted string instead of printing it.

```go
message := fmt.Sprintf("Name: %s, Age: %d", name, age)
fmt.Println(message)
```

## Common Formatting Verbs

| Verb | Meaning | Example |
| --- | --- | --- |
| `%s` | string | `"Go"` |
| `%d` | integer | `25` |
| `%f` | float | `99.500000` |
| `%.2f` | float with 2 decimal places | `99.50` |
| `%t` | boolean | `true` |
| `%v` | default value format | any value |
| `%+v` | struct with field names | `{Name:Hemant Age:22}` |
| `%#v` | Go-syntax representation | `main.User{Name:"Hemant"}` |
| `%T` | type of value | `string`, `int` |
| `%p` | pointer address | `0x...` |

## String Formatting

```go
language := "Go"

fmt.Printf("Language: %s\n", language)
```

Output:

```text
Language: Go
```

## Integer Formatting

```go
count := 50

fmt.Printf("Count: %d\n", count)
```

Output:

```text
Count: 50
```

## Float Formatting

```go
price := 99.9876

fmt.Printf("Price: %f\n", price)
fmt.Printf("Price: %.2f\n", price)
```

Output:

```text
Price: 99.987600
Price: 99.99
```

## Boolean Formatting

```go
isActive := true

fmt.Printf("Active: %t\n", isActive)
```

Output:

```text
Active: true
```

## Struct Formatting

```go
type User struct {
	Name string
	Age  int
}

user := User{Name: "Hemant", Age: 22}

fmt.Printf("Default: %v\n", user)
fmt.Printf("With fields: %+v\n", user)
fmt.Printf("Go syntax: %#v\n", user)
```

Example output:

```text
Default: {Hemant 22}
With fields: {Name:Hemant Age:22}
Go syntax: main.User{Name:"Hemant", Age:22}
```

## Width and Alignment

You can control spacing using width.

```go
fmt.Printf("|%10s|\n", "Go")
fmt.Printf("|%-10s|\n", "Go")
```

Output:

```text
|        Go|
|Go        |
```

`%10s` means right-align within 10 spaces.

`%-10s` means left-align within 10 spaces.

## Padding Numbers

```go
number := 7

fmt.Printf("%03d\n", number)
```

Output:

```text
007
```

## Escaping Percent Sign

Use `%%` to print a percent sign.

```go
discount := 10

fmt.Printf("Discount: %d%%\n", discount)
```

Output:

```text
Discount: 10%
```

## Complete Example

```go
package main

import "fmt"

type Product struct {
	Name  string
	Price float64
	Stock int
}

func main() {
	product := Product{Name: "T-Shirt", Price: 499.99, Stock: 25}

	fmt.Printf("Product: %s\n", product.Name)
	fmt.Printf("Price: %.2f\n", product.Price)
	fmt.Printf("Stock: %03d\n", product.Stock)
	fmt.Printf("Details: %+v\n", product)

	message := fmt.Sprintf("%s costs %.2f", product.Name, product.Price)
	fmt.Println(message)
}
```
