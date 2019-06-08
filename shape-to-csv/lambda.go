// +build lambda

package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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

	fmt.Printf("%v", record.S3.Object)

	zPath := "/tmp/shelters.zip"
	z, err := os.Create(zPath)
	if err != nil {
		panic(err)
	}
	defer z.Close()

	key, err := url.QueryUnescape(record.S3.Object.Key)
	if err != nil {
		panic(err)
	}

	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(z,
		&s3.GetObjectInput{
			Bucket: &record.S3.Bucket.Name,
			Key:    &key,
		})
	if err != nil {
		fmt.Printf("Unable to download item %s from %s, %v", key, record.S3.Bucket.Name, err)
		panic(err)
	}

	csvPath := "/tmp/shelters.csv"
	ExportShapeToCSV(zPath, csvPath)

	csvFile, err := os.Open(csvPath)
	if err != nil {
		panic(err)
	}

	uploader := s3manager.NewUploader(sess)

	csvName := strings.TrimSuffix(key, filepath.Ext(key)) + ".csv"
	upResp, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: &record.S3.Bucket.Name,
		Key:    &csvName,
		Body:   csvFile,
	})
	if err != nil {
		fmt.Printf("Unable to upload file %s to bucket %s, %v", csvName, record.S3.Bucket.Name, err)
		panic(err)
	}
	fmt.Printf("%v", upResp)
}

func main() {
	lambda.Start(handler)
}
