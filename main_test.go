package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wormi4ok/evernote2md/internal"
)

const sampleFile = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE en-export SYSTEM "http://xml.evernote.com/pub/evernote-export3.dtd">
<en-export>
<note>
  <title>Test</title>
  <content>
	<![CDATA[<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<!DOCTYPE en-note SYSTEM "http://xml.evernote.com/pub/enml2.dtd"><en-note><div><br /></div></en-note>]]>
  </content>
</note>
</en-export>
`

func Test_run(t *testing.T) {
	setLogLevel(false)
	tmpDir := t.TempDir()
	input := filepath.FromSlash(tmpDir + "/export.enex")
	err := ioutil.WriteFile(input, []byte(sampleFile), 0600)
	if err != nil {
		t.Fatalf("failed to create a test file at %s", input)
	}
	output := newNoteFilesDir(tmpDir, false, false)
	converter, _ := internal.NewConverter("", false)
	run(input, output, newProgressBar(false), converter)

	want := filepath.FromSlash(output.Path() + "/Test.md")
	_, err = os.Stat(want)
	if err != nil && os.IsNotExist(err) {
		t.Error("Test.md was not created")
	}
}

func Test_matchInput_cwd(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Create(filepath.Join(tmpDir, "test_export.enex")); err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(tmpDir, "test_export.enex")
	if got := matchInput(""); !matchPath(got, want) {
		t.Errorf("matchInput()\n got  %v\n want %v", got, want)
	}
}

func Test_matchInput_file(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(tmpDir, "export.enex")
	if _, err := os.Create(want); err != nil {
		t.Fatal(err)
	}

	if got := matchInput("export.enex"); !matchPath(got, want) {
		t.Errorf("matchInput()\n got  %v\n want %v", got, want)
	}
}

func Test_matchInput_dir(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(tmpDir, "test2"), 0777); err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(tmpDir, "test2", "in_dir.enex")
	if _, err := os.Create(want); err != nil {
		t.Fatal(err)
	}

	if got := matchInput("test2"); !matchPath(got, want) {
		t.Errorf("matchInput()\n got  %v\n want %v", got, want)
	}
}

func Test_matchInput_glob(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(tmpDir, "test3"), 0777); err != nil {
		t.Fatal(err)
	}

	want1 := filepath.Join(tmpDir, "test3", "glob1.enex")
	if _, err := os.Create(want1); err != nil {
		t.Fatal(err)
	}
	want2 := filepath.Join(tmpDir, "test3", "glob2.enex")
	if _, err := os.Create(want2); err != nil {
		t.Fatal(err)
	}

	if got := matchInput("test3/glob*.enex"); !matchPath(got, want1, want2) {
		t.Errorf("matchInput()\n got  %v\n want %v\n and  %v", got, want1, want2)
	}
}

func matchPath(got []string, want ...string) bool {
	for i, s := range want {
		if !strings.Contains(got[i], s) {
			return false
		}
	}

	return true
}
