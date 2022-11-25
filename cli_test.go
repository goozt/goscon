package goscon

import (
	"path/filepath"
	"testing"
)

func TestSetFormat(t *testing.T) {
	o := CliOptions{}
	_ = o.SetFormat("csv")
	got := o.Format
	want := "csv"
	if got != want {
		t.Fatalf("got %v, wanted %v", got, want)
	}
	_ = o.SetFormat("json")
	got = o.Format
	want = "json"
	if got != want {
		t.Fatalf("got %v, wanted %v", got, want)
	}
	err := o.SetFormat("txt")
	if err == nil {
		t.Fatalf("got %v, wanted %v", o.Format, want)
	}
}

func TestCliApp(t *testing.T) {
	filename := ""
	dir := "sample"
	format := "csv"
	app, err := cliApp(filename, dir, format)
	if err != nil {
		t.Fatalf("error %+v", err)
	}
	got := app.IsBatch
	want := true
	if got != want {
		t.Fatalf("got %v, wanted %v", got, want)
	}
	got2 := app.Batch
	want2 := []string{}
	if len(got2) != len(want2) {
		t.Fatalf("got %v, wanted %v", got2, want2)
	}
	got3 := app.Dir
	want3, _ := filepath.Abs(dir)
	if got3 != want3 {
		t.Fatalf("got %v, wanted %v", got3, want3)
	}
	got4 := app.File
	want4 := ""
	if got4 != want4 {
		t.Fatalf("got %v, wanted %v", got4, want4)
	}
	got5 := app.Format
	want5 := "csv"
	if got5 != want5 {
		t.Fatalf("got %v, wanted %v", got5, want5)
	}

	filename = "test.pdf"
	dir = ""
	format = "csv"
	app, err = cliApp(filename, dir, format)
	if err != nil {
		t.Fatalf("error %+v", err)
	}
	got = app.IsBatch
	want = false
	if got != want {
		t.Fatalf("got %v, wanted %v", got, want)
	}
	got3, _ = filepath.Abs(app.Dir)
	want3, _ = filepath.Abs(dir)
	if got3 != want3 {
		t.Fatalf("got %v, wanted %v", got3, want3)
	}
	got4 = app.File
	want4, _ = filepath.Abs(filename)
	if got4 != want4 {
		t.Fatalf("got %v, wanted %v", got4, want4)
	}
}
