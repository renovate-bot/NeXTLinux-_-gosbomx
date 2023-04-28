package file

import (
	"github.com/nextlinux/stereoscope/pkg/file"
	"github.com/nextlinux/gosbom/internal/log"
	"github.com/nextlinux/gosbom/gosbom/source"
)

func allRegularFiles(resolver source.FileResolver) (locations []source.Location) {
	for location := range resolver.AllLocations() {
		resolvedLocations, err := resolver.FilesByPath(location.RealPath)
		if err != nil {
			log.Warnf("unable to resolve %+v: %+v", location, err)
			continue
		}

		for _, resolvedLocation := range resolvedLocations {
			metadata, err := resolver.FileMetadataByLocation(resolvedLocation)
			if err != nil {
				log.Warnf("unable to get metadata for %+v: %+v", location, err)
				continue
			}

			if metadata.Type != file.TypeRegular {
				continue
			}
			locations = append(locations, resolvedLocation)
		}
	}
	return locations
}
