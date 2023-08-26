package main

import (
	"database/sql"
	"fmt"
	"time"
)

func main() {
	doServe()
}

func getBlock(db *sql.DB, block *BlockRequest) error {
	err := db.QueryRow("select contents, modTime from blocks where name = $1 order by modTime desc limit 1", block.Name).Scan(&block.Contents, &block.ModTime)
	if err != nil {
		// If there was an error, just assume we are creating a new block, so just return an empty one.
		block.Contents = ""
		block.ModTime = time.Now()
	}
	return nil
}

func putBlock(db *sql.DB, block *BlockRequest) error {
	result, err := db.Exec("insert into blocks (name, contents, modTime) values ($1, $2, $3)", block.Name, block.Contents, time.Now())
	if err != nil {
		return err
	}
	numRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// should not happen, but we'll check for it anyway
	if numRows != 1 {
		return fmt.Errorf("More than one row (%d) affected by insert\n", numRows)
	}
	return nil
}
