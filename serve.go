package main

import (
  "database/sql"
  "os"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  "github.com/brothertoad/btu"
)

func doServe() {
  db, err := sql.Open("pgx", os.Getenv("REST_DB_URL"))
	btu.CheckError(err)
	defer db.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS()) // allow all requests
	e.GET("/block/:name", func(c echo.Context) error {
		return getBlock(c, db)
	})
	e.POST("/block", func(c echo.Context) error {
		return postBlock(c, db)
	})
	// Should specify port in some kind of configuration.
	e.Logger.Fatal(e.Start(":9903"))
  // return nil
}
