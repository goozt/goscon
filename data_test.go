package goscon

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func Print[T any](d ...T) {
	fmt.Print(d)
}

func TestGetTime(t *testing.T) {
	trans := Transaction{}
	trans.Date, _ = time.Parse(ISTDATETIMEFORMAT, "30/12/2021 23:22:21 +0530")
	got := trans.GetTime()
	want := "30/12/2021 23:22:21 +0530"
	if got != want {
		t.Fatalf("got %v, wanted %v", got, want)
	}
	trans.Date, _ = time.Parse(ISTDATEFORMAT, "30/12/2021 +0530")
	got = trans.GetTime()
	want = "30/12/2021 00:00:00 +0530"
	if got != want {
		t.Fatalf("got %v, wanted %v", got, want)
	}
}

func TestSetTime(t *testing.T) {
	trans := Transaction{}
	trans.SetTime(ISTDATETIMEFORMAT, "30/12/2021 23:22:21 +0530")
	got := trans.GetTime()
	want := "30/12/2021 23:22:21 +0530"
	if got != want {
		t.Fatalf("got %v, wanted %v", got, want)
	}
}

func TestTotalDue(t *testing.T) {
	trans := []Transaction{
		{
			time.Now(),
			"Somethin Purchased",
			34645.45,
			false,
		},
		{
			time.Now(),
			"Somethin Purchased",
			6754.45,
			false,
		},
		{
			time.Now(),
			"Autopay",
			2363.56,
			true,
		},
		{
			time.Now(),
			"Refund",
			7346.45,
			true,
		},
	}
	s := Statement{}
	s.Opening = 10000.00
	s.Transactions = trans
	got := s.TotalDues()
	want := s.Opening
	for _, t := range trans {
		if t.Credited {
			want -= t.Amount
		} else {
			want += t.Amount
		}
	}
	want = math.Round(want)
	if got != want {
		t.Fatalf("expected %.2f, got %.2f", want, got)
	}
}

func TestPurchase(t *testing.T) {
	t1 := Transaction{
		time.Now(),
		"Bag Purchased",
		34645.45,
		false,
	}
	t2 := Transaction{
		time.Now(),
		"Shoe Purchased",
		6754.45,
		false,
	}
	s := Statement{}
	s.Transactions = []Transaction{t1, t2}
	got := math.Round(s.Purchase()*100) / 100
	want := math.Round((t1.Amount+t2.Amount)*100) / 100
	if got != want {
		t.Fatalf("expected %.2f, got %.2f", want, got)
	}
}

func TestPayment(t *testing.T) {
	t1 := Transaction{
		time.Now(),
		"Autopay",
		2363.56,
		true,
	}
	t2 := Transaction{
		time.Now(),
		"Refund",
		7346.45,
		true,
	}
	s := Statement{}
	s.Transactions = []Transaction{t1, t2}
	got := math.Round(s.Payment()*100) / 100
	want := math.Round((t1.Amount+t2.Amount)*100) / 100
	if got != want {
		t.Fatalf("expected %.2f, got %.2f", want, got)
	}
}

func TestFormatedStatement(t *testing.T) {
	trans := []Transaction{
		{
			time.Now(),
			"Bag Purchased",
			34645.45,
			false,
		},
		{
			time.Now(),
			"Shoe Purchased",
			6754.45,
			false,
		},
		{
			time.Now(),
			"Autopay",
			2363.56,
			true,
		},
		{
			time.Now(),
			"Refund",
			7346.45,
			true,
		},
	}
	s := Statement{}
	s.Opening = 10000.00
	s.Transactions = trans
	got := s.FormatedStatement()
	want := `---------------------------------------------------------------------------
Date                       |     Amount | Description
---------------------------------------------------------------------------
` + trans[0].GetTime() + `  |   34645.45 | Bag Purchased
` + trans[1].GetTime() + `  |    6754.45 | Shoe Purchased
` + trans[2].GetTime() + `  |    2363.56 | Autopay
` + trans[3].GetTime() + `  |    7346.45 | Refund
---------------------------------------------------------------------------
Total Dues: 41690.00
---------------------------------------------------------------------------
`
	if got != want {
		t.Fatalf("expected %s, got %s", want, got)
	}
}
