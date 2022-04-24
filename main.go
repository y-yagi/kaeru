package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/y-yagi/kaeru/internal/replacer"
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
	r := replacer.New(replacer.ReplacerOption{From: flags.Arg(0), To: flags.Arg(1), Verbose: false, Stdout: stdout, Stderr: stderr})
	filepath.Walk(".", func(path string, f os.FileInfo, err error) error {
		if path != "." && strings.HasPrefix(path, ".") {
			if f.IsDir() {
				return filepath.SkipDir
			}
		} else if !f.IsDir() {
			wg.Add(1)
			go r.Run(&wg, path)
		}
		return nil
	})

	wg.Wait()
	return 0
}
