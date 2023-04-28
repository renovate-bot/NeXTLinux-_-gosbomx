package generic

import (
	"github.com/nextlinux/gosbom/gosbom/artifact"
	"github.com/nextlinux/gosbom/gosbom/linux"
	"github.com/nextlinux/gosbom/gosbom/pkg"
	"github.com/nextlinux/gosbom/gosbom/source"
)

type Environment struct {
	LinuxRelease *linux.Release
}

type Parser func(source.FileResolver, *Environment, source.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error)
