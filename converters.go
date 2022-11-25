package goscon

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (s Statement) ToStrings() [][]string {
	records := [][]string{}
	records = append(records, []string{"DATE", "DESCRIPTION", "AMOUNT", "CREDITED"})
	for _, t := range s.Transactions {
		records = append(records, []string{
			t.Date.Format("02/01/2006 15:04:05 -0700"),
			t.Description,
			strconv.FormatFloat(t.Amount, 'f', 2, 64),
			strconv.FormatBool(t.Credited),
		})
	}
	return records
}

func (s Statement) WriteCSV(filename string) error {
	records := s.ToStrings()
	ext := filepath.Ext(filename)
	filename = strings.TrimSuffix(filename, ext) + ".csv"
	dir := filepath.Dir(filename)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	csvFile, err := os.Create(filename)
	defer func() { csvFile.Close() }()
	if err != nil {
		return err
	}
	w := csv.NewWriter(csvFile)
	w.WriteAll(records)

	if err = w.Error(); err != nil {
		return err
	}
	return nil
}

func (s Statement) WriteJSON() error {
	// TODO: JSON Converter
	return nil
}
