package replacer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
)

type Replacer interface {
	Run(wg *sync.WaitGroup, path string)
}

type FileReplacer struct {
	from   string
	to     string
	regexp bool
	dryrun bool
	quiet  bool
	stdout io.Writer
	stderr io.Writer
}

type ReplacerOption struct {
	From    string
	To      string
	Verbose bool
	Regexp  bool
	Dryrun  bool
	Quiet   bool
	Stdout  io.Writer
	Stderr  io.Writer
}

func New(p ReplacerOption) *FileReplacer {
	return &FileReplacer{from: p.From, to: p.To, quiet: p.Quiet, stderr: p.Stderr, stdout: p.Stdout, regexp: p.Regexp, dryrun: p.Dryrun}
}

func (r *FileReplacer) Run(wg *sync.WaitGroup, path string) {
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
			if !r.quiet {
				fmt.Fprintf(r.stdout, "Replace %s:%d: %s\n", path, i, matcher.colorizeFrom(t)) //nolint:errcheck
			}
			needUpdate = true
		}
		data += t + "\n"
	}

	f.Close() //nolint:errcheck

	if r.dryrun || !needUpdate {
		return
	}

	newData := matcher.replace(string(data))

	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(r.stderr, "file stat failed: %v\n", err) //nolint:errcheck
	}

	err = os.WriteFile(path, []byte(newData), info.Mode())
	if err != nil {
		fmt.Fprintf(r.stderr, "file update failed: %v\n", err) //nolint:errcheck
	}
}
