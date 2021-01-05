package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
	run(input, output, newProgressBar(false), internal.Converter{})

	want := filepath.FromSlash(output.Path() + "/Test.md")
	_, err = os.Stat(want)
	if err != nil && os.IsNotExist(err) {
		t.Error("Test.md was not created")
	}
}
