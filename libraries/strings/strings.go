package main

import (
	"fmt"
	"strings"
)

func main() {
	message := "  Go is simple, Go is fast  "

	fmt.Println("Original:", message)

	// Remove spaces from the beginning and end.
	trimmedMessage := strings.TrimSpace(message)
	fmt.Println("TrimSpace:", trimmedMessage)

	// Convert to lowercase and uppercase.
	fmt.Println("ToLower:", strings.ToLower(trimmedMessage))
	fmt.Println("ToUpper:", strings.ToUpper(trimmedMessage))

	// Check if a string contains another string.
	fmt.Println("Contains Go:", strings.Contains(trimmedMessage, "Go"))

	// Count how many times a word appears.
	fmt.Println("Count Go:", strings.Count(trimmedMessage, "Go"))

	// Replace text.
	replacedMessage := strings.ReplaceAll(trimmedMessage, "Go", "Golang")
	fmt.Println("ReplaceAll:", replacedMessage)

	// Split string into a slice.
	csvText := "apple,banana,mango,orange"
	fruits := strings.Split(csvText, ",")
	fmt.Println("Split:", fruits)

	// Join slice values into a string.
	joinedFruits := strings.Join(fruits, " | ")
	fmt.Println("Join:", joinedFruits)

	// Check prefix and suffix.
	fileName := "report.pdf"
	fmt.Println("HasPrefix report:", strings.HasPrefix(fileName, "report"))
	fmt.Println("HasSuffix .pdf:", strings.HasSuffix(fileName, ".pdf"))

	// Find index of a word.
	fmt.Println("Index of simple:", strings.Index(trimmedMessage, "simple"))

	fmt.Println("----- Strings Builder -----")

	var builder strings.Builder

	builder.WriteString("Go")
	builder.WriteByte(' ')
	builder.WriteString("is")
	builder.WriteByte(' ')
	builder.WriteString("powerful")

	finalMessage := builder.String()
	fmt.Println("Builder result:", finalMessage)

	words := []string{"Learn", "Go", "step", "by", "step"}
	var sentenceBuilder strings.Builder

	for index, word := range words {
		if index > 0 {
			sentenceBuilder.WriteByte(' ')
		}
		sentenceBuilder.WriteString(word)
	}

	fmt.Println("Sentence:", sentenceBuilder.String())
}
