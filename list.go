package main

import (
  "fmt"
  "github.com/urfave/cli/v2"
  "github.com/brothertoad/btu"
)

var listCommand = cli.Command {
  Name: "list",
  Usage: "lists the block names",
  Action: doList,
}

func doList(c *cli.Context) error {
  // We need a name, and possibly a file to write to.
  if c.Args().Len() != 0 {
    btu.Fatal("Usage: echo-rest list\n")
  }
  db := openDB();
  defer db.Close()
  names := listBlockNames(db)
  for _, name := range(names) {
    fmt.Printf("%s\n", name)
  }
  return nil
}
