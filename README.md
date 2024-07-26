# gobarchar

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/usrme/gobarchar)
[![test](https://github.com/usrme/gobarchar/actions/workflows/test.yml/badge.svg)](https://github.com/usrme/gobarchar/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/usrme/gobarchar)](https://goreportcard.com/report/github.com/usrme/gobarchar)

For whatever reason, I was enamored with [Alex Chan's snippet of code](https://alexwlchan.net/2018/ascii-bar-charts/) where a key-value pair list is turned into a passable bar chart for quick demonstration purposes. When writing a [blog post mentioning the number of books I've read](https://usrme.xyz/posts/glee-in-2023/#books-read) I wanted to quickly visualize the values, but didn't want to spend too much time on it. The solution above worked wonders! But I wanted the same thing without having to actually open a terminal (heresy, I know), thus this thing was born: the charting solution that might not suit you 📊

![GoBarChar - animated GIF demo](examples/demo.gif)

Try the same link as in the demo above, https://gobarchar.usrme.xyz/?2012=8&2013=6&2014=8&2015=14, or just load the site without any parameters to get random data.

## Features

- Send data and get chart back
- Perform calculation for average to the nearest integer
- Perform calculation for sum of all values
- Sort ascending or descending, or don't sort at all
  - Add the `sort` query parameter and pass either `asc` or `desc` as the value
  - The default is to keep the rows ordered as the query parameters are
- Replace `%20` HTML entity with a space
  - Add the `spaces` query parameter and pass `yes` as the value
    - Any other value will be disregarded
  - The default is to not do any replacing
- Maybe coming: change layout from horizontal to vertical!

## Usage

### On the Fly

(Excuse the pun). Instead of installing the thing, you can also access the same functionality using the instance hosted on [Fly](https://fly.io/) at https://gobarchar.usrme.xyz/:

```console
$ curl https://gobarchar.usrme.xyz/
April        7 █▊
February    88 ██████████████████████▏
June        16 ████
March       99 █████████████████████████
October     19 ████▊
September   98 ████████████████████████▋
Avg.        55 █████████████▉
Total      327 █████████████████████████
```

Loading the page without any query parameters randomly chooses some data and will do so upon every reload; a link will be presented at the bottom to link to the data. The site should always be available for usage, but if it's not then do open up an issue and I'll see what I can do.

### Local

After [installation](#installation), execute `gobarchar`, which by default starts a web server listening on port 8080, though you can specify a different port through the `PORT` environment variable. Once the web server is running, you can use something like `curl` to perform requests:

```console
$ gobarchar &
[1] 32357
2024/07/26 12:23:13 listening on: 8080
$ curl localhost:8080
2024/07/26 12:23:16 completed in: 87.958µs
August        32 ████████████████
February      50 █████████████████████████
July           7 ███▌
March          6 ███
November      44 ██████████████████████
September      8 ████
Avg.       24.50 ████████████▎
Total        147 █████████████████████████
```

If no query parameters are provided then random data is generated. The "completed in" log line isn't a part of the chart output, but rather the result of using [job control](https://jvns.ca/blog/2024/07/03/reasons-to-use-job-control/) and having the standard output of a background process being interleaved with the standard output of a foreground process.

## Installation

- using `go install`:

```bash
go install github.com/usrme/gobarchar/cmd/...@latest
```

- download a binary from the [releases](https://github.com/usrme/gobarchar/releases) page

- build it yourself (requires Go 1.17+):

```bash
git clone https://github.com/usrme/gobarchar.git
cd gobarchar
go build -o gobarchar cmd/gobarchar/main.go
```

## Removal

```bash
rm -f "${GOPATH}/bin/gobarchar"
rm -rf "${GOPATH}/pkg/mod/github.com/usrme/gobarchar*"
```

## Acknowledgments

Heavily inspired by the ['Drawing ASCII bar charts' blog post](https://alexwlchan.net/2018/ascii-bar-charts/) by Alex Chan. If there was any prior art that pretty much does the same thing (present an ASCII chart based on query parameters or request payload), then I honestly wasn't aware of it and just created this for fun.

## License

[MIT](/LICENSE)
