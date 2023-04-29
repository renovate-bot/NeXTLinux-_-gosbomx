package gosbomjson

import (
	"github.com/nextlinux/gosbom/internal"
	"github.com/nextlinux/gosbom/gosbom/sbom"
)

const ID sbom.FormatID = "gosbom-json"

func Format() sbom.Format {
	return sbom.NewFormat(
		internal.JSONSchemaVersion,
		encoder,
		decoder,
		validator,
		ID, "json", "gosbom",
	)
}
