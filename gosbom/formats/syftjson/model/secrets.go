package model

import (
	"github.com/nextlinux/syft/syft/file"
	"github.com/nextlinux/syft/syft/source"
)

type Secrets struct {
	Location source.Coordinates  `json:"location"`
	Secrets  []file.SearchResult `json:"secrets"`
}
