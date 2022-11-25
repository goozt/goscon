package goscon

import (
	"testing"

	"github.com/dslipak/pdf"
)

func TestIsLine(t *testing.T) {
	var line pdf.TextHorizontal = []pdf.Text{
		{X: 13, Y: -494, S: "29/12/2022 02:14:05"},
		{X: 112, Y: -494, S: "Bag Purchase           COCHIN"},
		{X: 542, Y: -494, S: "2000.00"},
		{X: 562.5, Y: -494, S: " "},
	}
	if !isLine(line) {
		t.Fatalf("got %v, wanted %v", true, false)
	}
}
func TestIsNotLine(t *testing.T) {
	var line pdf.TextHorizontal = []pdf.Text{
		{X: 13, Y: -494, S: "29/12/2022 02:14:05"},
		{X: 562.5, Y: -494, S: " "},
	}
	if isLine(line) {
		t.Fatalf("got %v, wanted %v", false, true)
	}
	var line2 pdf.TextHorizontal = []pdf.Text{
		{X: 13, Y: -494, S: "29/12/2022 02:14:05"},
		{X: 112, Y: -494, S: "84,152.64"},
		{X: 542, Y: -494, S: "Transaction"},
		{X: 562.5, Y: -494, S: " "},
	}
	if isLine(line2) {
		t.Fatalf("got %v, wanted %v", false, true)
	}
}

func TestIsCredit(t *testing.T) {
	var line pdf.TextHorizontal = []pdf.Text{
		{X: 13, Y: -494, S: "29/12/2022 02:14:05"},
		{X: 112, Y: -494, S: "Bag Purchase           COCHIN"},
		{X: 542, Y: -494, S: "2000.00"},
		{X: 562.5, Y: -494, S: "Cr"},
	}
	if !isCredit(line) {
		t.Fatalf("got %v, wanted %v", true, false)
	}
}
func TestIsNotCredit(t *testing.T) {
	var line pdf.TextHorizontal = []pdf.Text{
		{X: 13, Y: -494, S: "29/12/2022 02:14:05"},
		{X: 112, Y: -494, S: "Bag Purchase           COCHIN"},
		{X: 542, Y: -494, S: "2000.00"},
		{X: 562.5, Y: -494, S: "  "},
	}
	if isCredit(line) {
		t.Fatalf("got %v, wanted %v", false, true)
	}
	var line2 pdf.TextHorizontal = []pdf.Text{
		{X: 13, Y: -494, S: "29/12/2022 02:14:05"},
		{X: 562.5, Y: -494, S: " "},
	}
	if isCredit(line2) {
		t.Fatalf("got %v, wanted %v", false, true)
	}
}

func TestIsOpeningBalance(t *testing.T) {
	var line pdf.TextHorizontal = []pdf.Text{
		{X: 16, Y: -236, S: "(including fresh purchases, if any) on an average daily"},
		{X: 264, Y: -236.62, S: "41,954.43"},
		{X: 332.12, Y: -236.62, S: "41,954.00"},
		{X: 402.12, Y: -236.62, S: "4,162.00"},
		{X: 477.12, Y: -236.62, S: "8.00"},
		{X: 538.12, Y: -236.62, S: "4,170.00"},
	}
	if !isOpeningBalance(line) {
		t.Fatalf("got %v, wanted %v", true, false)
	}
}

func TestIsNotOpeningBalance(t *testing.T) {
	var line pdf.TextHorizontal = []pdf.Text{
		{X: 13, Y: -494, S: "29/12/2022 02:14:05"},
		{X: 112, Y: -494, S: "Bag Purchase           COCHIN"},
		{X: 542, Y: -494, S: "2000.00"},
		{X: 562.5, Y: -494, S: "  "},
	}
	if isOpeningBalance(line) {
		t.Fatalf("got %v, wanted %v", false, true)
	}
	var line2 pdf.TextHorizontal = []pdf.Text{
		{X: 16, Y: -236, S: "Welcome"},
		{X: 264, Y: -236.62, S: "Hello"},
		{X: 332.12, Y: -236.62, S: "  "},
		{X: 402.12, Y: -236.62, S: "123.00"},
		{X: 477.12, Y: -236.62, S: "  "},
		{X: 538.12, Y: -236.62, S: "Thank you"},
	}
	if isOpeningBalance(line2) {
		t.Fatalf("got %v, wanted %v", false, true)
	}
}
