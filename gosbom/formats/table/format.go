package table

import (
	"github.com/nextlinux/gosbom/gosbom/sbom"
)

const ID sbom.FormatID = "gosbom-table"

func Format() sbom.Format {
	return sbom.NewFormat(
		sbom.AnyVersion,
		encoder,
		nil,
		nil,
		ID, "table",
	)
}
