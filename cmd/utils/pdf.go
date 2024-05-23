package utils

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

const (
	aliasNumPages = "NP"
	logoPath      = "cmd/data/assets/revolut.png"
	qrPath        = "cmd/data/assets/qr.png"
	watermarkPath = "cmd/data/assets/watermark.png"
)

func CreatePDF(invoice *Invoice, outputName string) {

	// A4 width 219
	// margins 10mm each side
	// usable width 199

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)

	// Add header and footer
	header(pdf, invoice, logoPath, watermarkPath) // add watermark in the header so it will be on every page
	footer(pdf, qrPath)

	// Add a page
	pdf.AddPage()

	// Add client information and balance summaary
	clientInfo(pdf, invoice)
	balanceSummary(pdf, invoice)
	disclaimerBalance(pdf)

	// Add account transactions
	AllTransactions(pdf, invoice)

	// AliasNbPages replaces the placeholder with the total number of pages
	pdf.AliasNbPages(aliasNumPages)

	// Save the PDF to a file
	if err := pdf.OutputFileAndClose(outputName); err != nil {
		panic(err)
	}
}

func header(pdf *gofpdf.Fpdf, invoice *Invoice, logoPath string, watermarkPath string) {
	pdf.SetHeaderFunc(func() {
		// Add logo to the header
		pdf.Image(logoPath, 7, 7, 30, 0, false, "", 0, "")

		// Add "EUR STATEMENT" at the right corner
		pdf.SetFont("Arial", "B", 12)
		pdf.SetXY(170, 10) // Adjust X and Y positions as needed
		pdf.CellFormat(30, 10, "EUR STATEMENT", "", 0, "R", false, 0, "")

		// Add "GENERATED at May 12, 2024" below "EUR STATEMENT"
		pdf.SetXY(170, 15) // Adjust Y position to place the text below "EUR STATEMENT"
		pdf.SetFont("Arial", "I", 6)
		text := fmt.Sprintf("Generated at %v %v, %v", invoice.GeneratedAt.Month(), invoice.GeneratedAt.Day(), invoice.GeneratedAt.Year())
		pdf.CellFormat(30, 10, text, "", 0, "R", false, 0, "")

		// Add Revolut Ltd
		pdf.SetXY(170, 19) // Adjust Y position to place the text below "EUR STATEMENT"
		pdf.SetFont("Arial", "I", 6)
		pdf.CellFormat(30, 10, "Revolut Ltd", "", 1, "R", false, 0, "")

		// add watermark
		watermark(pdf, watermarkPath)
	})
}

func footer(pdf *gofpdf.Fpdf, qrPath string) {

	text1 := "Report lost/stolen card"
	text2 := "+44 20 3322 2222"
	text3 := "Get help in the app"
	text4 := "Scan the QR code"
	text5 := "Revolut Bank UAB is a credit institution licensed in the Republic of Lithuania with company number 304580906 and authorisation code LB002119, and whose registered office is at Konstitucijos ave. 21B, LT-08130,"

	text6 := "Vilnius the Republic of Lithuania. We are licensed and regulated by the Bank of Lithuania and the European Central Bank. The deposits are protected by Lithuanian Deposit Insurance System but some exceptions"

	text7 := "may apply. Please refer to our Deposit Insurance Information document here. More information on deposit insurance of the State Enterprise Deposit and Investment Insurance (VĮ “Indėlių ir investicijų"

	text8 := "draudimas”) is available at www.iidraudimas. lt. If you have any questions, please reach out to us via the in-app chat in the Revolut app."

	pdf.SetFooterFunc(func() {
		pdf.SetY(-35)

		// Add QR code to the footer
		pdf.Image(qrPath, 10, pdf.GetY(), 15, 15, false, "", 0, "")

		pdf.SetFont("Arial", "I", 6)
		pdf.CellFormat(20, 20, "", "", 0, "L", false, 0, "") // empty cell to simulate the image inserted
		pdf.CellFormat(30, 6.5, text1, "", 0, "L", false, 0, "")

		pdf.SetFont("Arial", "I", 4)
		pdf.CellFormat(129, 6.5, text5, "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 6.5, "", "", 1, "L", false, 0, "") //

		pdf.SetY(pdf.GetY() - 4)

		pdf.SetFont("Arial", "I", 6)
		pdf.CellFormat(20, 6.5, "", "", 0, "L", false, 0, "") // empty cell to simulate the image inserted
		pdf.CellFormat(30, 6.5, text2, "", 0, "L", false, 0, "")

		pdf.SetFont("Arial", "I", 4)
		pdf.CellFormat(129, 6.5, text6, "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 6.5, "", "", 1, "L", false, 0, "") //
		pdf.SetY(pdf.GetY() - 4)

		pdf.SetFont("Arial", "I", 6)
		pdf.CellFormat(20, 6.5, "", "", 0, "L", false, 0, "") // empty cell to simulate the image inserted
		pdf.CellFormat(30, 6.5, text3, "", 0, "L", false, 0, "")

		pdf.SetFont("Arial", "I", 4)
		pdf.CellFormat(129, 6.5, text7, "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 6.5, "", "", 1, "L", false, 0, "") //
		pdf.SetY(pdf.GetY() - 4)

		pdf.SetFont("Arial", "I", 6)
		pdf.CellFormat(20, 6.5, "", "", 0, "L", false, 0, "") // empty cell to simulate the image inserted
		pdf.CellFormat(30, 6.5, text4, "", 0, "L", false, 0, "")

		pdf.SetFont("Arial", "I", 4)
		pdf.CellFormat(129, 6.5, text8, "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 6.5, "", "", 1, "L", false, 0, "") //

		pdf.SetFont("Arial", "B", 8)
		pdf.SetY(pdf.GetY() + 2)

		tr := pdf.UnicodeTranslatorFromDescriptor("")
		pdf.CellFormat(95, 10, fmt.Sprintf("%v Revolut Ltd", tr("©")), "", 0, "L", false, 0, "")
		pdf.CellFormat(95, 10, fmt.Sprintf("Page %d of NP", pdf.PageNo()), "", 0, "R", false, 0, "")
	})

}

func clientInfo(pdf *gofpdf.Fpdf, invoice *Invoice) {
	// Set font
	pdf.SetFont("Arial", "B", 20)

	// name of the client
	//CellFormat(w, h, txtStr, borderStr, ln, alignStr, fill bool, link int, linkStr string)
	pdf.SetXY(10, 50)
	pdf.CellFormat(189, 20, invoice.Client.Name, "", 1, "L", false, 0, "")

	pdf.SetY(pdf.GetY() - 4)

	// address of the client
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(100, 9, invoice.Client.Address, "", 0, "L", false, 0, "")
	pdf.CellFormat(30, 9, "", "", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(10, 9, "IBAN", "", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(49, 9, invoice.Client.BankAccount.IBAN, "", 1, "L", false, 0, "")

	// Add a small vertical space between IBAN and BIC
	pdf.SetY(pdf.GetY() - 4)

	pdf.CellFormat(100, 9, invoice.Client.Email, "", 0, "L", false, 0, "")
	pdf.CellFormat(30, 9, "", "", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(10, 9, "BIC", "", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(49, 9, invoice.Client.BankAccount.BIC, "", 1, "L", false, 0, "")
}

func balanceSummary(pdf *gofpdf.Fpdf, invoice *Invoice) {

	// add title
	pdf.SetFont("Arial", "B", 15)
	pdf.SetY(pdf.GetY() + 10)

	pdf.CellFormat(189, 15, "Balance Summary", "", 1, "L", false, 0, "")

	// add table header (Product, Opening Balance, Money Out, Money In, Closing Balance)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)

	pdf.CellFormat(60, 10, "Account", "", 0, "L", true, 0, "")
	pdf.CellFormat(32.25, 10, "Opening Balance", "", 0, "L", true, 0, "")
	pdf.CellFormat(32.25, 10, "Money Out", "", 0, "L", true, 0, "")
	pdf.CellFormat(32.25, 10, "Money In", "", 0, "L", true, 0, "")
	pdf.CellFormat(32.25, 10, "Closing Balance", "", 1, "L", true, 0, "")

	// add line
	pdf.SetFillColor(255, 255, 255)
	pdf.SetFont("Arial", "", 10)
	// line width
	pdf.SetLineWidth(.5)
	pdf.Line(10, pdf.GetY(), 199, pdf.GetY())

	// add data
	data := invoice.Summary
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	var ob, mo, mi, cb float64
	for _, account := range data {

		pdf.CellFormat(60, 10, account.Name, "", 0, "L", false, 0, "")
		pdf.CellFormat(32.25, 10, fmt.Sprintf("%s %.2f", tr("€"), account.OpeningBalance), "", 0, "L", false, 0, "")
		pdf.CellFormat(32.25, 10, fmt.Sprintf("%s %.2f", tr("€"), account.MoneyOut), "", 0, "L", false, 0, "")
		pdf.CellFormat(32.25, 10, fmt.Sprintf("%s %.2f", tr("€"), account.MoneyIn), "", 0, "L", false, 0, "")
		pdf.CellFormat(32.25, 10, fmt.Sprintf("%s %.2f", tr("€"), account.ClosingBalance), "", 1, "L", false, 0, "")

		// add line
		pdf.SetFillColor(255, 255, 255)
		pdf.SetFont("Arial", "", 10)
		// line width
		pdf.SetLineWidth(.2)
		pdf.Line(10, pdf.GetY(), 199, pdf.GetY())

		ob += account.OpeningBalance
		mo += account.MoneyOut
		mi += account.MoneyIn
		cb += account.ClosingBalance
	}

	// add total
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(60, 10, "Total", "", 0, "L", false, 0, "")
	pdf.CellFormat(32.25, 10, fmt.Sprintf("%s %.2f", tr("€"), ob), "", 0, "L", false, 0, "")
	pdf.CellFormat(32.25, 10, fmt.Sprintf("%s %.2f", tr("€"), mo), "", 0, "L", false, 0, "")
	pdf.CellFormat(32.25, 10, fmt.Sprintf("%s %.2f", tr("€"), mi), "", 0, "L", false, 0, "")
	pdf.CellFormat(32.25, 10, fmt.Sprintf("%s %.2f", tr("€"), cb), "", 1, "L", false, 0, "")
}

func disclaimerBalance(pdf gofpdf.Pdf) {
	pdf.SetFont("Arial", "I", 6)
	pdf.SetY(pdf.GetY() + 5)
	pdf.CellFormat(189, 6, "The balance on your statement might differ from the balance shown in your app. The statement balance only reflects completed transactions, while the app shows the balance available for use,", "", 1, "L", false, 0, "")

	pdf.SetY(pdf.GetY() - 3)
	pdf.CellFormat(189, 6, "which accounts for pending transactions.", "", 1, "L", false, 0, "")
}

func AllTransactions(pdf gofpdf.Pdf, invoice *Invoice) {

	// title
	pdf.SetFont("Arial", "B", 15)
	pdf.SetY(pdf.GetY() + 10)

	text := fmt.Sprintf("Account Transactions from %v %v, %v to %v %v, %v",
		invoice.Transactions[0].Date.Month(),
		invoice.Transactions[0].Date.Day(),
		invoice.Transactions[0].Date.Year(),
		invoice.Transactions[len(invoice.Transactions)-1].Date.Month(),
		invoice.Transactions[len(invoice.Transactions)-1].Date.Day(),
		invoice.Transactions[len(invoice.Transactions)-1].Date.Year(),
	)
	pdf.CellFormat(189, 15, text, "", 1, "L", false, 0, "")

	// add table header (Date, Description, Money Out, Money In, Balance)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)

	pdf.CellFormat(32.25, 10, "Date", "", 0, "L", true, 0, "")
	pdf.CellFormat(60, 10, "Description", "", 0, "L", true, 0, "")
	pdf.CellFormat(32.25, 10, "Money Out", "", 0, "L", true, 0, "")
	pdf.CellFormat(32.25, 10, "Money In", "", 0, "L", true, 0, "")
	pdf.CellFormat(32.25, 10, "Balance", "", 1, "L", true, 0, "")

	// add line
	pdf.SetFillColor(255, 255, 255)
	pdf.SetFont("Arial", "", 10)
	// line width
	pdf.SetLineWidth(.5)
	pdf.Line(10, pdf.GetY(), 199, pdf.GetY())

	// add data
	data := invoice.Transactions
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	availableSpace := 297 - pdf.GetY() - 37 // remove footer height(37)

	for _, transaction := range data {

		if availableSpace < 20 {
			pdf.AddPage()
			pdf.SetY(40)
			availableSpace = 297 - 30 - 37 // remove header height(30) and footer height(37)

			// add table header (Date, Description, Money Out, Money In, Balance)
			pdf.SetFont("Arial", "B", 10)
			pdf.SetFillColor(240, 240, 240)

			pdf.CellFormat(32.25, 10, "Date", "", 0, "L", true, 0, "")
			pdf.CellFormat(60, 10, "Description", "", 0, "L", true, 0, "")
			pdf.CellFormat(32.25, 10, "Money Out", "", 0, "L", true, 0, "")
			pdf.CellFormat(32.25, 10, "Money In", "", 0, "L", true, 0, "")
			pdf.CellFormat(32.25, 10, "Balance", "", 1, "L", true, 0, "")

			// add line
			pdf.SetFillColor(255, 255, 255)
			pdf.SetFont("Arial", "", 10)
			// line width
			pdf.SetLineWidth(.5)
			pdf.Line(10, pdf.GetY(), 199, pdf.GetY())
			pdf.SetFont("Arial", "", 10)

			availableSpace -= 10
		}

		pdf.CellFormat(32.25, 10, transaction.Date.Format("Jan 02, 2006"), "", 0, "L", false, 0, "")
		pdf.CellFormat(60, 10, transaction.Description, "", 0, "L", false, 0, "")
		pdf.CellFormat(32.25, 10, fmt.Sprintf("%s %.2f", tr("€"), transaction.MoneyOut), "", 0, "L", false, 0, "")
		pdf.CellFormat(32.25, 10, fmt.Sprintf("%s %.2f", tr("€"), transaction.MoneyIn), "", 0, "L", false, 0, "")
		pdf.CellFormat(32.25, 10, fmt.Sprintf("%s %.2f", tr("€"), transaction.Balance), "", 1, "L", false, 0, "")

		// add line
		pdf.SetFillColor(255, 255, 255)
		pdf.SetFont("Arial", "", 10)
		// line width
		pdf.SetLineWidth(.2)
		pdf.Line(10, pdf.GetY(), 199, pdf.GetY())

		availableSpace -= 10
	}
}

func watermark(pdf *gofpdf.Fpdf, watermarkPath string) {

	// place in the middle of the page
	pdf.TransformBegin()
	pdf.TransformRotate(35, 105, 105)
	pdf.Image(watermarkPath, 40, 135, 90, 0, false, "", 0, "")
	pdf.TransformEnd()
}
