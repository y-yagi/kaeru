package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/y-yagi/kaeru/internal/finder"
	"github.com/y-yagi/kaeru/internal/replacer"
)

const cmd = "kaeru"

var (
	flags              *flag.FlagSet
	showVersion        bool
	filenamePattern    string
	appendedIgnoreFile string
	regexp             bool
	dryrun             bool
	quiet              bool

	version = "devel"
)

func main() {
	setFlags()
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}

func setFlags() {
	flags = flag.NewFlagSet(cmd, flag.ExitOnError)
	flags.BoolVar(&showVersion, "v", false, "print version number")
	flags.StringVar(&filenamePattern, "name", "", "file name pattern")
	flags.BoolVar(&regexp, "regexp", false, "treat FROM as a regexp")
	flags.BoolVar(&dryrun, "dry-run", false, "perform a trial run with no changes made")
	flags.BoolVar(&quiet, "quiet", false, "suppress output")
	flags.StringVar(&appendedIgnoreFile, "ignore-file", "", "pass a file name that specified ignores pattern")
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
	if err := flags.Parse(args[1:]); err != nil {
		return msg(err, stderr)
	}

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

	if flags.Arg(0) == flags.Arg(1) {
		fmt.Fprintln(stderr, "don't use the same value for FROM and TO.")
		return 1
	}

	r := replacer.New(
		replacer.ReplacerOption{From: flags.Arg(0), To: flags.Arg(1), Quiet: quiet,
			Stdout: stdout, Stderr: stderr, Regexp: regexp, Dryrun: dryrun},
	)
	f := finder.New(finder.FinderOption{Replacer: r, Pattern: filenamePattern, Stdout: stdout, Stderr: stderr, AppendedIgnoreFile: appendedIgnoreFile})

	if err := f.Run(); err != nil {
		return msg(err, stderr)
	}

	return 0
}
