package main

import (
  "fmt"
)

func main() {
  byteFeed := DownloadFile("https://gis-services.metria.se/msbfeed/skyddsrum.xml")

  var feed = ParseFeed(byteFeed)

  var shapefile = DownloadFile(feed.Entry.Link)
  path := StoreShapefile(shapefile)

  fmt.Printf("Stored shapefile at %s\n", path)
}

