package replacer

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

type Replacer struct {
	from    string
	to      string
	verbose bool
	regexp  bool
	dryrun  bool
	stdout  io.Writer
	stderr  io.Writer
}

type ReplacerOption struct {
	From    string
	To      string
	Verbose bool
	Regexp  bool
	Dryrun  bool
	Stdout  io.Writer
	Stderr  io.Writer
}

func New(p ReplacerOption) *Replacer {
	return &Replacer{from: p.From, to: p.To, verbose: p.Verbose, stderr: p.Stderr, stdout: p.Stdout, regexp: p.Regexp, dryrun: p.Dryrun}
}

func (r *Replacer) Run(wg *sync.WaitGroup, path string) {
	defer wg.Done()

	f, err := os.Open(path)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(f)
	needUpdate := false
	data := ""

	var matcher Matcher
	if r.regexp {
		matcher = &RegexpMatcher{from: r.from, to: r.to}
	} else {
		matcher = &StringMatcher{from: r.from, to: r.to}
	}

	for i := 1; scanner.Scan(); i++ {
		t := scanner.Text()
		if matcher.match(t) {
			if r.verbose {
				fmt.Fprintf(r.stdout, "Replace %s:%d: %s\n", path, i, matcher.colorizeFrom(t))
			}
			needUpdate = true
		}
		data += t + "\n"
	}

	f.Close()

	if r.dryrun || !needUpdate {
		return
	}

	newData := matcher.replace(string(data))

	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(r.stderr, "file stat failed: %v\n", err)
	}

	err = ioutil.WriteFile(path, []byte(newData), info.Mode())
	if err != nil {
		fmt.Fprintf(r.stderr, "file update failed: %v\n", err)
	}
}
