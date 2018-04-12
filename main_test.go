package main

import (
	"bytes"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

const (
	EXPECTED_MATCHES    = 8
	EXPECTED_NONMATCHES = 2
)

var data []byte

func TestReadFile(t *testing.T) {
	data = read_file("fixtures/access.log")
	if len(data) == 0 {
		t.Error("Log is empty")
	}
}

func TestParse(t *testing.T) {
	lines := bytes.Split(data, []byte("\n"))
	matches, nonmatches := parse(lines)

	if len(nonmatches) != EXPECTED_NONMATCHES {
		t.Errorf("Mismatch with nonmatches: expected %v, got %v\n",
			EXPECTED_NONMATCHES, len(nonmatches))
	}

	if len(matches) != EXPECTED_MATCHES {
		t.Errorf("Mismatch with matches: expected %v, got %v\n",
			EXPECTED_MATCHES, len(matches))
	}

	expected_entries := []log_entry{
		log_entry{Ip: "95.213.130.90", Date: "08/Apr/2018", Time: "07:54:55",
			Action: "\"GET /_asterisk/ HTTP/1.1\" 404 136 \"-\" \"python-requests/2.18.4\"",
			Request: request{Method: "GET", Endpoint: "/_asterisk/", HTTPVersion: "HTTP/1.1",
				ResponseCode: "404", Reserved: "136", UserAgent: "python-requests/2.18.4"}},

		log_entry{Ip: "184.105.139.70", Date: "08/Apr/2018", Time: "09:26:24",
			Action: "\"GET / HTTP/1.1\" 200 3997 \"-\" \"-\"",
			Request: request{Method: "GET", Endpoint: "/", HTTPVersion: "HTTP/1.1",
				ResponseCode: "200", Reserved: "3997", UserAgent: "-"}},

		log_entry{Ip: "216.218.206.66", Date: "08/Apr/2018", Time: "10:24:07",
			Action: "\"GET / HTTP/1.1\" 200 3997 \"-\" \"-\"",
			Request: request{Method: "GET", Endpoint: "/", HTTPVersion: "HTTP/1.1",
				ResponseCode: "200", Reserved: "3997", UserAgent: "-"}},

		log_entry{Ip: "185.234.15.88", Date: "08/Apr/2018", Time: "10:27:23",
			Action: "\"GET /xmlrpc.php HTTP/1.1\" 301 178 \"-\" \"Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.6) Gecko/20070725 Firefox/2.0.0.6\"",
			Request: request{Method: "GET", Endpoint: "/xmlrpc.php", HTTPVersion: "HTTP/1.1",
				ResponseCode: "301", Reserved: "178",
				UserAgent: "Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.6) Gecko/20070725 Firefox/2.0.0.6"}},

		log_entry{Ip: "201.6.107.126", Date: "08/Apr/2018", Time: "10:33:40",
			Action: "\"GET /admin/config.php HTTP/1.1\" 404 162 \"-\" \"curl/7.15.5 (x86_64-redhat-linux-gnu) libcurl/7.15.5 OpenSSL/0.9.8b zlib/1.2.3 libidn/0.6.5\"",
			Request: request{Method: "GET", Endpoint: "/admin/config.php", HTTPVersion: "HTTP/1.1",
				ResponseCode: "404", Reserved: "162",
				UserAgent: "curl/7.15.5 (x86_64-redhat-linux-gnu) libcurl/7.15.5 OpenSSL/0.9.8b zlib/1.2.3 libidn/0.6.5"}},

		log_entry{Ip: "159.203.121.40", Date: "08/Apr/2018", Time: "10:47:23",
			Action: "\"GET / HTTP/1.0\" 200 3997 \"-\" \"Mozilla/5.0 (compatible; NetcraftSurveyAgent/1.0; +info@netcraft.com)\"",
			Request: request{Method: "GET", Endpoint: "/", HTTPVersion: "HTTP/1.0",
				ResponseCode: "200", Reserved: "3997",
				UserAgent: "Mozilla/5.0 (compatible; NetcraftSurveyAgent/1.0; +info@netcraft.com)"}},

		log_entry{Ip: "185.234.15.88", Date: "08/Apr/2018", Time: "10:48:52",
			Action: "\"GET /xmlrpc.php HTTP/1.1\" 301 178 \"-\" \"Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.6) Gecko/20070725 Firefox/2.0.0.6\"",
			Request: request{Method: "GET", Endpoint: "/xmlrpc.php", HTTPVersion: "HTTP/1.1",
				ResponseCode: "301", Reserved: "178",
				UserAgent: "Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.6) Gecko/20070725 Firefox/2.0.0.6"}},

		log_entry{Ip: "5.79.69.140", Date: "08/Apr/2018", Time: "11:09:40",
			Action: "\"GET /recordings/ HTTP/1.1\" 400 264 \"-\" \"curl/7.29.0\"",
			Request: request{Method: "GET", Endpoint: "/recordings/", HTTPVersion: "HTTP/1.1",
				ResponseCode: "400", Reserved: "264", UserAgent: "curl/7.29.0"}},
	}

	for i, v := range expected_entries {
		expected_entries[i].Timestamp, _ = time.Parse("02/Jan/2006 15:04:05", v.Date+" "+v.Time)
	}

	for i, v := range matches {
		e := &expected_entries[i]
		if !cmp.Equal(*e, v) {
			t.Errorf("Mismatch in entry #%v", i)
			fmt.Printf("Entry #%v:\n%v\n", i, cmp.Diff(*e, v))
		}
	}
}
