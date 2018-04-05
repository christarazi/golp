package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"regexp"
	"sort"
	"time"
)

type log_entry struct {
	Ip        string
	Date      string
	Time      string
	Timestamp time.Time
	Action    string
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
	return t[i].Timestamp.Unix() < t[j].Timestamp.Unix()
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

func group_by(entries []log_entry, field reflect.StructField) [][]log_entry {

	// This function is quite ugly because of Golang's lack of generics. It
	// would have been extremely useful here.

	if len(entries) == 1 {
		return [][]log_entry{{entries[0]}}
	}

	i := 0 // Using a two iterator approach here.
	j := 1

	first := reflect.ValueOf(entries[i])
	second := reflect.ValueOf(entries[j])

	subgroup := []log_entry{first.Interface().(log_entry)}
	grouped := [][]log_entry{}

	for i < len(entries) && j < len(entries) {
		second = reflect.ValueOf(entries[j])

		a := first.FieldByName(field.Name).Interface()
		b := second.FieldByName(field.Name).Interface()

		if a == b {
			subgroup = append(subgroup, second.Interface().(log_entry))
		} else {
			grouped = append(grouped, subgroup)

			i = j
			first = reflect.ValueOf(entries[i])

			subgroup = nil // Clear slice for new subgroup.
			subgroup = append(subgroup, first.Interface().(log_entry))
		}

		j += 1
	}

	// Add remaining elements from subgroup.
	if len(subgroup) > 0 {
		grouped = append(grouped, subgroup)
	}

	return grouped
}

func output(resolvehost *bool, entries []log_entry, field reflect.StructField) {
	for _, g := range group_by(entries, field) {
		if len(g) == 0 {
			continue
		}

		fmt.Printf("=====\n")
		for _, v := range g {
			fmt.Printf("Timestamp: %v\n", v.Timestamp)

			if *resolvehost {
				names, err := net.LookupAddr(v.Ip)
				if err != nil {
					fmt.Printf("IP:        %v\nAction:    %v", v.Ip, v.Action)
				} else {
					fmt.Printf("Hostname:  %v\nAction:    %v", names[0],
						v.Action)
				}
			} else {
				fmt.Printf("IP:        %v\nAction:    %v", v.Ip, v.Action)
			}
			fmt.Println("")
		}
		fmt.Println("=====\n")
	}
}

func main() {
	// Set command line arg flags.
	file := flag.String("file", "", "log file to analyze/parse")
	resolvehost := flag.Bool("resolve", false, "resolve ip addr to hostnames")
	group := flag.String("group", "Ip",
		"category to group entries by (default ip address)")
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

	// Determine struct field to group the entries by.
	// Get element type since |m| is a Slice.
	mirror := reflect.TypeOf(m).Elem()
	field, success := mirror.FieldByName(*group)
	if !success {
		println("Unknown group type specified. Supported: ")
		for i := 0; i < mirror.NumField(); i++ {
			println(mirror.Field(i).Name)
		}
		os.Exit(1)
	}

	output(resolvehost, m, field)
}
