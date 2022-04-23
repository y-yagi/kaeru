package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const cmd = "kaeru"

var (
	flags       *flag.FlagSet
	showVersion bool

	version = "devel"
)

func main() {
	setFlags()
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}

func setFlags() {
	flags = flag.NewFlagSet(cmd, flag.ExitOnError)
	flags.BoolVar(&showVersion, "v", false, "print version number")
	flags.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] FROM TO\n\n", cmd)
	fmt.Fprintln(os.Stderr, "OPTIONS:")
	flags.PrintDefaults()
}

func msg(err error, stderr io.Writer) int {
	if err != nil {
		fmt.Fprintf(stderr, "%s: %v\n", cmd, err)
		return 1
	}
	return 0
}

func run(args []string, stdout, stderr io.Writer) int {
	flags.Parse(args[1:])

	if showVersion {
		fmt.Fprintf(stdout, "%s %s\n", cmd, version)
		return 0
	}

	if len(flags.Args()) != 2 {
		flags.Usage()
		return 1
	}

	var wg sync.WaitGroup
	filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if path != "." && strings.HasPrefix(path, ".") {
			if f.IsDir() {
				return filepath.SkipDir
			}
		} else if !f.IsDir() {
			wg.Add(1)
			go replace(&wg, path, flags.Arg(0), flags.Arg(1))
		}
		return nil
	})

	wg.Wait()
	return 0
}

func replace(wg *sync.WaitGroup, path, from, to string) {
	defer wg.Done()

	f, err := os.Open(path)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(f)
	needUpdate := false

	for i := 1; scanner.Scan(); i++ {
		t := scanner.Text()
		if strings.Contains(t, from) {
			fmt.Printf("Replace %s:%d: %s\n", path, i, t)
			needUpdate = true
		}
	}

	f.Close()
	if !needUpdate {
		return
	}

	// TODO: error check
	data, _ := ioutil.ReadFile(path)
	newData := strings.ReplaceAll(string(data), from, to)

	info, _ := os.Stat(path)
	ioutil.WriteFile(path, []byte(newData), info.Mode())
}
