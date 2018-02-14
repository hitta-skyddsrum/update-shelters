package main

func main() {
  byteFeed := DownloadFile("https://gis-services.metria.se/msbfeed/skyddsrum.xml")

  ParseFeed(byteFeed)
}

