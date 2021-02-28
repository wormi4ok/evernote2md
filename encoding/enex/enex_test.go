package enex_test

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"os"
	"reflect"
	"testing"

	"github.com/wormi4ok/evernote2md/encoding/enex"
)

var expect = &enex.Export{
	XMLName: xml.Name{
		Space: "",
		Local: "en-export",
	},
	Date: "20090101T202020Z",
	Notes: []enex.Note{{
		XMLName: xml.Name{
			Space: "",
			Local: "note",
		},
		Title:   "Sample note",
		Content: []byte(`<div>text in the note<br/><b>bold text</b><br/></div><en-media type="image/jpeg" hash="09dde741f3b38c1a954358172cad4c06"/>`),
		Updated: "20090101T050505Z",
		Created: "20090101T101010Z",
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
			ID:   "09dde741f3b38c1a954358172cad4c06",
			Type: "image",
			Data: enex.Data{
				XMLName: xml.Name{
					Space: "",
					Local: "data",
				},
				Encoding: "base64",
				Content:  readFile("testdata/img.gif"),
			},
			Mime:   "image/gif",
			Width:  16,
			Height: 16,
			Attributes: enex.Attributes{
				Timestamp: "20120515T051032Z",
				Filename:  "1.jpg",
				SourceUrl: "",
			},
			Recognition: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE recoIndex PUBLIC "SYSTEM" "http://xml.evernote.com/pub/recoIndex.dtd"><recoIndex docType="unknown" objType="image" objID="09dde741f3b38c1a954358172cad4c06" engineVersion="5.5.20.1" recoType="service" lang="en" objWidth="16" objHeight="16"/>
`),
		}},
	}},
}

func TestDecode(t *testing.T) {
	enexContent, err := os.Open("testdata/export.enex")
	if err != nil {
		t.Error(err)
	}
	got, err := enex.Decode(enexContent)
	if err != nil {
		t.Errorf("Error while Decodeing = %v", err)
	}

	if !reflect.DeepEqual(got, expect) {
		t.Errorf("Decode() = %+v,\nwant %+v", got, expect)
	}

}

func readFile(filename string) []byte {
	file, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &b)
	_, err = encoder.Write(file)
	if err != nil {
		panic(err)
	}
	return append(b.Bytes(), []byte(`Ow==`)...)
}
