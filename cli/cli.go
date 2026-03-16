package cli

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/goozt/goscon"
	"github.com/urfave/cli/v2"
)

type Options struct {
	IsBatch bool
	Batch   []string
	File    string
	Dir     string
	Format  string
}

func (c *Options) SetFormat(format string) error {
	formatsAvailable := []string{"csv", "json"}
	for _, f := range formatsAvailable {
		if format == f {
			c.Format = f
			return nil
		}
	}
	return errors.New("error: invalid format flag")
}

func loadFilenames(dir string, ch chan string) error {
	defer close(ch)
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			ch <- file.Name()
		}
	}
	return nil
}

func cliApp(filename string, dir string, format string) (*Options, error) {
	co := Options{}
	err := co.SetFormat(format)
	if err != nil {
		return nil, err
	}
	if dir != "" {
		filenames := make(chan string, 1)
		dir = goscon.CleanPath(dir)

		errCh := make(chan error, 1)
		go func() {
			errCh <- loadFilenames(dir, filenames)
		}()

		for filename := range filenames {
			if !goscon.IsPdfFile(filename) {
				continue
			}
			co.Batch = append(co.Batch, filepath.Join(dir, filename))
		}

		if err := <-errCh; err != nil {
			return nil, err
		}

		co.Dir = dir
		co.IsBatch = true
	} else {
		if !goscon.IsPdfFile(filename) {
			return nil, errors.New("invalid file format")
		}
		co.File = goscon.CleanPath(filename)
		co.Dir = filepath.Dir(filename)
	}
	return &co, nil
}

func Run() (*Options, error) {
	var dir string
	var format string
	var filename string
	app := &cli.App{
		Name:     "Statement Converter",
		Version:  "0.0.3",
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
