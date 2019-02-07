// +build lambda

package main

//go:generate go run include-migration.go

import (
  "database/sql"
  "context"
  "fmt"
  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  _ "github.com/go-sql-driver/mysql"
  "net/url"
  "os"
)

func handler(ctx context.Context, s3Event events.S3Event) {
  if len(s3Event.Records) > 1 {
    panic(fmt.Sprintf("Received %d records from S3 event", len(s3Event.Records)))
  }

  record := s3Event.Records[0]

  sess, err := session.NewSession()
  if err != nil {
    panic(err)
  }

  key, err := url.QueryUnescape(record.S3.Object.Key)
  if err != nil {
    panic(err)
  }

  cPath := "/tmp/" + key
  c, err := os.Create(cPath)
  if err != nil {
    panic(err)
  }
  defer c.Close()

  downloader := s3manager.NewDownloader(sess)
  _, err = downloader.Download(c,
    &s3.GetObjectInput{
      Bucket: &record.S3.Bucket.Name,
      Key:    &key,
    })
  if err != nil {
    fmt.Printf("Unable to download item %s from %s, %v", key, record.S3.Bucket.Name, err)
    panic(err)
  }
  fmt.Printf("%s downloaded.", key)
  fmt.Println()

  db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/", os.Getenv("DB_MASTER_USER"), os.Getenv("DB_MASTER_PASSWORD"), os.Getenv("DB_ADDRESS")))
  if err != nil {
    panic(err)
  }
  defer db.Close()
  fmt.Println("Database connection open to %s.", os.Getenv("DB_ADDRESS"))

  ImportCsvToMysql(db, schema, cPath)
}

func main() {
  lambda.Start(handler)
}

