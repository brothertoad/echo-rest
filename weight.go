package main

import (
  "database/sql"
  "fmt"
  "net/http"
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

// from browser: date is 2023-09-08 (string) and weight is 276.0 (string)

func addDailyWeight(c echo.Context, db *sql.DB) error {
  dateFromForm := c.FormValue("date")
  weightFromForm := c.FormValue("weight")
  // fmt.Printf("date is %v (%T) and weight is %v (%T)\n", dateFromForm, dateFromForm, weightFromForm, weightFromForm)
  year, month, day := parseDateString(dateFromForm)
  weight, err := parseWeightString(weightFromForm)
  if err != nil {
    return c.String(http.StatusBadRequest, err.Error())
  }
  date := year * 10000 + month * 100 + day
  // Try to insert first.  If that fails, try to update.
  _, err = db.Exec("insert into weightDaily (date, weight) values ($1, $2)", date, weight)
  if err != nil {
    // fmt.Printf("Err from insert is %v\n", err)
    _, err := db.Exec("update weightDaily set weight = $1 where date = $2", weight, date)
    btu.CheckError(err)
  }
  updateMonth(db, month, year, true)
  return c.String(http.StatusOK, "")
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
  total := 0
  count := 0
  for rows.Next() {
    var weight int
		if err := rows.Scan(&weight); err != nil {
			btu.Fatal(err.Error())
		}
    total += weight
    count++
  }
  avg := total / count
  // Try to update first - if that fails, then insert.
  // Use this order because updates will be much more common then inserts.
  result, err := db.Exec("update weightSum set total = $1, count = $2, avg = $3 where year = $4 and month = $5", total, count, avg, year, month)
  // If there was no existing row for this year/month, then err will be nil, but the row count will be zero.
  btu.CheckError(err)
  rowsAffected, _ := result.RowsAffected()
  if rowsAffected == 0 {
    // Try insert
    _, err = db.Exec("insert into weightSum (month, year, total, count, avg) values ($1, $2, $3, $4, $5)", month, year, total, count, avg)
    btu.CheckError(err)
  }
  if refreshYear {
    updateYear(db, year)
  }
}

func updateYear(db *sql.DB, year int) {
  rows, err := db.Query("select total, count from weightSum where year = $1 and month > 0", year)
  btu.CheckError(err)
  defer rows.Close()
  cumulativeSum := 0
  total := 0
  for rows.Next() {
    var sum, count int
		if err := rows.Scan(&sum, &count); err != nil {
			btu.Fatal(err.Error())
		}
    cumulativeSum += sum
    total += count
  }
  avg := cumulativeSum / total
  // Try to update first - if that fails, then insert.
  // Use this order because updates will be much more common then inserts.
  result, err := db.Exec("update weightSum set total = $1, count = $2, avg = $3 where year = $4 and month = 0", cumulativeSum, total, avg, year)
  // If there was no existing row for this year, then err will be nil, but the row count will be zero.
  btu.CheckError(err)
  rowsAffected, _ := result.RowsAffected()
  if rowsAffected == 0 {
    // Try insert
    _, err = db.Exec("insert into weightSum (month, year, total, count, avg) values ($1, $2, $3, $4, $5)", 0, year, cumulativeSum, total, avg)
    btu.CheckError(err)
  }
}

// String is in format yyyy-mm-dd
func parseDateString(s string) (int, int, int) {
  year := btu.Atoi(s[0:4])
  month := btu.Atoi(s[5:7])
  day := btu.Atoi(s[8:10])
  return year, month, day
}

func parseDate(d int) (int, int, int) {
  year := d / 10000
  month := (d % 10000) / 100
  day := d % 100
  return year, month, day
}

// String is in format nnn.n.  Decimal point and subsequent digit are optional,
// but if a decimal point is present, there must be mo more than one digit
// following it.
func parseWeightString(s string) (int, error) {
  dpCount := 0
  for _, ch := range(s) {
    if (ch < '0' || ch > '9') && ch != '.' {
      return 0, fmt.Errorf("weight %s is not a number\n", s)
    }
    if ch == '.' {
      dpCount++;
    }
  }
  if dpCount > 1 {
    return 0, fmt.Errorf("More than one decimal point in %s\n", s)
  }
  // Note that from here on, we know that the string has only digits and
  // (possibly) a decimal point, so we can use len to get the length of the
  // string.
  // If the string has no decimal point, just multiply by 10.
  if dpCount == 0 {
    return btu.Atoi(s) * 10, nil
  }
  n := strings.Index(s, ".")
  // If the decimal point is at the end, just ignore it and mulitply the rest of the value by 10.
  if n == (len(s) - 1) {
    return btu.Atoi(s[0:n]) * 10, nil
  }
  // If we have more than one digit after the decimal point, we have a bad value.
  if n < (len(s) - 2) {
    return 0, fmt.Errorf("Too many digits after decimal point in %s\n", s)
  }
  // OK, we have exactly one digit after the decimal point.  We're good to go.
  return btu.Atoi(s[0:n] + s[(n+1):(n+2)]), nil
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
  // Assume most recent daily is first.  Update months, then years.
  y0, m0, _ := parseDate(dailies[len(dailies)-1].date)
  y1, m1, _ := parseDate(dailies[0].date)
  year := y0
  month := m0
  for {
    updateMonth(db, month, year, false)
    if year == y1 && month == m1 {
      break
    }
    year, month = incrementMonth(year, month)
  }
  for year = y0; year <= y1; year++ {
    updateYear(db, year)
  }
}

func incrementMonth(year, month int) (int, int) {
  if month == 12 {
    return year + 1, 1
  }
  return year, month + 1
}

func weightStringToInt(s string) int {
  parts := strings.Split(s, ".")
  withoutDecimal := parts[0] + parts[1]
  return btu.Atoi(withoutDecimal)
}
