package integration

import (
	"testing"

	"github.com/nextlinux/syft/syft/source"
)

func TestRegressionJavaNoMainPackage(t *testing.T) { // Regression: https://github.com/nextlinux/syft/issues/252
	catalogFixtureImage(t, "image-java-no-main-package", source.SquashedScope, nil)
}
