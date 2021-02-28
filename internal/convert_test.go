package internal_test

import (
	"bytes"
	"encoding/base64"
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/encoding/markdown"
	"github.com/wormi4ok/evernote2md/internal"
)

var update = flag.Bool("update-golden-file", false, "Update golden file")

type testTemplate struct {
	name              string
	arg               *enex.Note
	want              *markdown.Note
	enableFrontMatter bool
	wantErr           bool
	markdownFile      string
}

const encodedImage = `R0lGODlhEAAQAPcAAAAAAAAAMwAAZgAAmQAAzAAA/wArAAArMwArZgArmQArzAAr/wBVAABVMwBVZgBVmQBVzABV/wCAAACAMwCAZgCAmQCAzACA/wCqAACqMwCqZgCqmQCqzACq/wDVAADVMwDVZgDVmQDVzADV/wD/AAD/MwD/ZgD/mQD/zAD//zMAADMAMzMAZjMAmTMAzDMA/zMrADMrMzMrZjMrmTMrzDMr/zNVADNVMzNVZjNVmTNVzDNV/zOAADOAMzOAZjOAmTOAzDOA/zOqADOqMzOqZjOqmTOqzDOq/zPVADPVMzPVZjPVmTPVzDPV/zP/ADP/MzP/ZjP/mTP/zDP//2YAAGYAM2YAZmYAmWYAzGYA/2YrAGYrM2YrZmYrmWYrzGYr/2ZVAGZVM2ZVZmZVmWZVzGZV/2aAAGaAM2aAZmaAmWaAzGaA/2aqAGaqM2aqZmaqmWaqzGaq/2bVAGbVM2bVZmbVmWbVzGbV/2b/AGb/M2b/Zmb/mWb/zGb//5kAAJkAM5kAZpkAmZkAzJkA/5krAJkrM5krZpkrmZkrzJkr/5lVAJlVM5lVZplVmZlVzJlV/5mAAJmAM5mAZpmAmZmAzJmA/5mqAJmqM5mqZpmqmZmqzJmq/5nVAJnVM5nVZpnVmZnVzJnV/5n/AJn/M5n/Zpn/mZn/zJn//8wAAMwAM8wAZswAmcwAzMwA/8wrAMwrM8wrZswrmcwrzMwr/8xVAMxVM8xVZsxVmcxVzMxV/8yAAMyAM8yAZsyAmcyAzMyA/8yqAMyqM8yqZsyqmcyqzMyq/8zVAMzVM8zVZszVmczVzMzV/8z/AMz/M8z/Zsz/mcz/zMz///8AAP8AM/8AZv8Amf8AzP8A//8rAP8rM/8rZv8rmf8rzP8r//9VAP9VM/9VZv9Vmf9VzP9V//+AAP+AM/+AZv+Amf+AzP+A//+qAP+qM/+qZv+qmf+qzP+q///VAP/VM//VZv/Vmf/VzP/V////AP//M///Zv//mf//zP///wAAAAAAAAAAAAAAACH5BAEAAPwALAAAAAAQABAAAAhuAPftg2PECEGDBQkKFFhwjZEgcIAcUWNkohGGcYAQ3GhQI8J9HxeKBAkxFEUgI0ciNBLnYsqFGoF0fLmQZcIjNBnGjJiT4sGWIGN2rAgnyEAjKBFabKkQpEIgTAtSVFPTYU6RM68K9KgVo0utAQEAOw==`

func TestConvert(t *testing.T) {
	image, _ := base64.StdEncoding.DecodeString(encodedImage)
	tests := []testTemplate{
		{
			name: "",
			arg: &enex.Note{
				Title:   "Test note",
				Content: goldenFile(t, "evernote.html"),
				Created: "20121202T112233Z",
				Updated: "20201220T223344Z",
				Tags:    []string{"tag1", "tag2"},
				Attributes: enex.NoteAttributes{
					Source:            "mobile.android",
					SourceApplication: "",
					Latitude:          "50.00000000000000",
					Longitude:         "30.00000000000000",
					Altitude:          "",
					Author:            "",
					SourceUrl:         "",
				},
				Resources: []enex.Resource{{
					ID:   "c9e6c70ea74388346ffa16ff8edbdf58",
					Mime: "image/png",
					Attributes: enex.Attributes{
						Filename: "1.jpg",
					},
					Data: enex.Data{
						Encoding: "base64",
						Content:  []byte(encodedImage),
					},
				}, {
					ID:   "90fdbde3hk91aff643883475tgh94bds1",
					Mime: "image/gif",
					Attributes: enex.Attributes{
						Filename: "1.jpg",
					},
					Data: enex.Data{
						Encoding: "base64",
						Content:  []byte(encodedImage),
					},
				}, {
					ID:   "1sdb49hgt574388346ffa19kh3edbdf09",
					Mime: "image/gif",
					Attributes: enex.Attributes{
						Filename: "complex?path=http://image.com/2.gif",
					},
					Data: enex.Data{
						Encoding: "base64",
						Content:  []byte(encodedImage),
					},
				}},
			},
			want: &markdown.Note{
				Content: []byte(""),
				CTime:   time.Date(2012, 12, 02, 11, 22, 33, 0, time.UTC),
				MTime:   time.Date(2020, 12, 20, 22, 33, 44, 0, time.UTC),
				Media: map[string]markdown.Resource{
					"c9e6c70ea74388346ffa16ff8edbdf58": {
						Name:    "1.jpg",
						Type:    "image",
						Content: image,
					},
					"90fdbde3hk91aff643883475tgh94bds1": {
						Name:    "1-1.jpg",
						Type:    "image",
						Content: image,
					},
					"1sdb49hgt574388346ffa19kh3edbdf09": {
						Name:    "complex?path=http-image-com-2.gif",
						Type:    "image",
						Content: image,
					},
				},
			},
			wantErr:           false,
			enableFrontMatter: false,
			markdownFile:      "golden.md",
		},
	}
	secondTestWithFrontMatter := tests[0]
	secondTestWithFrontMatter.markdownFile = "golden-frontmatter.md"
	secondTestWithFrontMatter.enableFrontMatter = true
	tests = append(tests, secondTestWithFrontMatter)
	for _, tt := range tests {
		c, _ := internal.NewConverter("", tt.enableFrontMatter, "", true)
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Convert(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *update {
				if err := os.WriteFile(filepath.Join("testdata", tt.markdownFile), got.Content, 0644); err == nil {
					t.SkipNow()
				}
			}
			content := goldenFile(t, tt.markdownFile)
			tt.want.Content = content
			if !reflect.DeepEqual(got, tt.want) {
				if !bytes.Equal(got.Content, tt.want.Content) {
					t.Errorf("Content mismatch! \nGot = %s, \nWant= %s", got.Content, tt.want.Content)
				} else {
					t.Errorf("Convert() = %s, want %+v", got.Media["c9e6c70ea74388346ffa16ff8edbdf58"].Content, tt.want)
				}
			}
		})
	}
}

func goldenFile(t *testing.T, filename string) []byte {
	golden := filepath.Join("testdata", filename)
	expected, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	return expected
}
