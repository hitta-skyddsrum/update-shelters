package main

import (
  "encoding/xml"
)

type Feed struct {
  XMLName xml.Name  `xml:"feed"`
  Entry   Entry     `xml:"entry"`
}

type Entry struct {
  XMLName xml.Name  `xml:"entry"`
  Link    string    `xml:"id"`
  Updated string    `xml:"updated"`
}

func ParseFeed(byteFeed []byte) Feed {
  var parsedFeed Feed

  xml.Unmarshal(byteFeed, &parsedFeed)

  return parsedFeed
}


