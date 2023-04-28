package source

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/wagoodman/go-partybus"
	"github.com/wagoodman/go-progress"

	"github.com/nextlinux/stereoscope/pkg/file"
	"github.com/nextlinux/stereoscope/pkg/filetree"
	"github.com/nextlinux/syft/internal"
	"github.com/nextlinux/syft/internal/bus"
	"github.com/nextlinux/syft/internal/log"
	"github.com/nextlinux/syft/syft/event"
)

type pathIndexVisitor func(string, os.FileInfo, error) error

type directoryIndexer struct {
	path              string
	base              string
	pathIndexVisitors []pathIndexVisitor
	errPaths          map[string]error
	tree              filetree.ReadWriter
	index             filetree.Index
}

func newDirectoryIndexer(path, base string, visitors ...pathIndexVisitor) *directoryIndexer {
	i := &directoryIndexer{
		path:              path,
		base:              base,
		tree:              filetree.New(),
		index:             filetree.NewIndex(),
		pathIndexVisitors: append([]pathIndexVisitor{requireFileInfo, disallowByFileType, disallowUnixSystemRuntimePath}, visitors...),
		errPaths:          make(map[string]error),
	}

	// these additional stateful visitors should be the first thing considered when walking / indexing
	i.pathIndexVisitors = append(
		[]pathIndexVisitor{
			i.disallowRevisitingVisitor,
			i.disallowFileAccessErr,
		},
		i.pathIndexVisitors...,
	)

	return i
}

func (r *directoryIndexer) build() (filetree.Reader, filetree.IndexReader, error) {
	return r.tree, r.index, indexAllRoots(r.path, r.indexTree)
}

func indexAllRoots(root string, indexer func(string, *progress.Stage) ([]string, error)) error {
	// why account for multiple roots? To cover cases when there is a symlink that references above the root path,
	// in which case we need to additionally index where the link resolves to. it's for this reason why the filetree
	// must be relative to the root of the filesystem (and not just relative to the given path).
	pathsToIndex := []string{root}
	fullPathsMap := map[string]struct{}{}

	stager, prog := indexingProgress(root)
	defer prog.SetCompleted()
loop:
	for {
		var currentPath string
		switch len(pathsToIndex) {
		case 0:
			break loop
		case 1:
			currentPath, pathsToIndex = pathsToIndex[0], nil
		default:
			currentPath, pathsToIndex = pathsToIndex[0], pathsToIndex[1:]
		}

		additionalRoots, err := indexer(currentPath, stager)
		if err != nil {
			return fmt.Errorf("unable to index filesystem path=%q: %w", currentPath, err)
		}

		for _, newRoot := range additionalRoots {
			if _, ok := fullPathsMap[newRoot]; !ok {
				fullPathsMap[newRoot] = struct{}{}
				pathsToIndex = append(pathsToIndex, newRoot)
			}
		}
	}

	return nil
}

func (r *directoryIndexer) indexTree(root string, stager *progress.Stage) ([]string, error) {
	log.WithFields("path", root).Trace("indexing filetree")

	var roots []string
	var err error

	root, err = filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	// we want to be able to index single files with the directory resolver. However, we should also allow for attempting
	// to index paths that do not exist (that is, a root that does not exist is not an error case that should stop indexing).
	// For this reason we look for an opportunity to discover if the given root is a file, and if so add a single root,
	// but continue forth with index regardless if the given root path exists or not.
	fi, err := os.Stat(root)
	if err != nil && fi != nil && !fi.IsDir() {
		// note: we want to index the path regardless of an error stat-ing the path
		newRoot, _ := r.indexPath(root, fi, nil)
		if newRoot != "" {
			roots = append(roots, newRoot)
		}
		return roots, nil
	}

	err = filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			stager.Current = path

			newRoot, err := r.indexPath(path, info, err)

			if err != nil {
				return err
			}

			if newRoot != "" {
				roots = append(roots, newRoot)
			}

			return nil
		})

	if err != nil {
		return nil, fmt.Errorf("unable to index root=%q: %w", root, err)
	}

	return roots, nil
}

func (r *directoryIndexer) indexPath(path string, info os.FileInfo, err error) (string, error) {
	// ignore any path which a filter function returns true
	for _, filterFn := range r.pathIndexVisitors {
		if filterFn == nil {
			continue
		}

		if filterErr := filterFn(path, info, err); filterErr != nil {
			if errors.Is(filterErr, fs.SkipDir) {
				// signal to walk() to skip this directory entirely (even if we're processing a file)
				return "", filterErr
			}
			// skip this path but don't affect walk() trajectory
			return "", nil
		}
	}

	if info == nil {
		// walk may not be able to provide a FileInfo object, don't allow for this to stop indexing; keep track of the paths and continue.
		r.errPaths[path] = fmt.Errorf("no file info observable at path=%q", path)
		return "", nil
	}

	// here we check to see if we need to normalize paths to posix on the way in coming from windows
	if runtime.GOOS == WindowsOS {
		path = windowsToPosix(path)
	}

	newRoot, err := r.addPathToIndex(path, info)
	if r.isFileAccessErr(path, err) {
		return "", nil
	}

	return newRoot, nil
}

func (r *directoryIndexer) disallowFileAccessErr(path string, _ os.FileInfo, err error) error {
	if r.isFileAccessErr(path, err) {
		return errSkipPath
	}
	return nil
}

func (r *directoryIndexer) isFileAccessErr(path string, err error) bool {
	// don't allow for errors to stop indexing, keep track of the paths and continue.
	if err != nil {
		log.Warnf("unable to access path=%q: %+v", path, err)
		r.errPaths[path] = err
		return true
	}
	return false
}

func (r directoryIndexer) addPathToIndex(p string, info os.FileInfo) (string, error) {
	switch t := file.TypeFromMode(info.Mode()); t {
	case file.TypeSymLink:
		return r.addSymlinkToIndex(p, info)
	case file.TypeDirectory:
		return "", r.addDirectoryToIndex(p, info)
	case file.TypeRegular:
		return "", r.addFileToIndex(p, info)
	default:
		return "", fmt.Errorf("unsupported file type: %s", t)
	}
}

func (r directoryIndexer) addDirectoryToIndex(p string, info os.FileInfo) error {
	ref, err := r.tree.AddDir(file.Path(p))
	if err != nil {
		return err
	}

	metadata := file.NewMetadataFromPath(p, info)
	r.index.Add(*ref, metadata)

	return nil
}

func (r directoryIndexer) addFileToIndex(p string, info os.FileInfo) error {
	ref, err := r.tree.AddFile(file.Path(p))
	if err != nil {
		return err
	}

	metadata := file.NewMetadataFromPath(p, info)
	r.index.Add(*ref, metadata)

	return nil
}

func (r directoryIndexer) addSymlinkToIndex(p string, info os.FileInfo) (string, error) {
	linkTarget, err := os.Readlink(p)
	if err != nil {
		return "", fmt.Errorf("unable to readlink for path=%q: %w", p, err)
	}

	if filepath.IsAbs(linkTarget) {
		// if the link is absolute (e.g, /bin/ls -> /bin/busybox) we need to
		// resolve relative to the root of the base directory
		linkTarget = filepath.Join(r.base, filepath.Clean(linkTarget))
	} else {
		// if the link is not absolute (e.g, /dev/stderr -> fd/2 ) we need to
		// resolve it relative to the directory in question (e.g. resolve to
		// /dev/fd/2)
		if r.base == "" {
			linkTarget = filepath.Join(filepath.Dir(p), linkTarget)
		} else {
			// if the base is set, then we first need to resolve the link,
			// before finding it's location in the base
			dir, err := filepath.Rel(r.base, filepath.Dir(p))
			if err != nil {
				return "", fmt.Errorf("unable to resolve relative path for path=%q: %w", p, err)
			}
			linkTarget = filepath.Join(r.base, filepath.Clean(filepath.Join("/", dir, linkTarget)))
		}
	}

	ref, err := r.tree.AddSymLink(file.Path(p), file.Path(linkTarget))
	if err != nil {
		return "", err
	}

	targetAbsPath := linkTarget
	if !filepath.IsAbs(targetAbsPath) {
		targetAbsPath = filepath.Clean(filepath.Join(path.Dir(p), linkTarget))
	}

	metadata := file.NewMetadataFromPath(p, info)
	metadata.LinkDestination = linkTarget
	r.index.Add(*ref, metadata)

	return targetAbsPath, nil
}

func (r directoryIndexer) hasBeenIndexed(p string) (bool, *file.Metadata) {
	filePath := file.Path(p)
	if !r.tree.HasPath(filePath) {
		return false, nil
	}

	exists, ref, err := r.tree.File(filePath)
	if err != nil || !exists || !ref.HasReference() {
		return false, nil
	}

	// cases like "/" will be in the tree, but not been indexed yet (a special case). We want to capture
	// these cases as new paths to index.
	if !ref.HasReference() {
		return false, nil
	}

	entry, err := r.index.Get(*ref.Reference)
	if err != nil {
		return false, nil
	}

	return true, &entry.Metadata
}

func (r *directoryIndexer) disallowRevisitingVisitor(path string, _ os.FileInfo, _ error) error {
	// this prevents visiting:
	// - link destinations twice, once for the real file and another through the virtual path
	// - infinite link cycles
	if indexed, metadata := r.hasBeenIndexed(path); indexed {
		if metadata.IsDir {
			// signal to walk() that we should skip this directory entirely
			return fs.SkipDir
		}
		return errSkipPath
	}
	return nil
}

func disallowUnixSystemRuntimePath(path string, _ os.FileInfo, _ error) error {
	if internal.HasAnyOfPrefixes(path, unixSystemRuntimePrefixes...) {
		return fs.SkipDir
	}
	return nil
}

func disallowByFileType(_ string, info os.FileInfo, _ error) error {
	if info == nil {
		// we can't filter out by filetype for non-existent files
		return nil
	}
	switch file.TypeFromMode(info.Mode()) {
	case file.TypeCharacterDevice, file.TypeSocket, file.TypeBlockDevice, file.TypeFIFO, file.TypeIrregular:
		return errSkipPath
		// note: symlinks that point to these files may still get by.
		// We handle this later in processing to help prevent against infinite links traversal.
	}

	return nil
}

func requireFileInfo(_ string, info os.FileInfo, _ error) error {
	if info == nil {
		return errSkipPath
	}
	return nil
}

func indexingProgress(path string) (*progress.Stage, *progress.Manual) {
	stage := &progress.Stage{}
	prog := progress.NewManual(-1)

	bus.Publish(partybus.Event{
		Type:   event.FileIndexingStarted,
		Source: path,
		Value: struct {
			progress.Stager
			progress.Progressable
		}{
			Stager:       progress.Stager(stage),
			Progressable: prog,
		},
	})

	return stage, prog
}
