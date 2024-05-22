package utils

import (
	"time"
)

type BankAccount struct {
	IBAN string `json:"iban"`
	BIC  string `json:"bic"`
}

type Client struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	BankAccount BankAccount `json:"bank_account"`
}

type Account struct {
	Name           string `json:"account"`
	OpeningBalance float64 `json:"opening_balance"`
	MoneyOut       float64 `json:"money_out"`
	MoneyIn        float64 `json:"money_in"`
	ClosingBalance float64 `json:"closing_balance"`
}

type Transaction struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	MoneyOut    float64   `json:"money_out,omitempty"`
	MoneyIn     float64   `json:"money_in,omitempty"`
	Balance     float64   `json:"balance"`
}

type Data struct {
	Client Client `json:"client"`
	Summary []Account `json:"summary"`
	Transactions []Transaction `json:"transactions"`
}


type Invoice struct {
	GeneratedAt    time.Time
	Client         Client
	Summary []Account
	Transactions   []Transaction //this Must be sorted by Date
}

func NewInvoice(client Client, summary []Account, transactions []Transaction) *Invoice {
	return &Invoice{
		GeneratedAt: time.Now(),
		Client:      client,
		Summary:     summary,
		Transactions: transactions,
	}
}