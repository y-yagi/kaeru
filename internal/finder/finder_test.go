package finder_test

import (
	"os"
	"sync"
	"testing"

	"github.com/y-yagi/kaeru/internal/finder"
)

type TestReplacer struct {
	Files []string
}

func (r *TestReplacer) Run(wg *sync.WaitGroup, path string) {
	defer wg.Done()
	r.Files = append(r.Files, path)
}

func TestFinder_string(t *testing.T) {
	tempdir, err := os.MkdirTemp("", "findertest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempdir)

	testdir := tempdir + "/abc"
	if err = os.Mkdir(testdir, 0755); err != nil {
		t.Fatal(err)
	}

	testfile := testdir + "/dummy.log"
	if err = os.WriteFile(testfile, []byte("Hello, world"), 0644); err != nil {
		t.Fatal(err)
	}

	ignorefile := tempdir + "/.gitignore"
	if err = os.WriteFile(ignorefile, []byte("public/"), 0644); err != nil {
		t.Fatal(err)
	}

	ignoreDirs := []string{".git", "log", "public"}
	for _, ignoredir := range ignoreDirs {
		dir := tempdir + "/" + ignoredir
		if err = os.Mkdir(dir, 0755); err != nil {
			t.Fatal(err)
		}
		testfile := dir + "/dummy.log"
		if err = os.WriteFile(testfile, []byte("Hello, world\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	r := &TestReplacer{}
	f := finder.New(finder.FinderOption{Replacer: r, Pattern: "", Stdout: os.Stdout, Stderr: os.Stderr})
	os.Chdir(tempdir)
	f.Run()

	if len(r.Files) != 1 {
		t.Fatalf("Exepectd files are one, but got %+v\n", r.Files)
	}

	expected := "abc/dummy.log"
	if r.Files[0] != expected {
		t.Fatalf("Exepectd %+v, but got %+v\n", expected, r.Files[0])
	}
}
