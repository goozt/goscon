package goscon

import (
	"testing"
	"time"
)

func TestToStrings(t *testing.T) {
	trans := []Transaction{
		{
			time.Now(),
			"Bag Purchased",
			34645.45,
			false,
		},
		{
			time.Now(),
			"Autopay",
			2363.56,
			true,
		},
	}
	s := Statement{}
	s.Opening = 10000.00
	s.Transactions = trans
	got := s.ToStrings()
	want := [][]string{{"DATE", "DESCRIPTION", "AMOUNT", "CREDITED"}}
	want = append(want, []string{trans[0].GetTime(), "Bag Purchased", "34645.45", "false"})
	want = append(want, []string{trans[1].GetTime(), "Autopay", "2363.56", "true"})
	for i, line := range got {
		for j, char := range line {
			if char != want[i][j] {
				t.Fatalf("expected %v, got %v", want, got)
			}
		}
	}
}
