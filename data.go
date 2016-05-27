package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	MaterialDesign map[string]string = map[string]string{
		"black":      "#000000",
		"red":        "#F44336",
		"pink":       "#E91E63",
		"purple":     "#9C27B0",
		"deeppurple": "#673AB7",
		"indigo":     "#3F51B5",
		"blue":       "#2196F3",
		"lightblue":  "#03A9F4",
		"cyan":       "#00BCD4",
		"teal":       "#009688",
		"green":      "#4CAF50",
		"kellygreen": "#00C853",
		"shamrock":   "#00E676",
		"lightgreen": "#8BC34A",
		"lime":       "#CDDC39",
		"yellow":     "#FFEB3B",
		"amber":      "#FFC107",
		"orange":     "#FF9800",
		"deeporange": "#FF5722",
		"brown":      "#795548",
		"grey":       "#9E9E9E",
		"bluegrey":   "#607D8B",
		"archblue":   "#1793D1",
	}
)

func getPwd() string {

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return pwd
}

func Log(err error) {
	if DoShutUp {
		return
	}
	fmt.Println(err)
}

func LogErr(err error) {
	if DoShutUp || err == nil {
		return
	}
	Log(err)
}

func mkDir(dir string) error {
	if _pathExists(dir) {
		return nil
	}
	return os.MkdirAll(dir, 0777)
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
		Log(err)
		return nil
	}

	srcPath := path

	if err = editFileFromPath(dstPath, srcPath); err != nil {
		Log(err)
		return nil
	}

	return nil
}

func _isPathSymlink(path string) bool {

	fi, err := os.Lstat(path)
	if err != nil {
		Log(err)
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

	fileContent, err := _fileToString(srcPath)
	if err != nil {

		if fileContent == "" {
			return err
		}

		Log(err)
		return _copyFromPath(dstPath, srcPath)
	}

	editedFileContent := replace(fileContent)
	if editedFileContent == "" || editedFileContent == fileContent {
		return nil
	}

	newFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	if err = _stringToFile(editedFileContent, newFile); err != nil {
		return err
	}

	TotalEdited++
	Progress(dstPath)

	return nil
}

// func _copySymlinkFromPath(dstPath string, srcPath string) error {

// 	if SrcDstSame {
// 		return nil
// 	}

// 	linkTarget, err := os.Readlink(srcPath)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = filepath.Rel(linkTarget, SrcDir)
// 	if isUnderSrcDir := (err == nil); !isUnderSrcDir {
// 		return nil
// 	}

// 	return _copyFromPath(dstPath, srcPath)
// }

func _stringToFile(s string, file *os.File) error {

	b := []byte(s)
	if _, err := file.Write(b); err != nil {
		return err
	}

	return file.Sync()
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

func _fileToString(fileName string) (string, error) {

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(file), nil
}

/*func _copyFromPath(srcPath, dstPath string) error {
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

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return
}*/
