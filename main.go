package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/y-yagi/kaeru/internal/replacer"
)

const cmd = "kaeru"

var (
	flags           *flag.FlagSet
	showVersion     bool
	enableVerbose   bool
	filenamePattern string

	version = "devel"
)

func main() {
	setFlags()
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}

func setFlags() {
	flags = flag.NewFlagSet(cmd, flag.ExitOnError)
	flags.BoolVar(&showVersion, "v", false, "print version number")
	flags.BoolVar(&enableVerbose, "verbose", false, "enable versbose log")
	flags.StringVar(&filenamePattern, "name", "", "file name pattern")
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

	if len(filenamePattern) != 0 {
		if _, err := path.Match(filenamePattern, ""); err != nil {
			fmt.Fprintf(stderr, "invalid file name pattern is specified: %v\n", err)
			return 1
		}
	}

	var wg sync.WaitGroup
	r := replacer.New(replacer.ReplacerOption{From: flags.Arg(0), To: flags.Arg(1), Verbose: enableVerbose, Stdout: stdout, Stderr: stderr})
	filepath.Walk(".", func(p string, f os.FileInfo, err error) error {
		if p != "." && strings.HasPrefix(p, ".") {
			if f.IsDir() {
				return filepath.SkipDir
			}
		}

		if f.IsDir() {
			return nil
		}

		if len(filenamePattern) != 0 {
			if matched, _ := path.Match(filenamePattern, path.Base(p)); !matched {
				return nil
			}
		}

		wg.Add(1)
		go r.Run(&wg, p)
		return nil
	})

	wg.Wait()
	return 0
}
