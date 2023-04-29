package model

import (
	"github.com/nextlinux/gosbom/gosbom/file"
	"github.com/nextlinux/gosbom/gosbom/source"
)

type Secrets struct {
	Location source.Coordinates  `json:"location"`
	Secrets  []file.SearchResult `json:"secrets"`
}
