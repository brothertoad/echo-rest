package main

import (
  "time"
)

type BlockRequest struct {
  Name string `json:"name" form:"name" param:"name"`
  Contents string `json:"contents" form:"contents" param:"contents"`
  ModTime time.Time `json:"modTime" form:"modTime" param:"modTime"`
}

type ListResponse struct {
  Name string `json:"name" form:"name" param:"name"`
  Items []string `json:"items" form:"items" param:"items"`
  ModTime time.Time `json:"modTime" form:"modTime" param:"modTime"`
}

type AvgWeight struct {
  Month int `json:"month"`
  Year int `json:"year"`
  Avg string `json:"avg"`
}
