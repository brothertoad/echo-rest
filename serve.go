package main

import (
  "database/sql"
  "fmt"
  "net/http"
  "os"
  _ "github.com/jackc/pgx/stdlib"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  "github.com/urfave/cli/v2"
  "github.com/brothertoad/btu"
)

var serveCommand = cli.Command {
  Name: "serve",
  Usage: "run as a REST service",
  Flags: []cli.Flag {
    &cli.IntFlag {Name: "port", Aliases: []string{"p"}, Value: 9903},
  },
  Action: doServe,
}

func doServe(c *cli.Context) error {
  port := c.Int("port")
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
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
  return nil
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
