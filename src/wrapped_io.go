package main

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func CopyFile(src string, dest string) error {
	destF, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destF.Close()

	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	if _, err = io.Copy(destF, srcF); err != nil {
		return err
	}

	return destF.Sync()
}

const MaxDepth = 20

func Reconcile(srcLoc string, dstLoc string, override bool) error {
	srcFiles, err := DeepLs(srcLoc)
	if err != nil {
		return err
	}
	dstFiles, err := DeepLs(dstLoc)
	if err != nil {
		return err
	}
	sort.Strings(srcFiles)
	sort.Strings(dstFiles)

	dstCursor := 0
	dstLen := len(dstFiles)

	for _, srcFile := range srcFiles {

		destFile := filepath.Join(dstLoc, srcFile)

		for dstCursor < dstLen && dstFiles[dstCursor] < destFile {
			dstCursor++
		}
		if !override && dstCursor < dstLen && dstFiles[dstCursor] == destFile {
			continue
		}

		if err = os.MkdirAll(filepath.Dir(destFile), 0755); err != nil {
			return err
		}

		if err = CopyFile(filepath.Join(srcLoc, srcFile), destFile); err != nil {
			return err
		}
	}

	return nil
}

const separator = string(filepath.Separator)

func DeepLs(loc string) ([]string, error) {
	var srcFiles []string
	if err := fs.WalkDir(os.DirFS(loc), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.Count(path, separator) > MaxDepth {
			return errors.New("maximum depth exceeded")
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
