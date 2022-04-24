package replacer

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type Replacer struct {
	from    string
	to      string
	verbose bool
	stdout  io.Writer
	stderr  io.Writer
}

type ReplacerOption struct {
	From    string
	To      string
	Verbose bool
	Stdout  io.Writer
	Stderr  io.Writer
}

func New(p ReplacerOption) *Replacer {
	return &Replacer{from: p.From, to: p.To, verbose: p.Verbose, stderr: p.Stderr, stdout: p.Stdout}
}

func (r *Replacer) Run(wg *sync.WaitGroup, path string) {
	defer wg.Done()

	f, err := os.Open(path)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(f)
	needUpdate := false

	for i := 1; scanner.Scan(); i++ {
		t := scanner.Text()
		if strings.Contains(t, r.from) {
			fmt.Fprintf(r.stdout, "Replace %s:%d: %s\n", path, i, t)
			needUpdate = true
		}
	}

	f.Close()
	if !needUpdate {
		return
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(r.stderr, "file read failed: %v\n", err)
	}

	newData := strings.ReplaceAll(string(data), r.from, r.to)

	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(r.stderr, "file stat failed: %v\n", err)
	}

	err = ioutil.WriteFile(path, []byte(newData), info.Mode())
	if err != nil {
		fmt.Fprintf(r.stderr, "file update failed: %v\n", err)
	}
}
