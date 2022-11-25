package goscon

import (
	"path/filepath"
	"testing"
)

func TestCleanString(t *testing.T) {
	got := cleanString("  234.87  ABC  , DE  334 50 ,dsg 234")
	want := "234.87 ABC,DE 334 50,dsg 234"
	if got != want {
		t.Fatalf("got %v, wanted %v", got, want)
	}
}

func TestCleanPath(t *testing.T) {
	dir, _ := filepath.Abs(".")
	got := cleanPath("./sample/statements")
	want := dir + "/sample/statements"
	if got != want {
		t.Fatalf("got %v, wanted %v", got, want)
	}
}

func TestIsPdfFile(t *testing.T) {
	g1 := isPdfFile("file.pdf")
	g2 := isPdfFile("file.txt")
	if !g1 {
		t.Errorf("got %v, wanted %v", g1, true)
	}
	if g2 {
		t.Errorf("got %v, wanted %v", g2, false)
	}
}

func TestParseFloat(t *testing.T) {
	got, err := parseFloat(" 34,45,324.54    ")
	if err != nil {
		t.Fatal(err)
	}
	want := 3445324.54
	if got != want {
		t.Fatalf("got %v, wanted %v", got, want)
	}
}
