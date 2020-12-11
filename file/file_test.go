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

	// File should be created
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

const (
	longName            = `SeemerideoutofthesunsetOnyourcolorTVscreenOutofallthatIcangetIfyouknowwhatImeanWomentotheleftofmeAndwomentotherightAintgotnogunAintgotnoknifeDontyoustartnofightCauseImTNTImdynamiteTNTandIllwinthefightTNTImapowerloadTNTwatchmeexplodeImdirtymeanandmightyuncleanImawantedmanPublicenemynumberoneUnderstandSolockupyourdaughterLockupyourwifeLockupyourbackdoorAndrunforyourlifeThemanisbackintownSodontyoumessmeround`
	wantLongName        = `SeemerideoutofthesunsetOnyourcolorTVscreenOutofallthatIcangetIfyouknowwhatImeanWomentotheleftofmeAndwomentotherightAintgotn`
	multibyteString     = `能影岩界手月体載髪種実献旅約客法。長式補谷人億娘民襲分続三指造声付中配。在紛意学予禁底決喜報都情漁住止歳出。断読面法狙芸現東学博必官代身限社。食念伝公並民京勧情聞個約。少日政税省聞型員易身意作。住度勤里総明証神聞取権育意燃米別意解図同。立訴技森会売多転料仕著不推。世格安円中護裁伴止売前現趣平無`
	wantMultibyteString = `能影岩界手月体載髪種実献旅約客法。長式補谷人億娘民襲分続三指造声付中配。在紛意学予`
)

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
		{"very long names", longName, wantLongName},
		{"multibyte encoding", multibyteString, wantMultibyteString},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := file.BaseName(tt.input); got != tt.want {
				t.Errorf("BaseName for %s\ngot  = %v\nwant = %v", tt.name, got, tt.want)
			}
		})
	}
}
