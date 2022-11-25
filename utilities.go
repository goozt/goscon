package goscon

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func cleanString(data string) string {
	data = strings.TrimSpace(data)
	spaceExp := regexp.MustCompile(`( ){2,}`)
	data = spaceExp.ReplaceAllString(data, " ")
	commaExp := regexp.MustCompile(`(,)+|( )*,( )*`)
	return commaExp.ReplaceAllString(data, ",")
}

func cleanPath(filename string) string {
	file := filepath.Clean(filename)
	file, _ = filepath.Abs(file)
	return file
}

func isPdfFile(file string) bool {
	return strings.ToLower(filepath.Ext(file)) == ".pdf"
}

func parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	return strconv.ParseFloat(strings.Replace(s, ",", "", -1), 64)
}
