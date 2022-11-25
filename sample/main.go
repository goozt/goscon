package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/goozt/goscon"
)

func main() {
	cli, err := goscon.Cli()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if cli.IsBatch {
		for _, file := range cli.Batch {
			statement, err := goscon.Read(file)
			if err != nil {
				log.Fatalf("%+v", err)
			}
			dir, filename := filepath.Split(file)
			err = statement.WriteCSV(dir + "csv/" + filename)
			if err != nil {
				log.Fatalf("%+v", err)
			}
		}
	} else {
		statement, err := goscon.Read(cli.File)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		err = statement.WriteCSV(cli.File)
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}
	fmt.Println("Converted to", cli.Format)
}
