package syftjson

import (
	"github.com/nextlinux/syft/internal"
	"github.com/nextlinux/syft/syft/sbom"
)

const ID sbom.FormatID = "syft-json"

func Format() sbom.Format {
	return sbom.NewFormat(
		internal.JSONSchemaVersion,
		encoder,
		decoder,
		validator,
		ID, "json", "syft",
	)
}
