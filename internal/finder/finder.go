package finder

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/y-yagi/goext/osext"
	"github.com/y-yagi/kaeru/internal/replacer"
)

// TODO(y-yagi): allow configuring
const ignoreFile = ".gitignore"

var ignoreDirectories = map[string]bool{
	".git":         true,
	"log":          true,
	"tmp":          true,
	"node_modules": true,
}

type Finder struct {
	replacer           replacer.Replacer
	pattern            string
	stderr             io.Writer
	stdout             io.Writer
	gitignore          *ignore.GitIgnore
	appendedIgnore     *ignore.GitIgnore
	appendedIgnoreFile string
}

type FinderOption struct {
	Replacer           replacer.Replacer
	Pattern            string
	AppendedIgnoreFile string
	Stdout             io.Writer
	Stderr             io.Writer
}

func New(f FinderOption) *Finder {
	return &Finder{replacer: f.Replacer, pattern: f.Pattern, stderr: f.Stderr, stdout: f.Stdout, appendedIgnoreFile: f.AppendedIgnoreFile}
}

func (f *Finder) Run() error {
	var wg sync.WaitGroup
	var err error

	if len(f.appendedIgnoreFile) != 0 {
		if !osext.IsExist(f.appendedIgnoreFile) {
			return fmt.Errorf("ignore file did'nt find: '%s'", f.appendedIgnoreFile)
		}

		f.appendedIgnore, err = ignore.CompileIgnoreFile(f.appendedIgnoreFile)
		if err != nil {
			return err
		}
	}

	if osext.IsExist(ignoreFile) {
		f.gitignore, err = ignore.CompileIgnoreFile(ignoreFile)
		if err != nil {
			return err
		}
	}

	err = filepath.Walk(".", func(p string, fi os.FileInfo, err error) error {
		if f.isIgnorePath(p) {
			if fi.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if fi.IsDir() {
			return nil
		}

		if len(f.pattern) != 0 {
			if matched, _ := path.Match(f.pattern, path.Base(p)); !matched {
				return nil
			}
		}

		wg.Add(1)
		go f.replacer.Run(&wg, p)
		return nil
	})

	wg.Wait()

	return err
}

func (f *Finder) isIgnorePath(path string) bool {
	if _, found := ignoreDirectories[path]; found {
		return true
	}

	if f.gitignore != nil && f.gitignore.MatchesPath(path) {
		return true
	}

	if f.appendedIgnore != nil && f.appendedIgnore.MatchesPath(path) {
		return true
	}

	return false
}
