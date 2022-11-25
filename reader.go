package goscon

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dslipak/pdf"
)

func isLine(content pdf.TextHorizontal) bool {
	if len(content) < 4 {
		return false
	}
	expDate := `([^.]\d+( )+?)?\d\d/\d\d/\d{4}`
	expTime := `(( )*?\d\d:\d\d:\d\d)?`
	startLine := regexp.MustCompile(expDate + expTime)
	NotHead := regexp.MustCompile(`(\d\d,)*\d{3}\.\d\d`)
	isNotHead := !(NotHead.MatchString(content[1].S))
	return startLine.MatchString(content[0].S) && isNotHead
}

func isCredit(content pdf.TextHorizontal) bool {
	if len(content) < 4 {
		return false
	}
	return content[len(content)-1].S == "Cr"
}

func isOpeningBalance(content pdf.TextHorizontal) bool {

	if content.Len() == 6 {
		opening, err := parseFloat(content[1].S)
		if err != nil {
			return false
		}
		credit, err := parseFloat(content[2].S)
		if err != nil {
			return false
		}
		debit, err := parseFloat(content[3].S)
		if err != nil {
			return false
		}
		charges, err := parseFloat(content[4].S)
		if err != nil {
			return false
		}
		total, err := parseFloat(content[5].S)
		if err != nil {
			return false
		}
		return math.Round(total) == math.Round(opening-credit+debit+charges)
	}
	return false
}

func openPdfFile(path string) (*os.File, *pdf.Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, nil, err
	}
	r, err := pdf.NewReader(f, fi.Size())
	if err != nil {
		return nil, nil, err
	}
	return f, r, nil
}

func Read(file string) (Statement, error) {
	f, r, err := openPdfFile(file)
	defer func() { f.Close() }()
	if err != nil {
		return Statement{}, err
	}

	var openingBalance float64
	statement := Statement{}
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			fmt.Printf("%#v\n", row.Content)
			if isOpeningBalance(row.Content) {
				o, err := parseFloat(row.Content[1].S)
				if err == nil {
					openingBalance = o
				}
				continue
			}

			if isLine(row.Content) {
				amt := strings.Replace(cleanString(row.Content[len(row.Content)-2].S), ",", "", -1)
				amount, err := strconv.ParseFloat(amt, 64)
				if err != nil {
					return Statement{}, err
				}

				date, err := time.Parse(ISTDATETIMEFORMAT, cleanString(row.Content[0].S)+" +0530")
				if err != nil {
					date, err = time.Parse(ISTDATEFORMAT, cleanString(row.Content[0].S)+" +0530")
				}
				if err != nil {
					return Statement{}, err
				}

				statement.Transactions = append(statement.Transactions, Transaction{
					Date:        date,
					Description: cleanString(row.Content[1].S),
					Amount:      amount,
					Credited:    isCredit(row.Content),
				})
			}
		}
	}

	exp := regexp.MustCompile(`[A-Z][a-z]{2,3} \d{4}`)
	statement.MonthYear = exp.FindString(file)
	statement.Opening = openingBalance

	return statement, nil
}
