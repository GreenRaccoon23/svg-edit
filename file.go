package main

import (
	"io"
	"log"
	"os"
)

func getPwd() string {

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return pwd
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

func _isPathSymlink(path string) bool {

	fi, err := os.Lstat(path)
	if err != nil {
		LogErr(err)
		return false
	}

	return (fi.Mode()&os.ModeSymlink == os.ModeSymlink)
}

func _copyFromPath(dstPath string, srcPath string) error {

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

// func _copySymlinkFromPath(dstPath string, srcPath string) error {

//  if SrcDstSame {
//      return nil
//  }

//  linkTarget, err := os.Readlink(srcPath)
//  if err != nil {
//      return err
//  }

//  _, err = filepath.Rel(linkTarget, SrcDir)
//  if isUnderSrcDir := (err == nil); !isUnderSrcDir {
//      return nil
//  }

//  return _copyFromPath(dstPath, srcPath)
// }
