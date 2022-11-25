package goscon

import (
	"fmt"
	"math"
	"time"
)

const (
	ISTDATETIMEFORMAT = "02/01/2006 15:04:05 -0700"
	ISTDATEFORMAT     = "02/01/2006 -0700"
)

type Transaction struct {
	Date        time.Time
	Description string
	Amount      float64
	Credited    bool
}

func (t *Transaction) SetTime(format string, date string) error {
	var err error
	t.Date, err = time.Parse(format, date)
	return err
}

func (t *Transaction) GetTime() string {
	return t.Date.Format(ISTDATETIMEFORMAT)
}

type Statement struct {
	Transactions []Transaction
	Opening      float64
	MonthYear    string
}

func (s Statement) TotalDues() float64 {
	var total float64 = s.Opening
	for _, t := range s.Transactions {
		if t.Credited {
			total -= t.Amount
		} else {
			total += t.Amount
		}
	}
	return math.Round(total)
}

func (s Statement) Payment() float64 {
	var total float64
	for _, t := range s.Transactions {
		if t.Credited {
			total += t.Amount
		}
	}
	return math.Round(total*100) / 100
}

func (s Statement) Purchase() float64 {
	var total float64
	for _, t := range s.Transactions {
		if !t.Credited {
			total += t.Amount
		}
	}
	return math.Round(total*100) / 100
}

func (s Statement) FormatedStatement() string {
	divider := func() string {
		var out string
		for i := 0; i < 75; i++ {
			out += "-"
		}
		out += "\n"
		return out
	}
	output := divider()
	output += fmt.Sprintf("%-26s | %10s | %v\n", "Date", "Amount", "Description")
	output += divider()
	for _, t := range s.Transactions {
		output += fmt.Sprintf("%-26s | %10.2f | %v\n", t.GetTime(), t.Amount, t.Description)
	}
	output += divider()
	output += fmt.Sprintf("Total Dues: %.2f\n", s.TotalDues())
	output += divider()
	return output
}
