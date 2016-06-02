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

	err, _ := editFileFromPath(DstSvg, SrcSvg)
	return err
}

func editRecursive() error {
	return _editRecursiveFast()
}

func _editRecursiveFast() error {
	svgPaths := getSvgPaths()
	return _editSvgs(svgPaths)
}

func _editSvgs(svgPaths []string) error {

	lenSvgPaths := len(svgPaths)

	var wg sync.WaitGroup
	wg.Add(lenSvgPaths)
	chanEdited := make(chan bool, lenSvgPaths)

	for i := 0; i < lenSvgPaths; i++ {
		go func(i int) {
			defer wg.Done()

			svgPath := svgPaths[i]

			err, wasEdited := _editSvg(svgPath)
			if err != nil {
				LogErr(err)
			} else if wasEdited {
				Log(svgPath)
			}

			chanEdited <- wasEdited
		}(i)
	}

	wg.Wait()
	close(chanEdited)

	for wasEdited := range chanEdited {
		if wasEdited {
			TotalEdited++
		}
	}

	return nil
}

func _editSvg(svgPath string) (error, bool) {

	dstPath := svgPath
	if err := _mkDstDir(dstPath); err != nil {
		return err, false
	}

	srcPath := svgPath

	err, wasEdited := editFileFromPath(dstPath, srcPath)
	if err != nil {
		return err, wasEdited
	}

	return nil, wasEdited
}

func _mkDstDir(dstPath string) error {
	dstDir := filepath.Dir(dstPath)
	return mkDir(dstDir)
}

func editFileFromPath(dstPath string, srcPath string) (error, bool) {

	fileBytes, err := ioutil.ReadFile(srcPath)
	if failedToReadFile := (err != nil); failedToReadFile {
		LogErr(err)
		// return copyFromPath(dstPath, srcPath)
		return nil, false
	}

	if isEmptyFile := (len(fileBytes) == 0); isEmptyFile {
		return nil, false
	}

	var editedFileBytes []byte
	if wasEdited := _editFileBytes(&fileBytes, &editedFileBytes); !wasEdited {
		return nil, false
	}

	if somethingTerribleHappened := (len(editedFileBytes) == 0); somethingTerribleHappened {
		return nil, false
	}

	newFile, err := os.Create(dstPath)
	if err != nil {
		return err, false
	}
	defer newFile.Close()

	if err = _bytesToFile(&editedFileBytes, newFile); err != nil {
		return err, false
	}

	return nil, true
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
