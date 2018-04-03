package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

type log_entry struct {
	ip        string
	date      string
	time      string
	timestamp time.Time
	action    string
}

func create_entry(ip, date, timestr, action []byte) log_entry {
	_date := string(date)
	_timestr := string(timestr)
	dt := _date + " " + _timestr
	ts, _ := time.Parse("02/Jan/2006 15:04:05", dt)

	return log_entry{string(ip),
		string(date),
		string(timestr),
		ts,
		string(action)}
}

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

func parse(content [][]byte) ([]log_entry, [][]byte) {
	restr := "(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}) - - " +
		"\\[(\\d{1,2}\\/\\w{3}\\/\\d{4}):(\\d{2}:\\d{2}:\\d{2}).+" +
		"(\\\"GET.+\\\"){1,3}"

	matches := []log_entry{}
	nonmatches := [][]byte{}

	regex, _ := regexp.Compile(restr)
	for _, v := range content {
		if len(regex.Find(v)) != 0 {
			submatches := regex.FindAllSubmatch(v, -1)[0]
			ip := submatches[1]
			date := submatches[2]
			time := submatches[3]
			action := submatches[4]
			matches = append(matches, create_entry(ip, date, time, action))
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

	data := read_file(os.Args[1])
	lines := bytes.Split(data, []byte("\n"))

	m, n := parse(lines)

	fmt.Printf("Matches:     %v\n"+
		"Nonmatches:  %v\n"+
		"Total lines: %v\n", len(m), len(n), len(lines))
}
