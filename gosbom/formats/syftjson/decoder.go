package gosbomjson

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/nextlinux/gosbom/gosbom/formats/gosbomjson/model"
	"github.com/nextlinux/gosbom/gosbom/sbom"
)

func decoder(reader io.Reader) (*sbom.SBOM, error) {
	dec := json.NewDecoder(reader)

	var doc model.Document
	err := dec.Decode(&doc)
	if err != nil {
		return nil, fmt.Errorf("unable to decode gosbom-json: %w", err)
	}

	return toGosbomModel(doc)
}
