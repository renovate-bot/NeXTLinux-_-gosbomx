package table

import (
	"github.com/nextlinux/syft/syft/sbom"
)

const ID sbom.FormatID = "syft-table"

func Format() sbom.Format {
	return sbom.NewFormat(
		sbom.AnyVersion,
		encoder,
		nil,
		nil,
		ID, "table",
	)
}
