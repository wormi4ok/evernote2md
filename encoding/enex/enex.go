package enex

import (
	"bytes"
	"encoding/xml"
	"io"
	"log"
)

type (
	// Export represents Evernote enex file structure
	Export struct {
		XMLName xml.Name `xml:"en-export"`
		Date    string   `xml:"export-date,attr"`
		Notes   []Note   `xml:"note"`
	}

	// Note is one note in Evernote
	Note struct {
		XMLName    xml.Name   `xml:"note"`
		Title      string     `xml:"title"`
		Content    []byte     `xml:"content"`
		Updated    string     `xml:"created"`
		Created    string     `xml:"updated"`
		Tags       []string   `xml:"tag"`
		Attributes string     `xml:"note-attributes"`
		Resources  []Resource `xml:"resource"`
	}

	// Resource embedded in the note
	Resource struct {
		ID          string
		Type        string
		Data        Data       `xml:"data"`
		Mime        string     `xml:"mime"`
		Width       int        `xml:"width"`
		Height      int        `xml:"height"`
		Attributes  Attributes `xml:"resource-attributes"`
		Recognition []byte     `xml:"recognition"`
	}

	// Attributes of the resource
	Attributes struct {
		Timestamp string `xml:"timestamp"`
		Filename  string `xml:"file-name"`
	}
	// Recognition for the resource
	Recognition struct {
		XMLName xml.Name `xml:"recoIndex"`
		ObjID   string   `xml:"objID,attr"`
		ObjType string   `xml:"objType,attr"`
	}

	// Data object in base64
	Data struct {
		XMLName  xml.Name `xml:"data"`
		Encoding string   `xml:"encoding,attr"`
		Content  []byte   `xml:",innerxml"`
	}

	// Content of Evernote Notes
	Content struct {
		Text []byte `xml:",innerxml"`
	}
)

// Decode will return an Export from evernote
func Decode(data io.Reader) (*Export, error) {
	var e Export

	decoder := xml.NewDecoder(data)
	decoder.Strict = false
	err := decoder.Decode(&e)

	for i := range e.Notes {
		var c Content
		decoder := xml.NewDecoder(bytes.NewReader(e.Notes[i].Content))
		decoder.Strict = false
		err = decoder.Decode(&c)
		if err != nil {
			log.Fatal(err)
		}
		e.Notes[i].Content = c.Text

		for j := range e.Notes[i].Resources {
			var r Recognition
			if len(e.Notes[i].Resources[j].Recognition) == 0 {
				continue
			}
			decoder := xml.NewDecoder(bytes.NewReader(e.Notes[i].Resources[j].Recognition))
			decoder.Strict = false
			err = decoder.Decode(&r)
			if err != nil {
				log.Fatal(err)
			}
			e.Notes[i].Resources[j].ID = r.ObjID
			e.Notes[i].Resources[j].Type = r.ObjType
		}

	}
	return &e, err
}
