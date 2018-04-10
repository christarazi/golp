# golp
A web server log parser / analyzer written in Go for easier inspection of your
website's visitors.

### Usage

Simply download / clone the repo then run:

```shell
$ cd path/to/repo
$ go build && ./golp -file path/to/access.log

# ./golp -h for usage information
```

## TODO

 - Allow for more log parser formats (only nginx logs are supported for now)
 - ~~Add ability to convert timestamp into local time~~
 - ~~Breakdown the "action" into: method, endpoint, user-agent, etc.~~
 - ~~Show unmatched lines (--verbose option?)~~

## License

This program is free software, distributed under the terms of the [GNU] General
Public License as published by the Free Software Foundation, version 3 of the
License (or any later version).  For more information, see the file LICENSE.
