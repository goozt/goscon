package goscon

import (
	"io"
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
		opening, err := ParseFloat(content[1].S)
		if err != nil {
			return false
		}
		credit, err := ParseFloat(content[2].S)
		if err != nil {
			return false
		}
		debit, err := ParseFloat(content[3].S)
		if err != nil {
			return false
		}
		charges, err := ParseFloat(content[4].S)
		if err != nil {
			return false
		}
		total, err := ParseFloat(content[5].S)
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
	if err != nil {
		return Statement{}, err
	}
	defer f.Close()

	statement, err := Parse(r)
	if err != nil {
		return Statement{}, err
	}

	if statement.MonthYear == "" {
		exp := regexp.MustCompile(`[A-Z][a-z]{2,3} \d{4}`)
		statement.MonthYear = exp.FindString(file)
	}

	return statement, nil
}

func ReadFrom(r io.ReaderAt, size int64) (Statement, error) {
	pdfReader, err := pdf.NewReader(r, size)
	if err != nil {
		return Statement{}, err
	}
	return Parse(pdfReader)
}

func Parse(r *pdf.Reader) (Statement, error) {
	var openingBalance float64
	statement := Statement{}
	totalPage := r.NumPage()
	monthYearExp := regexp.MustCompile(`[A-Z][a-z]{2,3} \d{4}`)

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			if statement.MonthYear == "" {
				for _, content := range row.Content {
					if monthYearExp.MatchString(content.S) {
						statement.MonthYear = monthYearExp.FindString(content.S)
						break
					}
				}
			}
			if isOpeningBalance(row.Content) {
				o, err := ParseFloat(row.Content[1].S)
				if err == nil {
					openingBalance = o
				}
				continue
			}

			if isLine(row.Content) {
				amt := strings.Replace(CleanString(row.Content[len(row.Content)-2].S), ",", "", -1)
				amount, err := strconv.ParseFloat(amt, 64)
				if err != nil {
					return Statement{}, err
				}

				date, err := time.Parse(ISTDATETIMEFORMAT, CleanString(row.Content[0].S)+" +0530")
				if err != nil {
					date, err = time.Parse(ISTDATEFORMAT, CleanString(row.Content[0].S)+" +0530")
				}
				if err != nil {
					return Statement{}, err
				}

				statement.Transactions = append(statement.Transactions, Transaction{
					Date:        date,
					Description: CleanString(row.Content[1].S),
					Amount:      amount,
					Credited:    isCredit(row.Content),
				})
			}
		}
	}

	statement.Opening = openingBalance

	return statement, nil
}
