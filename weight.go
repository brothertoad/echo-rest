package main

import (
  "database/sql"
  "fmt"
  "github.com/labstack/echo/v4"
)

func addDailyWeight(c echo.Context, db *sql.DB) error {
  date := c.FormValue("date")
  weight := c.FormValue("weight")
  fmt.Printf("date is %v (%T) and weight is %v (%T)\n", date, date, weight, weight)
  return nil
}
