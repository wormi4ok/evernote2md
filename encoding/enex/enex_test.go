package enex_test

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"io"
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

var expectHash = "084f886210557e19eafc72449154331e"

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

func TestDecodeEmptyNote(t *testing.T) {
	enexContent, err := os.Open("testdata/empty.enex")
	if err != nil {
		t.Error(err)
	}
	_, err = enex.Decode(enexContent)
	if err != nil {
		t.Errorf("Error while Decoding = %v", err)
	}
}

func TestDecodeWithMissingRecognition(t *testing.T) {
	enexContent, err := os.Open("testdata/missing_recognition.enex")
	if err != nil {
		t.Error(err)
	}
	got, err := enex.Decode(enexContent)
	if err != nil {
		t.Errorf("Error while Decodeing = %v", err)
	}

	id := got.Notes[0].Resources[0].ID
	if id != expectHash {
		t.Errorf("Decoded resource id = %s,\nexpected resource id = %s", id, expectHash)
	}
}

func TestStreamDecoder(t *testing.T) {
	enexContent, err := os.Open("testdata/export.enex")
	if err != nil {
		t.Fatal(err)
	}
	d, err := enex.NewStreamDecoder(enexContent)
	if err != nil {
		t.Fatal(err)
	}
	var got enex.Note
	err = d.Next(&got)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(got, expect.Notes[0]) {
		t.Errorf("Next() = %+v,\nwant %+v", got, expect.Notes[0])
	}
	err = d.Next(&got)
	if err != io.EOF {
		t.Errorf("Next() second call = %+v,\nwant %+v", got, io.EOF)
	}
}

func TestStreamDecodeEmptyNote(t *testing.T) {
	enexContent, err := os.Open("testdata/empty.enex")
	if err != nil {
		t.Fatal(err)
	}
	d, err := enex.NewStreamDecoder(enexContent)
	if err != nil {
		t.Errorf("Error while Decoding = %v", err)
	}
	var got enex.Note
	err = d.Next(&got)
	if err != nil {
		t.Error(err)
	}
}

func TestStreamDecodeWrongFile(t *testing.T) {
	fakeFile := bytes.NewReader([]byte("Not an XML file"))
	_, err := enex.NewStreamDecoder(fakeFile)
	if err == nil {
		t.Errorf("Expected error, got = %v", err)
	}
}
func TestStreamDecodeAutofixCDATA(t *testing.T) {
	enexContent, err := os.Open("testdata/cdata.issue.enex")
	if err != nil {
		t.Fatal(err)
	}
	d, err := enex.NewStreamDecoder(enexContent)
	if err != nil {

		t.Errorf("Error while Decoding = %v", err)
	}
	var got enex.Note
	err = d.Next(&got)
	if err != nil {
		t.Error(err)
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
