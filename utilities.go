package goscon

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	spaceExp = regexp.MustCompile(`( ){2,}`)
	commaExp = regexp.MustCompile(`(,)+|( )*,( )*`)
)

func CleanString(data string) string {
	data = strings.TrimSpace(data)
	data = spaceExp.ReplaceAllString(data, " ")
	return commaExp.ReplaceAllString(data, ",")
}

func CleanPath(filename string) (string, error) {
	file := filepath.Clean(filename)
	return filepath.Abs(file)
}

func IsPdfFile(file string) bool {
	return strings.ToLower(filepath.Ext(file)) == ".pdf"
}

func ParseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	return strconv.ParseFloat(strings.Replace(s, ",", "", -1), 64)
}
