package main

import (
	"bytes"
	"io"
	"io/ioutil"
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

func walkReplace(path string, fi os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	if filepath.Ext(path) != ".svg" {
		return nil
	}

	dstPath := fmtDst(path)
	if err = _mkDstDir(dstPath); err != nil {
		LogErr(err)
		return nil
	}

	srcPath := path

	if err = editFileFromPath(dstPath, srcPath); err != nil {
		LogErr(err)
		return nil
	}

	return nil
}

func _isPathSymlink(path string) bool {

	fi, err := os.Lstat(path)
	if err != nil {
		LogErr(err)
		return false
	}

	return (fi.Mode()&os.ModeSymlink == os.ModeSymlink)
}

func _mkDstDir(path string) error {
	dir := filepath.Dir(path)
	return mkDir(dir)
}

func editFileFromPath(dstPath string, srcPath string) error {

	if _isPathSymlink(srcPath) {
		// return _copySymlinkFromPath(dstPath, srcPath)
		return nil
	}

	fileBytes, err := ioutil.ReadFile(srcPath)
	if failedToReadFile := (err != nil); failedToReadFile {
		LogErr(err)
		return _copyFromPath(dstPath, srcPath)
	}

	if isEmptyFile := (len(fileBytes) == 0); isEmptyFile {
		return nil
	}

	var editedFileBytes []byte
	if wasEdited := _editFileBytes(&fileBytes, &editedFileBytes); !wasEdited {
		return nil
	}

	if somethingTerribleHappened := (len(editedFileBytes) == 0); somethingTerribleHappened {
		return nil
	}

	newFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	if err = _bytesToFile(&editedFileBytes, newFile); err != nil {
		return err
	}

	TotalEdited++
	Log(dstPath)

	return nil
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

func _fileToString(path string, fileContent *string) error {

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	*fileContent = string(fileBytes)
	return nil
}

func _stringToFile(editedFileContent *string, newFile *os.File) error {

	b := []byte(*editedFileContent)
	if _, err := newFile.Write(b); err != nil {
		return err
	}

	return newFile.Sync()
}

func _editFileBytes(fileBytes *[]byte, editedFileBytes *[]byte) (wasEdited bool) {

	wasEdited = _replace(editedFileBytes, fileBytes)
	if wasEdited {
		return
	}

	if shouldAddFill := (DoAddFill && !_hasFill(fileBytes)); !shouldAddFill {
		return
	}

	wasEdited = _addFill(editedFileBytes, fileBytes)
	return
}

func _replace(fileBytes *[]byte, editedFileBytes *[]byte) (wasEdited bool) {

	*editedFileBytes = bytes.Replace(*fileBytes, ToFindBytes, ToReplaceBytes, -1)
	// *editedFileBytes = ReToFind.ReplaceAll(*fileBytes, ToReplaceBytes)

	wasEdited = (!bytes.Equal(*editedFileBytes, *fileBytes))
	return
}

func _addFill(fileBytes *[]byte, editedFileBytes *[]byte) (wasEdited bool) {

	*editedFileBytes = ReAddFill.ReplaceAll(*fileBytes, ToFillBytes)

	wasEdited = (!bytes.Equal(*editedFileBytes, *fileBytes))
	return
}

func _hasFill(fileBytes *[]byte) bool {
	return bytes.Contains(*fileBytes, []byte("fill=")) ||
		bytes.Contains(*fileBytes, []byte("fill:"))
}

func _bytesToFile(editedFileBytes *[]byte, newFile *os.File) error {

	if _, err := newFile.Write(*editedFileBytes); err != nil {
		return err
	}

	return newFile.Sync()
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
