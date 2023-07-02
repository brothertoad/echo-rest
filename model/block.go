package model

import (
  "time"
)

type BlockRequest struct {
  Name string `json:"name" form:"name" param:"name"`
  Kind string `json:"kind" form:"kind" param:"kind"`
  Contents string `json:"contents" form:"contents" param:"contents"`
  ModTime time.Time `json:"modTime" form:"modTime" param:"modTime"`
}
