package dbFlags

import (
  "flag"
  "fmt"
  "os"
)

func GetDbFlags() (*string, *string, *string) {
  dbU := flag.String("db-user", "", "Database user")
  dbP := flag.String("db-password", "", "Database password")
  dbA := flag.String("db-address", "", "Database address, including port number.")
  flag.Parse()

  missReq := false
  flag.VisitAll(func(f *flag.Flag) {
    if f.Value.String() == "" && f.Name != "db-password" {
      fmt.Printf("Flag hasnt been set %s", f.Name)
      fmt.Println()
      missReq = true
    }
  })
  if missReq != false {
    os.Exit(0)
  }

  return dbU, dbP, dbA
}
