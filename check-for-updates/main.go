package main

import (
  "net/http"
  "io/ioutil"
  "fmt"
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

func downloadFeed() []byte {
  rs, err := http.Get("https://gis-services.metria.se/msbfeed/skyddsrum.xml")

  if err != nil {
    panic(err)
  }
  defer rs.Body.Close()

  bodyBytes, err := ioutil.ReadAll(rs.Body)
  if err != nil {
    panic(err)
  }

//  bodyString := string(bodyBytes)

  return bodyBytes
}

func parseFeed(byteFeed []byte) {
  var parsedFeed Feed

  xml.Unmarshal(byteFeed, &parsedFeed)

  fmt.Print(parsedFeed.Entry.Link)
  fmt.Print(parsedFeed.Entry.Updated)
}

func main() {
  byteFeed := downloadFeed()

  parseFeed(byteFeed)
}

