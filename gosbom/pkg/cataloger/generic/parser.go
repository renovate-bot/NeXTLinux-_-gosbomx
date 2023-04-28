package generic

import (
	"github.com/nextlinux/syft/syft/artifact"
	"github.com/nextlinux/syft/syft/linux"
	"github.com/nextlinux/syft/syft/pkg"
	"github.com/nextlinux/syft/syft/source"
)

type Environment struct {
	LinuxRelease *linux.Release
}

type Parser func(source.FileResolver, *Environment, source.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error)
