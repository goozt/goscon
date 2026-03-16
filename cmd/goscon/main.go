package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/goozt/goscon"
	"github.com/goozt/goscon/cli"
)

func main() {
	opts, err := cli.Run()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if opts.IsBatch {
		for _, file := range opts.Batch {
			statement, err := goscon.Read(file)
			if err != nil {
				log.Fatalf("%+v", err)
			}
			dir, filename := filepath.Split(file)
			outputDir := filepath.Join(dir, "converted")
			outputPath := filepath.Join(outputDir, filename)

			if opts.Format == "json" {
				err = statement.WriteJSON(outputPath)
			} else {
				err = statement.WriteCSV(outputPath)
			}

			if err != nil {
				log.Fatalf("%+v", err)
			}
		}
	} else {
		statement, err := goscon.Read(opts.File)
		if err != nil {
			log.Fatalf("%+v", err)
		}

		if opts.Format == "json" {
			err = statement.WriteJSON(opts.File)
		} else {
			err = statement.WriteCSV(opts.File)
		}

		if err != nil {
			log.Fatalf("%+v", err)
		}
	}
	fmt.Println("Converted to", opts.Format)
}
