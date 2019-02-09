// +build lambda

package main

import (
  "context"
  "database/sql"
  "encoding/json"
  "fmt"
  "github.com/aws/aws-lambda-go/cfn"
  "github.com/aws/aws-lambda-go/lambda"
  _ "github.com/go-sql-driver/mysql"
  "os"
)

func handler(ctx context.Context, event cfn.Event) (rid string, data map[string]interface{}, err error) {
  jsonEvent, e := json.Marshal(event)
  rid = "DB_MIGRATION"
  if e != nil {
    fmt.Printf("json Marshal of event failed %v", e)
    jsonEvent = []byte{}
  }

  fmt.Printf("Starting with\ncontext %v\njsonEvent %v\n", ctx, string(jsonEvent[:]))

  if event.RequestType != cfn.RequestCreate {
    fmt.Println("Request type is not of type create, returning")
    fmt.Printf("physicalResourceID: %s\n", rid)
    return
  }
  fmt.Printf("physicalResourceID: %s\n", rid)

  fmt.Printf("Connection to database at %s.\n", os.Getenv("DB_ADDRESS"))
  db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/", os.Getenv("DB_MASTER_USER"), os.Getenv("DB_MASTER_PASSWORD"), os.Getenv("DB_ADDRESS")))
  if err != nil {
    panic(err)
  }
  defer db.Close()

  fmt.Printf("Database connection open to %s.\n", os.Getenv("DB_ADDRESS"))

  defer func() {
    if r := recover(); r != nil {
      err = fmt.Errorf("handler: Failed to migrate DB: %v", r)
    }
  }()

  MigrateDb(db)

  return
}

func main() {
  lambda.Start(cfn.LambdaWrap(handler))
}
