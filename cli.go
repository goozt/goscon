package goscon

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/urfave/cli/v2"
)

type CliOptions struct {
	IsBatch bool
	Batch   []string
	File    string
	Dir     string
	Format  string
}

func (c *CliOptions) SetFormat(format string) error {
	formatsAvailable := []string{"csv", "json"}
	for _, f := range formatsAvailable {
		if format == f {
			c.Format = f
			return nil
		}
	}
	return errors.New("error: invalid format flag")
}

func loadFilenames(dir string, ch chan string) {
	defer close(ch)
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			ch <- file.Name()
		}
	}
}

func cliApp(filename string, dir string, format string) (*CliOptions, error) {
	co := CliOptions{}
	err := co.SetFormat(format)
	if err != nil {
		return nil, err
	}
	if dir != "" {
		filenames := make(chan string, 1)
		dir = cleanPath(dir)
		go loadFilenames(dir, filenames)
		for filename := range filenames {
			if !isPdfFile(filename) {
				continue
			}
			co.Batch = append(co.Batch, dir+"/"+filename)
		}
		co.Dir = dir
		co.IsBatch = true
	} else {
		if !isPdfFile(filename) {
			return nil, errors.New("invalid file format")
		}
		co.File = cleanPath(filename)
		co.Dir = filepath.Dir(filename)
	}
	return &co, nil
}

func Cli() (*CliOptions, error) {
	var dir string
	var format string
	var filename string
	app := &cli.App{
		Name:     "Statement Converter",
		Version:  "0.0.1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Nikhil John",
				Email: "me@nikz.in",
			},
		},
		Description: "A format converter for Bank Statement(PDF)",
		Usage:       "Converts Bank Statement(PDF) to desired structured formats",
		UsageText:   "goscon [filname|-d directory] [-f format]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "directory",
				Aliases:     []string{"d"},
				Usage:       "directory inwhich all statements are stored",
				Destination: &dir,
			},
			&cli.StringFlag{
				Name:        "format",
				Value:       "csv",
				Aliases:     []string{"f"},
				Usage:       "format to which the statement in converted",
				Destination: &format,
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.String("directory") == "" {
				values := ctx.Args()
				if values.Len() < 1 {
					return errors.New("missing filename")
				}
				filename = values.Get(0)
			}
			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	if err := app.Run(os.Args); err != nil {
		return nil, err
	}

	return cliApp(filename, dir, format)
}
