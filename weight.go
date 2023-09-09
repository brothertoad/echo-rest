package main

import (
  "database/sql"
  "fmt"
  "github.com/labstack/echo/v4"
  "github.com/urfave/cli/v2"
  "github.com/brothertoad/btu"
)

var weightCommand = cli.Command {
  Name: "weight",
  Usage: "seeds the weight tables",
  Action: doWeight,
}

func addDailyWeight(c echo.Context, db *sql.DB) error {
  date := c.FormValue("date")
  weight := c.FormValue("weight")
  fmt.Printf("date is %v (%T) and weight is %v (%T)\n", date, date, weight, weight)
  return nil
}

func doWeight(c *cli.Context) error {
  // We need a file name.
  if c.Args().Len() != 1 {
    btu.Fatal("Usage: echo-rest weight file\n")
  }
  return nil
}
