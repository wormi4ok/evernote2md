package file_test

import (
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
	b, err := os.ReadFile(filePath)
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
	longName            = `ImtiredofbeingwhatyouwantmetobeFeelingsofaithlesslostunderthesurfaceDontknowwhatyoureexpectingofmePutunderthepressureofwalkinginyourshoesCaughtintheundertowjustcaughtintheundertowEverystepthatItakeisanothermistaketoyouCaughtintheundertowjustcaughtintheundertowIvebecomesonumbIcantfeelyouthereBecomesotiredsomuchmoreawareImbecomingthisallIwanttodoIsbemorelikemeandbelesslikeyouCantyouseethatyouresmotheringmeHoldingtootightlyafraidtolosecontrolCauseeverythingthatyouthoughtIwouldbeHasfallenapartrightinfrontofyouCaughtintheundertowjustcaughtintheundertowEverystepthatItakeisanothermistaketoyouCaughtintheundertowjustcaughtintheundertowAndeverysecondIwasteismorethanIcantakeIvebecomesonumbIcantfeelyouthereBecomesotiredsomuchmoreawareImbecomingthisallIwanttodoIsbemorelikemeandbelesslikeyouAndIknowImayendupfailingtooButIknowYouwerejustlikemewithsomeonedisappointedinyouIvebecomesonumbIcantfeelyouthereBecomesotiredsomuchmoreawareImbecomingthisallIwanttodoIsbemorelikemeandbelesslikeyouIvebecomesonumbIcantfeelyouthereImtiredofbeingwhatyouwantmetobeIvebecomesonumbIcantfeelyouthereImtiredofbeingwhatyouwantmetobe`
	wantLongName        = `ImtiredofbeingwhatyouwantmetobeFeelingsofaithlesslostunderthesurfaceDontknowwhatyoureexpectingofmePutunderthepressureofwalkinginyourshoesCaughtintheundertowjustcaughtintheundertowEverystepthatItakeisanothermistaketoyouCaughtintheundertowjustcaughtinthe`
	multibyteString     = `上村脱史肝感円詰人子意情涯局需日曜予楽雑割応内提健門畔和機読展民安模新越権政評図行部提政賛必能新国統本覧応測属対賀君光間断難察展任下産位断物氷与英良図数訪死変度稲与個高職車察渉沢次読康寧平展申赤購刊皿正大想球米読量減阪誠投乏止圧北馬委日江変紀京墳判賞別側稿沖商終成伝館内御交返交換臓多積社計覚大的豊冷役権校宴要援年善真走橋生閥映氷東毎止育自前動歳線定見保劇寺触室幕作書刺山設経録車応探済本候困求止長年子追年意経燃図楽南未向連横業籠結遇趣挑見健証級代卒料済株迭済変物歌会多意生式童要更事記火卒動職寺誉第康真資板遅政南済記試表戒匿軍入読的池要処情一払注野原謙湖高由題質親報質振炎計界能転火鹿午京授需国供北容期字芝性木完文案合的根間津選減野部校足石正治解木払多歩郵難故助竹廟版場開間代選映参参住反向高常然受京報左厚責点南惑中選玲転的良芸玲離自回楽情勢州左見昨応無連同世重整算潮搭区議投研止載芳断独出国兵庭真歩時雷度実営機的合瞬晋民意投日暮紙用仙世標回左浦日馬八問川止以地聞地告朝童特`
	wantMultibyteString = `上村脱史肝感円詰人子意情涯局需日曜予楽雑割応内提健門畔和機読展民安模新越権政評図行部提政賛必能新国統本覧応測属対賀君光間断難察展任下産位断物氷与英良図数訪死変度稲与個高職車察渉沢次読康寧平展申赤購刊皿正大想球米読量減阪誠投乏止圧北馬委日江変紀京墳判賞別側稿沖商終成伝館内御交返交換臓多積社計覚大的豊冷役権校宴要援年善真走橋生閥映氷東毎止育自前動歳線定見保劇寺触室幕作書刺山設経録車応探済本候困求止長年子追年意経燃図楽南未向連横業籠結遇趣挑見健証級代卒料済株迭済変物歌会多`
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
		{"blacklisted chars", "1 <3 6:014|\\|6", "1_<3_6_014|\\|6"},
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
