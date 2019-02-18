package main

import (
  "database/sql"
  "strings"
  "time"
  "fmt"
  "strconv"
  "github.com/jonas-p/go-shp"
  _ "github.com/go-sql-driver/mysql"
  "github.com/JamesStewy/go-mysqldump"
)

func createTable(db *sql.DB, fields []shp.Field) {
  strFieldNames := make([]string, 0)

  strFieldNames = append(strFieldNames, "position_lat")
  strFieldNames = append(strFieldNames, "position_long")
  for _, f := range fields {
    strFieldNames = append(strFieldNames, f.String())
  }
  strFieldNames = append(strFieldNames, "")

  colNms := strings.TrimSuffix(strings.Join(strFieldNames, " varchar(255), "), ", ");
  _, err := db.Exec("CREATE TABLE shelters (" + colNms + ")")

  if err != nil {
    panic(err)
  }
}

func getPrepareStatement(db *sql.DB, length int) *sql.Stmt {
  s := make([]string, length)

  stmtIns, err := db.Prepare("INSERT INTO shelters VALUES(?" + strings.Join(s, ", ?") + ")")

  if err != nil {
    fmt.Println(s)
    panic(err)
  }

  return stmtIns
}

func dumpDb(db *sql.DB) {
	dumper, err := mysqldump.Register(db, "dumps", time.ANSIC)
	if err != nil {
		fmt.Println("Error registering databse:", err)
		return
	}

	_, err = dumper.Dump()
	if err != nil {
		fmt.Println("Error dumping:", err)
		return
	}

	dumper.Close()
}

func ShapeToSql(zipShape *shp.ZipReader) int {
  db, err := sql.Open("mysql", "root:hitta@tcp(mysql:3306)/hitta_skyddsrum")
  if err != nil {
    panic(err)
  }
  defer db.Close()

  nrShapes := 0
  fields := zipShape.Fields()

  createTable(db, fields)

  for zipShape.Next() {
    _, shape := zipShape.Shape()

    if shape.BBox().MinX != shape.BBox().MaxX {
      panic("Shelter BBox data validation failed")
    }

    if shape.BBox().MinY != shape.BBox().MaxY {
      panic("Shelter BBox data validation failed")
    }

    vals := make([]interface{}, 0)

    coors := Sweref99ToLatLon([2]float64{shape.BBox().MinX, shape.BBox().MinY})
    vals = append(vals, strconv.FormatFloat(coors[0], 'f', -1, 64))
    vals = append(vals, strconv.FormatFloat(coors[1], 'f', -1, 64))

    for k := range fields {
      val := zipShape.Attribute(k)
      vals = append(vals, val)
    }

    stmtIns := getPrepareStatement(db, len(vals))
    _, err := stmtIns.Exec(vals...)
    stmtIns.Close()

    if err != nil {
      panic(err)
    }

    nrShapes++
  }

  dumpDb(db)

  return nrShapes
}
