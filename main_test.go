package main

import (
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
	tmpDir := tDir(t)
	input := filepath.FromSlash(tmpDir + "/export.enex")
	err := os.WriteFile(input, []byte(sampleFile), 0600)
	if err != nil {
		t.Fatalf("failed to create a test file at %s", input)
	}
	files, _ := matchInput(input)
	output := newNoteFilesDir(tmpDir, false, false)
	converter, _ := internal.NewConverter("", true, "", false)
	run(files, output, newProgressBar(false), converter)

	want := filepath.FromSlash(output.Path() + "/Test.md")
	_, err = os.Stat(want)
	if err != nil && os.IsNotExist(err) {
		t.Error("Test.md was not created")
	}
}

func Test_matchInput_cwd(t *testing.T) {
	tmpDir := tDir(t)
	want := wantFile(t, tmpDir, "test_export.enex")

	if got, _ := matchInput(""); !matchPath(got, want) {
		t.Errorf("matchInput()\n got  %v\n want %v", got, want)
	}
}

func Test_matchInput_file(t *testing.T) {
	tmpDir := tDir(t)
	want := wantFile(t, tmpDir, "test_export.enex")

	if got, _ := matchInput("test_export.enex"); !matchPath(got, want) {
		t.Errorf("matchInput()\n got  %v\n want %v", got, want)
	}
}

func Test_matchInput_dir(t *testing.T) {
	tmpDir := tDir(t)
	wDir := wantDir(t, tmpDir, "testDir")
	want := wantFile(t, wDir, "in_dir.enex")

	if got, _ := matchInput("testDir"); !matchPath(got, want) {
		t.Errorf("matchInput()\n got  %v\n want %v", got, want)
	}
}

func Test_matchInput_glob(t *testing.T) {
	tmpDir := tDir(t)
	wDir := wantDir(t, tmpDir, "testDir2")
	want1 := wantFile(t, wDir, "glob1.enex")
	want2 := wantFile(t, wDir, "glob2.enex")

	if got, _ := matchInput("testDir2/glob*.enex"); !matchPath(got, want1, want2) {
		t.Errorf("matchInput()\n got  %v\n want %v\n and  %v", got, want1, want2)
	}
}

func Test_matchInput_fail(t *testing.T) {
	_ = tDir(t)

	if _, err := matchInput("not_exist.enex"); err == nil || !strings.HasPrefix(err.Error(), "[ERROR]") {
		t.Errorf("matchInput() got unexpected eror %v", err)
	}
}

func Test_matchInput_ignoreDirectories(t *testing.T) {
	tmpDir := tDir(t)
	_ = wantDir(t, tmpDir, "test1")
	want := wantFile(t, tmpDir, "test1.enex")

	if got, _ := matchInput("test*"); !matchPath(got, want) {
		t.Errorf("matchInput()\n got  %v\n want %v", got, want)
	}
}

func tDir(t *testing.T) string {
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	return tmpDir
}

func wantFile(t *testing.T, path ...string) string {
	filePath := filepath.Join(path...)
	if _, err := os.Create(filePath); err != nil {
		t.Fatal(err)
	}
	return filePath
}

func wantDir(t *testing.T, path ...string) string {
	dirPath := filepath.Join(path...)
	if err := os.MkdirAll(filepath.Join(path...), 0777); err != nil {
		t.Fatal(err)
	}
	return dirPath
}

func matchPath(got []string, want ...string) bool {
	for i, s := range want {
		if !strings.Contains(got[i], s) {
			return false
		}
	}

	return true
}
