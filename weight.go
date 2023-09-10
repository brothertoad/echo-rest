package main

import (
  "database/sql"
  "fmt"
  "strings"
  "github.com/labstack/echo/v4"
  "github.com/urfave/cli/v2"
  "github.com/brothertoad/btu"
)

var weightCommand = cli.Command {
  Name: "weight",
  Usage: "seeds the weight tables",
  Action: doWeight,
}

///////////////////////////////////////////
// Logic for web
///////////////////////////////////////////

func addDailyWeight(c echo.Context, db *sql.DB) error {
  date := c.FormValue("date")
  weight := c.FormValue("weight")
  fmt.Printf("date is %v (%T) and weight is %v (%T)\n", date, date, weight, weight)
  return nil
}

///////////////////////////////////////////
// Common functions
///////////////////////////////////////////

func updateMonth(db *sql.DB, month, year int, refreshYear bool) {
  d0 := year * 10000 + month * 100  // start date
  d1 := d0 + 99  // end date
  rows, err := db.Query("select weight from weightDaily where date >= $1 and date <= $2", d0, d1)
  btu.CheckError(err)
  defer rows.Close()
  sum := 0
  count := 0
  for rows.Next() {
    var weight int
		if err := rows.Scan(&weight); err != nil {
			btu.Fatal(err.Error())
		}
    sum += weight
    count++
  }
  avg := sum / count
  // Try to update first - if that fails, then insert.
  // Use this order because updates will be much more common then inserts.
  _, err = db.Exec("update weightSum set sum = $1, count = $2, avg = $3 where year = $4 and month = $5", sum, count, avg, year, month)
  if err != nil {
    // Try insert
    _, err = db.Exec("insert into weightSum (month, year, sum, count, avg) values ($1, $2, $3, $4, $5)", month, year, sum, count, avg)
    btu.CheckError(err)
  }
  if refreshYear {
    updateYear(db, year)
  }
}

func updateYear(db *sql.DB, year int) {

}

///////////////////////////////////////////
// Logic for weight command
///////////////////////////////////////////

type daily struct {
  date int
  weight int
}

// Loads seed weights into the weightDaily and weightSum tables.
// We assume the input file is sorted.
func doWeight(c *cli.Context) error {
  // We need a file name.
  if c.Args().Len() != 1 {
    btu.Fatal("Usage: echo-rest weight file\n")
  }
  // read input file into a string, then split at new lines
  s := btu.ReadFileS(c.Args().Get(0))
  lines := strings.Split(s, "\n")
  fmt.Printf("Found %d lines.\n", len(lines))
  dailies := make([]daily, 0, len(lines))
  for _, line := range(lines) {
    if len(line) == 0 {
      continue
    }
    parts := strings.Split(line, ",")
    if len(parts) != 2 {
      continue
    }
    var d daily
    d.date = btu.Atoi(parts[0])
    d.weight = weightStringToInt(parts[1])
    dailies = append(dailies, d)
  }
  fmt.Printf("Found %d daily weights.\n", len(dailies))
  loadDailies(dailies)
  return nil
}

func loadDailies(dailies []daily) {
  db := openDB()
  defer db.Close()
  for _, d := range(dailies) {
    _, err := db.Exec("insert into weightDaily (date, weight) values ($1, $2)", d.date, d.weight)
    btu.CheckError(err)
  }
  // Need to update sums.
}

func weightStringToInt(s string) int {
  parts := strings.Split(s, ".")
  withoutDecimal := parts[0] + parts[1]
  return btu.Atoi(withoutDecimal)
}
