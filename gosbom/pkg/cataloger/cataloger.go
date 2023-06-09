/*
Package cataloger provides the ability to process files from a container image or file system and discover packages
(gems, wheels, jars, rpms, debs, etc). Specifically, this package contains both a catalog function to utilize all
catalogers defined in child packages as well as the interface definition to implement a cataloger.
*/
package cataloger

import (
	"strings"

	"github.com/nextlinux/gosbom/internal/log"
	"github.com/nextlinux/gosbom/gosbom/pkg"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/alpm"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/apkdb"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/binary"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/cpp"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/dart"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/deb"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/dotnet"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/elixir"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/erlang"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/golang"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/haskell"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/java"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/javascript"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/kernel"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/nix"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/php"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/portage"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/python"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/rpm"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/ruby"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/rust"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/sbom"
	"github.com/nextlinux/gosbom/gosbom/pkg/cataloger/swift"
)

const AllCatalogersPattern = "all"

// ImageCatalogers returns a slice of locally implemented catalogers that are fit for detecting installations of packages.
func ImageCatalogers(cfg Config) []pkg.Cataloger {
	return filterCatalogers([]pkg.Cataloger{
		alpm.NewAlpmdbCataloger(),
		ruby.NewGemSpecCataloger(),
		python.NewPythonPackageCataloger(),
		php.NewComposerInstalledCataloger(),
		javascript.NewPackageCataloger(),
		deb.NewDpkgdbCataloger(),
		rpm.NewRpmDBCataloger(),
		java.NewJavaCataloger(cfg.Java()),
		java.NewNativeImageCataloger(),
		apkdb.NewApkdbCataloger(),
		golang.NewGoModuleBinaryCataloger(cfg.Go()),
		dotnet.NewDotnetDepsCataloger(),
		portage.NewPortageCataloger(),
		nix.NewStoreCataloger(),
		sbom.NewSBOMCataloger(),
		binary.NewCataloger(),
		kernel.NewLinuxKernelCataloger(cfg.Kernel()),
	}, cfg.Catalogers)
}

// DirectoryCatalogers returns a slice of locally implemented catalogers that are fit for detecting packages from index files (and select installations)
func DirectoryCatalogers(cfg Config) []pkg.Cataloger {
	return filterCatalogers([]pkg.Cataloger{
		alpm.NewAlpmdbCataloger(),
		ruby.NewGemFileLockCataloger(),
		python.NewPythonIndexCataloger(),
		python.NewPythonPackageCataloger(),
		php.NewComposerLockCataloger(),
		javascript.NewLockCataloger(),
		deb.NewDpkgdbCataloger(),
		rpm.NewRpmDBCataloger(),
		rpm.NewFileCataloger(),
		java.NewJavaCataloger(cfg.Java()),
		java.NewJavaPomCataloger(),
		java.NewNativeImageCataloger(),
		java.NewJavaGradleLockfileCataloger(),
		apkdb.NewApkdbCataloger(),
		golang.NewGoModuleBinaryCataloger(cfg.Go()),
		golang.NewGoModFileCataloger(cfg.Go()),
		rust.NewCargoLockCataloger(),
		dart.NewPubspecLockCataloger(),
		dotnet.NewDotnetDepsCataloger(),
		swift.NewCocoapodsCataloger(),
		cpp.NewConanCataloger(),
		portage.NewPortageCataloger(),
		haskell.NewHackageCataloger(),
		sbom.NewSBOMCataloger(),
		binary.NewCataloger(),
		elixir.NewMixLockCataloger(),
		erlang.NewRebarLockCataloger(),
		kernel.NewLinuxKernelCataloger(cfg.Kernel()),
		nix.NewStoreCataloger(),
	}, cfg.Catalogers)
}

// AllCatalogers returns all implemented catalogers
func AllCatalogers(cfg Config) []pkg.Cataloger {
	return filterCatalogers([]pkg.Cataloger{
		alpm.NewAlpmdbCataloger(),
		ruby.NewGemFileLockCataloger(),
		ruby.NewGemSpecCataloger(),
		python.NewPythonIndexCataloger(),
		python.NewPythonPackageCataloger(),
		javascript.NewLockCataloger(),
		javascript.NewPackageCataloger(),
		deb.NewDpkgdbCataloger(),
		rpm.NewRpmDBCataloger(),
		rpm.NewFileCataloger(),
		java.NewJavaCataloger(cfg.Java()),
		java.NewJavaPomCataloger(),
		java.NewNativeImageCataloger(),
		java.NewJavaGradleLockfileCataloger(),
		apkdb.NewApkdbCataloger(),
		golang.NewGoModuleBinaryCataloger(cfg.Go()),
		golang.NewGoModFileCataloger(cfg.Go()),
		rust.NewCargoLockCataloger(),
		rust.NewAuditBinaryCataloger(),
		dart.NewPubspecLockCataloger(),
		dotnet.NewDotnetDepsCataloger(),
		php.NewComposerInstalledCataloger(),
		php.NewComposerLockCataloger(),
		swift.NewCocoapodsCataloger(),
		cpp.NewConanCataloger(),
		portage.NewPortageCataloger(),
		haskell.NewHackageCataloger(),
		sbom.NewSBOMCataloger(),
		binary.NewCataloger(),
		elixir.NewMixLockCataloger(),
		erlang.NewRebarLockCataloger(),
		kernel.NewLinuxKernelCataloger(cfg.Kernel()),
		nix.NewStoreCataloger(),
	}, cfg.Catalogers)
}

func RequestedAllCatalogers(cfg Config) bool {
	for _, enableCatalogerPattern := range cfg.Catalogers {
		if enableCatalogerPattern == AllCatalogersPattern {
			return true
		}
	}
	return false
}

func filterCatalogers(catalogers []pkg.Cataloger, enabledCatalogerPatterns []string) []pkg.Cataloger {
	// if cataloger is not set, all applicable catalogers are enabled by default
	if len(enabledCatalogerPatterns) == 0 {
		return catalogers
	}
	for _, enableCatalogerPattern := range enabledCatalogerPatterns {
		if enableCatalogerPattern == AllCatalogersPattern {
			return catalogers
		}
	}
	var keepCatalogers []pkg.Cataloger
	for _, cataloger := range catalogers {
		if contains(enabledCatalogerPatterns, cataloger.Name()) {
			keepCatalogers = append(keepCatalogers, cataloger)
			continue
		}
		log.Infof("skipping cataloger %q", cataloger.Name())
	}
	return keepCatalogers
}

func contains(enabledPartial []string, catalogerName string) bool {
	catalogerName = strings.TrimSuffix(catalogerName, "-cataloger")
	for _, partial := range enabledPartial {
		partial = strings.TrimSuffix(partial, "-cataloger")
		if partial == "" {
			continue
		}
		if hasFullWord(partial, catalogerName) {
			return true
		}
	}
	return false
}

func hasFullWord(targetPhrase, candidate string) bool {
	if targetPhrase == "cataloger" || targetPhrase == "" {
		return false
	}
	start := strings.Index(candidate, targetPhrase)
	if start == -1 {
		return false
	}

	if start > 0 && candidate[start-1] != '-' {
		return false
	}

	end := start + len(targetPhrase)
	if end < len(candidate) && candidate[end] != '-' {
		return false
	}
	return true
}
