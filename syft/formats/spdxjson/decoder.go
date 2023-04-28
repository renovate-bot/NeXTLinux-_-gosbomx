package spdxjson

import (
	"fmt"
	"io"

	"github.com/spdx/tools-golang/json"

	"github.com/nextlinux/syft/syft/formats/common/spdxhelpers"
	"github.com/nextlinux/syft/syft/sbom"
)

func decoder(reader io.Reader) (s *sbom.SBOM, err error) {
	doc, err := json.Read(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to decode spdx-json: %w", err)
	}

	return spdxhelpers.ToSyftModel(doc)
}
