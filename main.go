package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"sort"
	"time"
)

type log_entry struct {
	ip        string
	date      string
	time      string
	timestamp time.Time
	action    string
}

// Methods for sort.Interface.
type byTimestamp []log_entry

func (t byTimestamp) Len() int {
	return len(t)
}

func (t byTimestamp) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t byTimestamp) Less(i, j int) bool {
	return t[i].timestamp.Unix() < t[j].timestamp.Unix()
}

// End of sort.Interface.

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
	// Set command line arg flags.
	file := flag.String("file", "", "log file to analyze/parse")
	resolvehost := flag.Bool("resolve", false, "resolve ip addr to hostnames")
	flag.Parse()

	if len(os.Args) < 2 || flag.NFlag() == 0 {
		println("Usage:", os.Args[0], "[OPTIONS] --file <file>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	data := read_file(*file)
	lines := bytes.Split(data, []byte("\n"))

	m, n := parse(lines)

	fmt.Printf("Matches:     %v\n"+
		"Nonmatches:  %v\n"+
		"Total lines: %v\n\n", len(m), len(n), len(lines))

	sort.Sort(byTimestamp(m))
	for _, v := range m {
		fmt.Printf("%v\n", v.timestamp)

		if *resolvehost {
			names, err := net.LookupAddr(v.ip)
			if err != nil {
				fmt.Printf("%v\n%v", v.ip, v.action)
			} else {
				fmt.Printf("%v\n%v", names[0], v.action)
			}
		} else {
			fmt.Printf("%v\n%v", v.ip, v.action)
		}
		fmt.Println("\n")
	}
}
