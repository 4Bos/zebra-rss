package parser

import (
	"encoding/xml"
	"time"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Title   string   `xml:"title"`
	Items   []Item   `xml:"item"`
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Guid        string   `xml:"guid"`
	PubDate     PubDate  `xml:"pubDate"`
	Hash        string   `xml:"-"`
}

type PubDate time.Time

func (pd *PubDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string

	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	parsedPubDate, err := ParsePubDate(s)

	if err == nil {
		*pd = PubDate(*parsedPubDate)
		return nil
	}

	return err
}

func (pd PubDate) String() string {
	return time.Time(pd).String()
}
