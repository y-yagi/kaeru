package finder

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/y-yagi/kaeru/internal/replacer"
)

type Finder struct {
	replacer *replacer.Replacer
	pattern  string
	stderr   io.Writer
	stdout   io.Writer
}

type FinderOption struct {
	Replacer *replacer.Replacer
	Pattern  string
	Stdout   io.Writer
	Stderr   io.Writer
}

func New(f FinderOption) *Finder {
	return &Finder{replacer: f.Replacer, pattern: f.Pattern, stderr: f.Stderr, stdout: f.Stdout}
}

func (f *Finder) Run() error {
	var wg sync.WaitGroup

	err := filepath.Walk(".", func(p string, fi os.FileInfo, err error) error {
		if p != "." && strings.HasPrefix(p, ".") {
			if fi.IsDir() {
				return filepath.SkipDir
			}
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
