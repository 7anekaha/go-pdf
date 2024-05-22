package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	utils "github.com/7anekaha/go-pdf/cmd/utils"
)

func main() {

	// read transactions from a file or database
	file, err := os.Open("cmd/data/transactions.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// decode the json file
	var data utils.Data
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		panic(err)
	}

	// Create an invoice
	invoice := utils.NewInvoice(data.Client, data.Summary, data.Transactions)

	// create a PDF
	outputName := fmt.Sprintf("invoice_%s.pdf", invoice.GeneratedAt.Format("2006-01-02"))
	utils.CreatePDF(invoice, outputName)

	log.Println("PDF created successfully:", outputName)
}
