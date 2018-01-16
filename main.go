//
// expected file list
//
// <<EOF
// href="https://example.com"
// title="file title this is append suffix"
// middle title
// --
// EOF
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func get(url string, out string) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

const version = "0.0.0"
const name = "gather"

type option struct {
	version bool
	dryRun  bool
	dir     string
	list    string

	randomDelay bool
	delayMin    int64
	delayMax    int64

	trimNumber bool
}

var opt option

func init() {
	flag.BoolVar(&opt.version, "version", false, "")
	flag.BoolVar(&opt.dryRun, "dry-run", true, "dry-run")
	flag.StringVar(&opt.dir, "dir", "", "specify output directory")
	flag.StringVar(&opt.list, "list", "", "specify list file")

	flag.BoolVar(&opt.randomDelay, "delay", true, "random delay between -delay-min to -delay-max")
	flag.Int64Var(&opt.delayMin, "delay-min", 1, "random delay minimal seconds")
	flag.Int64Var(&opt.delayMax, "delay-max", 10, "random delay maximal seconds")

	flag.BoolVar(&opt.trimNumber, "trim-number", false, "trim prefix number(TODO: implementation)")
}

func main() {
	// TODO: consider exit code
	var exitcode int
	// TODO: consider logger
	logger := log.New(os.Stderr, "["+name+"]:", log.LstdFlags)

	// flag
	flag.Parse()
	if flag.NArg() == 1 && opt.list == "" {
		opt.list = flag.Arg(0)
	} else if flag.NArg() != 0 {
		log.Fatalf("unknown arguments: %v", flag.Args())
	}

	// version
	if opt.version {
		fmt.Printf("%s version %s\n", name, version)
	}

	// output directory
	outdir, err := filepath.Abs(opt.dir)
	if err != nil {
		logger.Fatal(err)
	}
	if err := os.MkdirAll(outdir, 0777); err != nil {
		logger.Fatal(err)
	}

	// list
	if opt.list == "" {
		flag.Usage()
		logger.Fatal("expected list file")
	}
	f, err := os.Open(opt.list)
	if err != nil {
		logger.Fatal(err)
	}
	defer f.Close()

	// loop
	sc := bufio.NewScanner(f)
	var i int
	var (
		url, midTitle, title, out string
		trunc                     = func() {
			url = ""
			midTitle = ""
			title = ""
			out = ""
		}
		mkout = func(out *string, add string) {
			switch {
			case *out == "":
				*out = add
			default:
				*out = fmt.Sprintf("%s-%s", *out, add)
			}
		}
	)
	for sc.Scan() {
		if sc.Err() != nil {
			logger.Fatal(sc.Err())
		}

		// get!
		text := strings.TrimSpace(sc.Text())
		switch {
		case text == "":
			continue
		case text == "--" && url != "":
			i++

			// TODO: consider
			// name
			if !opt.trimNumber {
				out = fmt.Sprintf("%d", i)
			}
			if midTitle != "" {
				mkout(&out, midTitle)
			}
			if title != "" {
				mkout(&out, title)
			}
			if out == "" {
				out = fmt.Sprintf("%d", i)
			}
			out, err := filepath.Abs(filepath.Join(outdir, out))
			if err != nil {
				logger.Fatal(err)
			}

			// check
			if _, err := os.Stat(out); err == nil {
				logger.Printf("[get %d]: [still exists]: %s [skipped url]: %s", i, out, url)
				exitcode = 2
				trunc()
				continue
			}

			// get
			delay := time.Second * time.Duration(opt.delayMin+rand.Int63n(opt.delayMax))
			if opt.dryRun {
				logger.Printf("[dry-run get %d]:\n\t[url]: %s\n\t[out]: %s\n\t[delay]: %v\n", i, url, out, delay)
			} else {
				logger.Printf("[get %d]:\n\t[url]: %s\n\t[out]: %s\n\t[delay]: %v\n", i, url, out, delay)
				if err := get(url, out); err != nil {
					logger.Fatal(err)
				}
				time.Sleep(delay)
			}
			trunc()
		case strings.HasPrefix(text, "href=\""):
			url = strings.TrimPrefix(text, "href=\"")
			url = strings.TrimSuffix(url, "\"")
		case strings.HasPrefix(text, "title=\""):
			title = strings.TrimPrefix(text, "title=\"")
			title = strings.TrimSuffix(title, "\"")
		default:
			midTitle = text
		}
	}
	os.Exit(exitcode)
}
