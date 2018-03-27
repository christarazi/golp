package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func read_file(filename string) []byte {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		println("File", filename, "does not exist.")
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		println("Failed to read file:", err)
		os.Exit(1)
	}

	return data
}

func parse(content []string) ([]string, []string) {
	restr := "(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) - - " +
		"\\[(\\d{1,2}\\/\\w{3}\\/\\d{4}):(\\d{2}:\\d{2}:\\d{2}).+" +
		"(\\\"GET.+\\\"){1,3}"

	matches := []string{}
	nonmatches := []string{}

	regex, _ := regexp.Compile(restr)
	for _, v := range content {
		if len(regex.FindString(v)) != 0 {
			matches = append(matches, v)
		} else {
			nonmatches = append(nonmatches, v)
		}
	}

	return matches, nonmatches
}

func main() {
	if len(os.Args) != 2 {
		println("Usage:", os.Args[0], "<file>")
		os.Exit(1)
	}

	data := string(read_file(os.Args[1]))
	lines := strings.Split(data, "\n")

	m, n := parse(lines)

	fmt.Printf("Matches:     %v\n"+
		"Nonmatches:  %v\n"+
		"Total lines: %v\n", len(m), len(n), len(lines))
}
