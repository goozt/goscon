package main

import (
	"fmt"
	"log"

	"github.com/goozt/goscon"
)

func main() {
	// Example of using the library to parse a PDF file
	file := "path/to/statement.pdf"
	statement, err := goscon.Read(file)
	if err != nil {
		// Handle error (in this sample we just print it as the file likely doesn't exist)
		fmt.Printf("Error reading statement: %v\n", err)
		return
	}

	// Calculate some stats
	fmt.Printf("Statement Month: %s\n", statement.MonthYear)
	fmt.Printf("Opening Balance: %.2f\n", statement.Opening)
	fmt.Printf("Total Purchases: %.2f\n", statement.Purchase())
	fmt.Printf("Total Payments: %.2f\n", statement.Payment())
	fmt.Printf("Total Dues: %.2f\n", statement.TotalDues())

	// Convert to CSV
	err = statement.WriteCSV("statement.csv")
	if err != nil {
		log.Fatalf("Error writing CSV: %v", err)
	}

	// Convert to JSON
	err = statement.WriteJSON("statement.json")
	if err != nil {
		log.Fatalf("Error writing JSON: %v", err)
	}

	fmt.Println("Successfully converted statement to CSV and JSON.")
}
