package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

func editOne() error {

	if filepath.Ext(SrcSvg) != ".svg" {
		return fmt.Errorf("%v does not appear to be an svg", SrcSvg)
	}

	if isPathSymlink(SrcSvg) {
		return fmt.Errorf("Cannot edit a symlink")
	}

	return editFileFromPath(DstSvg, SrcSvg)
}

func editRecursive() error {
	return editRecursiveFast()
	// return editRecursiveSlow()
}

func editRecursiveFast() error {

	svgPaths := getSvgPaths()
	lenSvgPaths := len(svgPaths)

	var wg sync.WaitGroup
	wg.Add(lenSvgPaths)

	for i := 0; i < lenSvgPaths; i++ {
		go func(i int) {
			defer wg.Done()

			svgPath := svgPaths[i]

			err := _editSvg(svgPath)
			if err != nil {
				LogErr(err)
			}
		}(i)
	}

	wg.Wait()

	return nil
}

func _editSvg(svgPath string) error {

	dstPath := fmtDst(svgPath)
	if err := _mkDstDir(dstPath); err != nil {
		return err
	}

	srcPath := svgPath

	if err := editFileFromPath(dstPath, srcPath); err != nil {
		return err
	}

	return nil
}

func _mkDstDir(dstPath string) error {
	dstDir := filepath.Dir(dstPath)
	return mkDir(dstDir)
}

func editFileFromPath(dstPath string, srcPath string) error {

	fileBytes, err := ioutil.ReadFile(srcPath)
	if failedToReadFile := (err != nil); failedToReadFile {
		LogErr(err)
		return copyFromPath(dstPath, srcPath)
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

func _editFileBytes(fileBytes *[]byte, editedFileBytes *[]byte) (wasEdited bool) {

	wasEdited = _replace(fileBytes, editedFileBytes)
	if wasEdited {
		return
	}

	if shouldAddFill := (DoAddFill && !_hasFill(fileBytes)); !shouldAddFill {
		return
	}

	wasEdited = _addFill(fileBytes, editedFileBytes)
	return
}

func _replace(fileBytes *[]byte, editedFileBytes *[]byte) (wasEdited bool) {

	*editedFileBytes = bytes.Replace(*fileBytes, ToFindBytes, ToReplaceBytes, -1)

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

func editRecursiveSlow() error {
	return filepath.Walk(SrcDir, _walkReplace)
}

func _walkReplace(path string, fi os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	if !isEditable(fi) {
		return nil
	}

	return _editSvg(path)
}
