package internal_test

import (
	"bytes"
	"encoding/base64"
	"reflect"
	"testing"

	"github.com/wormi4ok/evernote2md/encoding/enex"
	"github.com/wormi4ok/evernote2md/encoding/markdown"
	"github.com/wormi4ok/evernote2md/internal"
)

const evernoteContent = `<p>abc</p>

<en-media type="image/jpeg" hash="c9e6c70ea74388346ffa16ff8edbdf58"/><en-media type="image/jpeg" hash="1sdb49hgt574388346ffa19kh3edbdf09"/>
`

const markdownContent = `# Test note

` + "`tag1`" + `

abc

![](img/1.jpg)

![](img/2.jpg)
`
const encodedImage = `R0lGODlhEAAQAPcAAAAAAAAAMwAAZgAAmQAAzAAA/wArAAArMwArZgArmQArzAAr/wBVAABVMwBVZgBVmQBVzABV/wCAAACAMwCAZgCAmQCAzACA/wCqAACqMwCqZgCqmQCqzACq/wDVAADVMwDVZgDVmQDVzADV/wD/AAD/MwD/ZgD/mQD/zAD//zMAADMAMzMAZjMAmTMAzDMA/zMrADMrMzMrZjMrmTMrzDMr/zNVADNVMzNVZjNVmTNVzDNV/zOAADOAMzOAZjOAmTOAzDOA/zOqADOqMzOqZjOqmTOqzDOq/zPVADPVMzPVZjPVmTPVzDPV/zP/ADP/MzP/ZjP/mTP/zDP//2YAAGYAM2YAZmYAmWYAzGYA/2YrAGYrM2YrZmYrmWYrzGYr/2ZVAGZVM2ZVZmZVmWZVzGZV/2aAAGaAM2aAZmaAmWaAzGaA/2aqAGaqM2aqZmaqmWaqzGaq/2bVAGbVM2bVZmbVmWbVzGbV/2b/AGb/M2b/Zmb/mWb/zGb//5kAAJkAM5kAZpkAmZkAzJkA/5krAJkrM5krZpkrmZkrzJkr/5lVAJlVM5lVZplVmZlVzJlV/5mAAJmAM5mAZpmAmZmAzJmA/5mqAJmqM5mqZpmqmZmqzJmq/5nVAJnVM5nVZpnVmZnVzJnV/5n/AJn/M5n/Zpn/mZn/zJn//8wAAMwAM8wAZswAmcwAzMwA/8wrAMwrM8wrZswrmcwrzMwr/8xVAMxVM8xVZsxVmcxVzMxV/8yAAMyAM8yAZsyAmcyAzMyA/8yqAMyqM8yqZsyqmcyqzMyq/8zVAMzVM8zVZszVmczVzMzV/8z/AMz/M8z/Zsz/mcz/zMz///8AAP8AM/8AZv8Amf8AzP8A//8rAP8rM/8rZv8rmf8rzP8r//9VAP9VM/9VZv9Vmf9VzP9V//+AAP+AM/+AZv+Amf+AzP+A//+qAP+qM/+qZv+qmf+qzP+q///VAP/VM//VZv/Vmf/VzP/V////AP//M///Zv//mf//zP///wAAAAAAAAAAAAAAACH5BAEAAPwALAAAAAAQABAAAAhuAPftg2PECEGDBQkKFFhwjZEgcIAcUWNkohGGcYAQ3GhQI8J9HxeKBAkxFEUgI0ciNBLnYsqFGoF0fLmQZcIjNBnGjJiT4sGWIGN2rAgnyEAjKBFabKkQpEIgTAtSVFPTYU6RM68K9KgVo0utAQEAOw==`

func TestConvert(t *testing.T) {
	image, _ := base64.StdEncoding.DecodeString(encodedImage)
	tests := []struct {
		name    string
		arg     *enex.Note
		want    *markdown.Note
		wantErr bool
	}{
		{
			name: "",
			arg: &enex.Note{
				Title:   "Test note",
				Content: []byte(evernoteContent),
				Tags:    []string{"tag1"},
				Resources: []enex.Resource{{
					ID: "c9e6c70ea74388346ffa16ff8edbdf58",
					Attributes: enex.Attributes{
						Filename: "1.jpg",
					},
					Data: enex.Data{
						Encoding: "base64",
						Content:  []byte(encodedImage),
					},
				}, {
					ID: "1sdb49hgt574388346ffa19kh3edbdf09",
					Attributes: enex.Attributes{
						Filename: "2.jpg",
					},
					Data: enex.Data{
						Encoding: "base64",
						Content:  []byte(encodedImage),
					},
				}},
			},
			want: &markdown.Note{
				Content: []byte(markdownContent),
				Media: map[string]markdown.Resource{
					"c9e6c70ea74388346ffa16ff8edbdf58": markdown.Resource{
						Name:    "1.jpg",
						Content: image,
					},
					"1sdb49hgt574388346ffa19kh3edbdf09": markdown.Resource{
						Name:    "2.jpg",
						Content: image,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := internal.Converter{AssetsDir: "img"}.Convert(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
