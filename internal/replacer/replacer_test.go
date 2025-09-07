package replacer_test

import (
	"os"
	"sync"
	"testing"

	"github.com/y-yagi/kaeru/internal/replacer"
)

func TestReplacer_string(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "replacertest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir) //nolint:errcheck

	testfile := tempDir + "/dummy.log"
	if err = os.WriteFile(testfile, []byte("Hello, wolrd"), 0644); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	replacer := replacer.New(replacer.ReplacerOption{From: "wolrd", To: "world", Stdout: os.Stdout, Stderr: os.Stderr, Quiet: true})
	replacer.Run(&wg, testfile)

	newtext, err := os.ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}

	expected := "Hello, world"
	if string(newtext) != expected {
		t.Fatalf("Exepectd \n\n%s\nbut got\n\n%s\n", expected, newtext)
	}
}

func TestReplacer_regexp(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "replacertest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir) //nolint:errcheck

	testfile := tempDir + "/dummy.log"
	if err = os.WriteFile(testfile, []byte("Hello, world from 2022/01/31"), 0644); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	replacer := replacer.New(replacer.ReplacerOption{From: "(\\d{4})/(\\d{2})/(\\d{2})", To: "$2/$3/$1", Stdout: os.Stdout, Stderr: os.Stderr, Quiet: true, Regexp: true})
	replacer.Run(&wg, testfile)

	newtext, err := os.ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}

	expected := "Hello, world from 01/31/2022"
	if string(newtext) != expected {
		t.Fatalf("Exepectd \n\n%s\nbut got\n\n%s\n", expected, newtext)
	}
}
func TestReplacer_preserve_trailing_newline(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "replacertest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir) //nolint:errcheck

	testfile := tempDir + "/dummy.log"
	// File ends with a newline
	if err = os.WriteFile(testfile, []byte("foo\nbar\n"), 0644); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	replacer := replacer.New(replacer.ReplacerOption{From: "bar", To: "baz", Stdout: os.Stdout, Stderr: os.Stderr, Quiet: true})
	replacer.Run(&wg, testfile)

	newtext, err := os.ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}

	expected := "foo\nbaz\n"
	if string(newtext) != expected {
		t.Fatalf("Expected trailing newline to be preserved.\nExpected: %q\nGot: %q", expected, newtext)
	}
}

func TestReplacer_omit_trailing_newline(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "replacertest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir) //nolint:errcheck

	testfile := tempDir + "/dummy.log"
	// File does not end with a newline
	if err = os.WriteFile(testfile, []byte("foo\nbar"), 0644); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	replacer := replacer.New(replacer.ReplacerOption{From: "bar", To: "baz", Stdout: os.Stdout, Stderr: os.Stderr, Quiet: true})
	replacer.Run(&wg, testfile)

	newtext, err := os.ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}

	expected := "foo\nbaz"
	if string(newtext) != expected {
		t.Fatalf("Expected no trailing newline to be added.\nExpected: %q\nGot: %q", expected, newtext)
	}
}
