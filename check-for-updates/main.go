package main

func main() {
  byteFeed := DownloadFeed("https://gis-services.metria.se/msbfeed/skyddsrum.xml")

  ParseFeed(byteFeed)
}

