package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

var (
	csvFilename  string
	csvDelimiter string
	comma        rune
)

func main() {
	parseCL()

	var csvfd *os.File
	var err error
	var all [][]string

	if csvfd, err = os.Open(csvFilename); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer csvfd.Close()

	reader := csv.NewReader(csvfd)
	reader.Comma = comma
	reader.LazyQuotes = true

	fmt.Printf("Read file: %s\n", csvFilename)
	if all, err = reader.ReadAll(); err != nil {
		fmt.Println(err)
		return
	}

	arr := make([]Patient, len(all))
	for i, rec := range all {
		p := NewPatient(rec)
		p.ComputeMd5()
		arr[i] = *p
	}

	if out, err := json.Marshal(arr); err == nil {
		fmt.Println(string(out))
	}
}

func parseCL() {
	flag.StringVar(&csvFilename, "f", "patients.csv", "a file to import (.csv)")
	flag.StringVar(&csvDelimiter, "s", "C", "field separator to use (C=Comma, T=Tab)")
	flag.Parse()

	// -----------------------------
	switch csvDelimiter {
	case "C":
		comma = ','
	default:
		comma = '\t'
	}
}
