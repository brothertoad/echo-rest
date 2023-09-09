package main

import (
  "database/sql"
  "fmt"
  "strconv"
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
  w, _ := strconv.Atoi(withoutDecimal)
  return w
}
