package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/labstack/echo/v4"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/brothertoad/btu"
	"github.com/brothertoad/echo-rest/model"
)

func main() {
	db, err := sql.Open("pgx", os.Getenv("REST_DB_URL"))
	btu.CheckError(err)
	defer db.Close()

	e := echo.New()
	e.GET("/block/:kind/:name", func(c echo.Context) error {
		return getRawBlock(c, db)
	})
	e.POST("/block", func(c echo.Context) error {
		return postBlock(c, db)
	})
	e.Logger.Fatal(e.Start(":9903"))
}

func getRawBlock(c echo.Context, db *sql.DB) error {
	block := new(model.BlockRequest)
	if err := c.Bind(block); err != nil {
		return err
	}
	err := db.QueryRow("select contents, modTime from blocks where name = $1 and kind = $2", block.Name, block.Kind).Scan(&block.Contents, &block.ModTime)
	if err != nil {
		// If there was an error, just assume we are creating a new block, so just return an empty one.
		block.Contents = ""
		block.ModTime = time.Now()
	}
	fmt.Printf("getRawBlock: returning %+v\n", block)
	return c.JSON(http.StatusOK,  block)
}

func postBlock(c echo.Context, db *sql.DB) error {
	block := new(model.BlockRequest)
	if err := c.Bind(block); err != nil {
		return err
	}
	result, err := db.Exec("insert into blocks (name, kind, contents, modTime) values ($1, $2, $3, $4)", block.Name, block.Kind, block.Contents, time.Now())
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, fmt.Sprintf("%d", insertId))
}
