package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	"github.com/urfave/cli/v2"
	"github.com/brothertoad/btu"
)

func main() {
	app := &cli.App {
    Name: "echo-rest",
    Compiled: time.Now(),
    Usage: "a REST service, also used to manually get and put",
		// TASK: Add log-level flag
    Commands: []*cli.Command {
			&getCommand,
      &serveCommand,
    },
    // Before: command.Init,
  }
  app.Run(os.Args)
}

func openDB() *sql.DB {
	db, err := sql.Open("pgx", os.Getenv("REST_DB_URL"))
	btu.CheckError(err)
	return db
}

func getBlock(db *sql.DB, block *BlockRequest) error {
	err := db.QueryRow("select contents, modTime from blocks where name = $1 order by modTime desc limit 1", block.Name).Scan(&block.Contents, &block.ModTime)
	if err != nil {
		// If there was an error, assume we are creating a new block, so just return an empty one.
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
