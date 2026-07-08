package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func main() {
	text := "Go नमस्ते 世界"

	fmt.Println("Text:", text)

	// len gives the number of bytes, not characters.
	fmt.Println("Byte length:", len(text))

	// utf8.RuneCountInString gives the number of runes.
	fmt.Println("Rune count:", utf8.RuneCountInString(text))

	fmt.Println("----- Loop By Bytes -----")

	for i := 0; i < len(text); i++ {
		fmt.Printf("Byte index: %d, Byte value: %d\n", i, text[i])
	}

	fmt.Println("----- Loop By Runes -----")

	for index, char := range text {
		fmt.Printf("Byte index: %d, Rune: %c, Unicode: %U\n", index, char, char)
	}

	fmt.Println("----- Unicode Checks -----")

	for _, char := range text {
		switch {
		case unicode.IsLetter(char):
			fmt.Printf("%c is a letter\n", char)
		case unicode.IsDigit(char):
			fmt.Printf("%c is a digit\n", char)
		case unicode.IsSpace(char):
			fmt.Println("space is a space")
		default:
			fmt.Printf("%c is another symbol\n", char)
		}
	}

	fmt.Println("----- Convert Rune Case -----")

	letter := 'g'
	fmt.Println("Original:", string(letter))
	fmt.Println("Uppercase:", string(unicode.ToUpper(letter)))

	upperLetter := 'H'
	fmt.Println("Original:", string(upperLetter))
	fmt.Println("Lowercase:", string(unicode.ToLower(upperLetter)))
}
