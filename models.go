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

type Link struct {
  Text string `json:"text"`
  URL string `json:"url"`
}

type LinkListResponse struct {
  Name string `json:"name" form:"name" param:"name"`
  Link []Link `json:"links" form:"links" param:"links"`
  ModTime time.Time `json:"modTime" form:"modTime" param:"modTime"`
}

type DailyWeight struct {
  Day int `json:"day"`
  Month int `json:"month"`
  Year int `json:"year"`
  Weight string `json:"weight"`
}

type AvgWeight struct {
  Month int `json:"month"`
  Year int `json:"year"`
  Avg string `json:"avg"`
}
