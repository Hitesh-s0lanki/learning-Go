package main

import (
	"fmt"
	"regexp"
)

func main() {
	text := "Contact Hemant at hemant@example.com or support@example.org. Order ID: ORD-2026-7845"

	fmt.Println("Text:", text)

	fmt.Println("----- Basic Match -----")

	hasEmail, err := regexp.MatchString(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`, text)
	if err != nil {
		fmt.Println("Regex error:", err)
		return
	}

	fmt.Println("Has email:", hasEmail)

	fmt.Println("----- Find Values -----")

	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	emails := emailRegex.FindAllString(text, -1)
	fmt.Println("Emails:", emails)

	fmt.Println("----- Replace Values -----")

	hiddenEmails := emailRegex.ReplaceAllString(text, "[hidden-email]")
	fmt.Println("Hidden emails:", hiddenEmails)

	fmt.Println("----- Split String -----")

	csvText := "apple, banana; mango|orange"
	separatorRegex := regexp.MustCompile(`[,;|]\s*`)
	fruits := separatorRegex.Split(csvText, -1)
	fmt.Println("Fruits:", fruits)

	fmt.Println("----- Capturing Groups -----")

	orderRegex := regexp.MustCompile(`ORD-(\d{4})-(\d{4})`)
	matches := orderRegex.FindStringSubmatch(text)

	if len(matches) > 0 {
		fmt.Println("Full order ID:", matches[0])
		fmt.Println("Year:", matches[1])
		fmt.Println("Number:", matches[2])
	}

	fmt.Println("----- Validate Phone Number -----")

	phoneRegex := regexp.MustCompile(`^\d{10}$`)
	fmt.Println("9876543210 valid:", phoneRegex.MatchString("9876543210"))
	fmt.Println("98765 valid:", phoneRegex.MatchString("98765"))
}
