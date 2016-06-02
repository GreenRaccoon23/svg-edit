package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func getPwd() string {

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return pwd
}

func getSvgPaths() (svgPaths []string) {

	filepath.Walk(SrcDir, func(path string, fi os.FileInfo, err error) error {

		if !isEditable(fi) {
			return nil
		}

		svgPaths = append(svgPaths, path)
		return nil
	})

	return
}

func isEditable(fi os.FileInfo) bool {

	if filepath.Ext(fi.Name()) != ".svg" {
		return false
	}

	if isSymlink := (fi.Mode()&os.ModeSymlink == os.ModeSymlink); isSymlink {
		return nil
	}
}

func mkDir(path string) error {
	if _pathExists(path) {
		return nil
	}
	return os.MkdirAll(path, 0777)
}

func _pathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func isPathSymlink(path string) bool {

	fi, err := os.Lstat(path)
	if err != nil {
		LogErr(err)
		return false
	}

	return (fi.Mode()&os.ModeSymlink == os.ModeSymlink)
}

func copyFromPath(dstPath string, srcPath string) error {

	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	return _copyFile(dst, src)
}

func _copyFile(dst *os.File, src *os.File) error {

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return dst.Sync()
}
