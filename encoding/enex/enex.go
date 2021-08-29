package enex

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"regexp"
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
		XMLName    xml.Name       `xml:"note"`
		Title      string         `xml:"title"`
		Content    []byte         `xml:"content"`
		Updated    string         `xml:"updated"`
		Created    string         `xml:"created"`
		Tags       []string       `xml:"tag"`
		Attributes NoteAttributes `xml:"note-attributes"`
		Resources  []Resource     `xml:"resource"`
	}

	// NoteAttributes contain the note metadata
	NoteAttributes struct {
		Source            string `xml:"source"`
		SourceApplication string `xml:"source-application"`
		Latitude          string `xml:"latitude"`
		Longitude         string `xml:"longitude"`
		Altitude          string `xml:"altitude"`
		Author            string `xml:"author"`
		SourceUrl         string `xml:"source-url"`
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
		SourceUrl string `xml:"source-url"`
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

var hashRe = regexp.MustCompile(`\b[0-9a-f]{32}\b`)

// Decode will return an Export from evernote
func Decode(data io.Reader) (*Export, error) {
	var e Export
	err := newDecoder(data).Decode(&e)

	for i := range e.Notes {
		var c Content
		var reader = bytes.NewReader(e.Notes[i].Content)

		if err := newDecoder(reader).Decode(&c); err != nil {
			// EOF is a known case when the content is empty
			if !errors.Is(err, io.EOF) {
				return nil, fmt.Errorf("decoding note %s: %w", e.Notes[i].Title, err)
			}
		}
		e.Notes[i].Content = c.Text

		for j := range e.Notes[i].Resources {
			if res := e.Notes[i].Resources[j]; len(res.Recognition) == 0 {
				hash := hashRe.FindString(res.Attributes.SourceUrl)
				if len(hash) > 0 {
					e.Notes[i].Resources[j].ID = hash
				}
				continue
			}
			var rec Recognition
			decoder := newDecoder(bytes.NewReader(e.Notes[i].Resources[j].Recognition))
			err = decoder.Decode(&rec)
			if err != nil {
				return nil, fmt.Errorf("decoding resource %s: %w", e.Notes[i].Resources[j].Attributes.Filename, err)
			}
			e.Notes[i].Resources[j].ID = rec.ObjID
			e.Notes[i].Resources[j].Type = rec.ObjType
		}
	}

	return &e, err
}

func newDecoder(r io.Reader) *xml.Decoder {
	d := xml.NewDecoder(r)
	d.Strict = false
	return d
}
