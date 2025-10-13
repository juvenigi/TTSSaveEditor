package src

import (
	"ReallyDumbCopyPaste/src/wrapped_io"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// reconcile recursively searches for files in srcLoc and copies them to dstLoc while preserving subfolder structure.
// todo: we follow symlinks so it might be a good idea to add a cap...
// override == true  -> will write to dstLoc without checking
func reconcile(srcLoc string, dstLoc string, override bool) error {
	srcFiles, err := deepLs(srcLoc)
	if err != nil {
		return err
	}
	dstFiles, err := deepLs(dstLoc)
	if err != nil {
		return err
	}
	sort.Strings(srcFiles)
	sort.Strings(dstFiles)

	dstCursor := 0

	for _, srcFile := range srcFiles {
		srcRel, err := filepath.Rel(srcLoc, srcFile)
		if err != nil {
			return err
		}
		destFile := filepath.Join(dstLoc, srcRel)
		// we need to try to find dstFile first. Once we are sure that it does not exist (since our slices are sorted) we can continue on with it
		// todo: working with dstCursor here...
		if !override && false {
			continue
		}

		// all clear, can write
		if err = os.MkdirAll(filepath.Dir(destFile), 0777); err != nil {
			if !errors.Is(err, fs.ErrExist) {
				return err
			}
		}
		if err = wrapped_io.CopyFile(srcFile, destFile); err != nil {
			return err
		}
	}

	return nil
}

func deepLs(loc string) ([]string, error) {
	var srcFiles []string
	if err := fs.WalkDir(os.DirFS(loc), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			srcFiles = append(srcFiles, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return srcFiles, nil
}
