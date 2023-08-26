package main

import (
  "github.com/urfave/cli/v2"
  "github.com/brothertoad/btu"
)

var putCommand = cli.Command {
  Name: "put",
  Usage: "puts a new value for a block",
  Action: doPut,
}

func doPut(c *cli.Context) error {
  // We need a name and a file to read from.
  if c.Args().Len() != 2 {
    btu.Fatal("Usage: echo-rest put name file\n")
  }
  db := openDB();
  defer db.Close()
  block := new(BlockRequest)
  block.Name = c.Args().Get(0)
  block.Contents = btu.ReadFileS(c.Args().Get(1))
  err := putBlock(db, block)
  btu.CheckError(err)
  return nil
}
