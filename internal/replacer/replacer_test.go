package replacer_test

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/y-yagi/kaeru/internal/replacer"
)

func TestReplacer_string(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "replacertest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	testfile := tempDir + "/dummy.log"
	if err = ioutil.WriteFile(testfile, []byte("Hello, wolrd"), 0644); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	replacer := replacer.New(replacer.ReplacerOption{From: "wolrd", To: "world", Stdout: os.Stdout, Stderr: os.Stderr, Quiet: true})
	replacer.Run(&wg, testfile)

	newtext, err := ioutil.ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}

	expected := "Hello, world\n"
	if string(newtext) != expected {
		t.Fatalf("Exepectd \n\n%s\nbut got\n\n%s\n", expected, newtext)
	}
}

func TestReplacer_regexp(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "replacertest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	testfile := tempDir + "/dummy.log"
	if err = ioutil.WriteFile(testfile, []byte("Hello, world from 2022/01/31"), 0644); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	replacer := replacer.New(replacer.ReplacerOption{From: "(\\d{4})/(\\d{2})/(\\d{2})", To: "$2/$3/$1", Stdout: os.Stdout, Stderr: os.Stderr, Quiet: true, Regexp: true})
	replacer.Run(&wg, testfile)

	newtext, err := ioutil.ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}

	expected := "Hello, world from 01/31/2022\n"
	if string(newtext) != expected {
		t.Fatalf("Exepectd \n\n%s\nbut got\n\n%s\n", expected, newtext)
	}
}
