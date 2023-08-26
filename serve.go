package main

import (
  "database/sql"
  "net/http"
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
		return getBlockForREST(c, db)
	})
	e.POST("/block", func(c echo.Context) error {
		return putBlockForREST(c, db)
	})
	// Should specify port in some kind of configuration.
	e.Logger.Fatal(e.Start(":9903"))
  // return nil
}

func getBlockForREST(c echo.Context, db *sql.DB) error {
	block := new(BlockRequest)
	if err := c.Bind(block); err != nil {
		return err
	}
  getBlock(db, block)
	return c.JSON(http.StatusOK, block)
}

func putBlockForREST(c echo.Context, db *sql.DB) error {
	block := new(BlockRequest)
	if err := c.Bind(block); err != nil {
		return err
	}
  err := putBlock(db, block)
  if err != nil {
    return c.String(http.StatusInternalServerError, err.Error())
  }
	return c.String(http.StatusOK, "")
}
