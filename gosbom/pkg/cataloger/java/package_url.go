package java

import (
	"github.com/nextlinux/packageurl-go"
	"github.com/nextlinux/syft/syft/pkg"
	"github.com/nextlinux/syft/syft/pkg/cataloger/common/cpe"
)

// PackageURL returns the PURL for the specific java package (see https://github.com/package-url/purl-spec)
func packageURL(name, version string, metadata pkg.JavaMetadata) string {
	var groupID = name
	groupIDs := cpe.GroupIDsFromJavaMetadata(metadata)
	if len(groupIDs) > 0 {
		groupID = groupIDs[0]
	}

	pURL := packageurl.NewPackageURL(
		packageurl.TypeMaven, // TODO: should we filter down by package types here?
		groupID,
		name,
		version,
		nil, // TODO: there are probably several qualifiers that can be specified here
		"")
	return pURL.ToString()
}
