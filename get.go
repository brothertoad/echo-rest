package main

import (
  "fmt"
  "os"
  "github.com/urfave/cli/v2"
  "github.com/brothertoad/btu"
)

var getCommand = cli.Command {
  Name: "get",
  Usage: "gets the latest value for a block",
  Action: doGet,
}

func doGet(c *cli.Context) error {
  // We need a name, and possibly a file to write to.
  if c.Args().Len() < 1 || c.Args().Len() > 2 {
    btu.Fatal("Usage: echo-rest get name <file to write to>\n")
  }
  db := openDB();
  defer db.Close()
  block := new(BlockRequest)
  block.Name = c.Args().Get(0)
  getBlock(db, block)
  if c.Args().Len() == 2 {
    // write to a file
    err := os.WriteFile(c.Args().Get(1), []byte(block.Contents), 0644)
    btu.CheckError(err)
  } else {
    fmt.Printf("%s", block.Contents)
  }
  return nil
}
