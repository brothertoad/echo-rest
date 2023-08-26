package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"github.com/labstack/echo/v4"
)

func main() {
	doServe()
}

func getBlock(c echo.Context, db *sql.DB) error {
	block := new(BlockRequest)
	if err := c.Bind(block); err != nil {
		return err
	}
	err := db.QueryRow("select contents, modTime from blocks where name = $1 order by modTime desc limit 1", block.Name).Scan(&block.Contents, &block.ModTime)
	if err != nil {
		// If there was an error, just assume we are creating a new block, so just return an empty one.
		block.Contents = ""
		block.ModTime = time.Now()
	}
	fmt.Printf("getBlock: returning %+v\n", block)
	return c.JSON(http.StatusOK,  block)
}

func postBlock(c echo.Context, db *sql.DB) error {
	block := new(BlockRequest)
	if err := c.Bind(block); err != nil {
		return err
	}
	result, err := db.Exec("insert into blocks (name, contents, modTime) values ($1, $2, $3)", block.Name, block.Contents, time.Now())
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	numRows, err := result.RowsAffected()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if numRows != 1 {
		// should not happen
	}
	fmt.Printf("postBlock: status OK\n")
	return c.String(http.StatusOK, "")
}
