package enex

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
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
	err := NewDecoder(data).Decode(&e)

	for i := range e.Notes {
		if err := decodeContent(&e.Notes[i]); err != nil {
			// EOF is a known case when the content is empty
			if !errors.Is(err, io.EOF) {
				e.Notes = append(e.Notes[:i], e.Notes[+1:]...)
				return nil, fmt.Errorf("decoding note %s: %w", e.Notes[i].Title, err)
			}
		}

		err = decodeRecognition(&e.Notes[i])
		if err != nil {
			return nil, err
		}
	}

	return &e, err
}

type Decoder struct {
	xml *xml.Decoder
}

func NewDecoder(r io.Reader) *Decoder {
	d := xml.NewDecoder(r)
	d.Strict = false

	return &Decoder{xml: d}
}

func (d Decoder) Decode(v any) error {
	return d.xml.Decode(v)
}

type StreamDecoder struct {
	xml *xml.Decoder
}

func NewStreamDecoder(r io.Reader) (*StreamDecoder, error) {
	buf := bytes.Buffer{}
	if _, err := buf.ReadFrom(r); err != nil {
		return nil, err
	}
	clean := removeNestedCDATA(buf.String())

	d := xml.NewDecoder(strings.NewReader(clean))
	d.Strict = false

	for {
		token, err := d.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, fmt.Errorf("failed to initialise stream reader: no en-export data found: %w", err)
			}
			return nil, err
		}

		element, ok := token.(xml.StartElement)
		if ok && element.Name.Local == "en-export" {
			break
		}
	}

	return &StreamDecoder{xml: d}, nil
}

func (d StreamDecoder) Next(n *Note) error {
	for {
		token, err := d.xml.Token()
		if err != nil {
			return err
		}
		element, ok := token.(xml.StartElement)

		if ok && element.Name.Local == "note" {
			err = d.xml.DecodeElement(n, &element)
			if err != nil {
				return err
			}
			err = decodeContent(n)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}

			return decodeRecognition(n)
		}
	}
}

func decodeContent(n *Note) error {
	var c Content
	var reader = bytes.NewReader(n.Content)

	if err := NewDecoder(reader).Decode(&c); err != nil {
		return err
	}
	n.Content = c.Text
	return nil
}

func decodeRecognition(n *Note) error {
	for j := range n.Resources {
		if res := n.Resources[j]; len(res.Recognition) == 0 {
			hash := hashRe.FindString(res.Attributes.SourceUrl)
			if len(hash) > 0 {
				n.Resources[j].ID = hash
			}
			continue
		}
		var rec Recognition
		decoder := NewDecoder(bytes.NewReader(n.Resources[j].Recognition))
		err := decoder.Decode(&rec)
		if err != nil {
			return fmt.Errorf("decoding resource %s: %w", n.Resources[j].Attributes.Filename, err)
		}
		n.Resources[j].ID = rec.ObjID
		n.Resources[j].Type = rec.ObjType
	}

	return nil
}

var reCDATA = regexp.MustCompile(`<!\[CDATA\[(.*?)\]\]>`)

// removeNestedCDATA tags in the note content
//
// Nested CDATA tags are not allowed by XML specification
// but Evernote puts them anyway, causing "Unexpected EOF" errors during decoding
func removeNestedCDATA(input string) string {
	output := reCDATA.ReplaceAllStringFunc(input, func(match string) string {
		submatch := reCDATA.FindStringSubmatch(match)
		if len(submatch) > 1 {
			return submatch[1]
		}
		return match
	})

	// Recursively remove nested CDATA tags
	if output != input {
		return removeNestedCDATA(output)
	}

	return output
}
