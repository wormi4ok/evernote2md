package file_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wormi4ok/evernote2md/file"
)

const dirName = "testdata"
const fileName = "testfile"
const content = "testdataInsideFile"

func TestSave(t *testing.T) {
	err := file.Save(dirName, fileName, strings.NewReader(content))
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dirName)

	// Directory should be created
	_, err = os.Stat(dirName)
	if err != nil {
		t.Error("directory was not created")
	}

	//File should be created
	filePath := filepath.FromSlash(dirName + "/" + fileName)
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Error(err)
	}
	// Content in the file should match original
	if string(b) != content {
		t.Errorf("Want content = %v, got = %v", content, string(b))
	}
}

func TestSaveWithEmptyName(t *testing.T) {
	err := file.Save(dirName, "", strings.NewReader(content))
	if err != nil {
		t.Error("Should skip without error")
	}
}

const longName = `SeemerideoutofthesunsetOnyourcolorTVscreenOutofallthatIcangetIfyouknowwhatImeanWomentotheleftofmeAndwomentotherightAintgotnogunAintgotnoknifeDontyoustartnofightCauseImTNTImdynamiteTNTandIllwinthefightTNTImapowerloadTNTwatchmeexplode`

func TestBaseName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"valid input should return the same", "input", "input"},
		{"values with separator", "ac/dc", "ac-dc"},
		{"multiple separators in the input", "http://s.petrashov.ru", "http-s-petrashov-ru"},
		{"blacklisted chars", "1 <3 6014|\\|6", "1-3_6014_\\_6"},
		{"complicated case", "/.-./.-./.com   ", "-com"},
		{"very long names", longName, longName[:200]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := file.BaseName(tt.input); got != tt.want {
				t.Errorf("BaseName for %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
