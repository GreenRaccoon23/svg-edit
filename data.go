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
	if DoShutUp {
		return
	}
	if err == nil {
		return
	}
	Log(err)
}

func mkDir(dir string) {
	if _, err := os.Stat(dir); err != nil {
		err := os.MkdirAll(dir, 0777)
		LogErr(err)
	}
}

func walkReplace(path string, file os.FileInfo, err error) error {
	if ext := filepath.Ext(path); ext != ".svg" {
		return nil
	}

	src := path
	dst := fmtDest(path)
	_genDest(dst)

	if _isSymlink(file) {
		_copy(src, dst)
		return nil
	}
	editFileFromPath(src, dst)
	return nil
}

func _isSymlink(file os.FileInfo) bool {
	return (file.Mode()&os.ModeSymlink == os.ModeSymlink)
}

func _genDest(path string) {
	dir := filepath.Dir(path)
	mkDir(dir)
}

func editFileFromPath(src string, dst string) {
	content, err := _fileToString(src)
	if err != nil {
		_copy(src, dst)
		return
	}

	edited := replace(content)
	if edited == "" {
		return
	}

	newFile, err := os.Create(dst)
	if err != nil {
		Log(err)
		return
	}
	defer newFile.Close()

	_stringToFile(edited, newFile)
	TotalEdited += 1
	Progress(dst)
}

func _fileToString(fileName string) (fileString string, err error) {
	var file []byte
	file, err = ioutil.ReadFile(fileName)
	if err != nil {
		Log(err)
		return
	}
	fileString = string(file)
	return
}

func _stringToFile(s string, file *os.File) {
	b := []byte(s)
	_, err := file.Write(b)
	LogErr(err)

	err = file.Sync()
	LogErr(err)
}

func _copy(srcPath, dstPath string) {

	if dstPath == srcPath {
		return
	}

	src, err := os.Open(srcPath)
	if err != nil {
		Log(err)
		return
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		Log(err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		Log(err)
		return
	}

	err = dst.Sync()
	LogErr(err)
}

/*func Copy(srcPath, dstPath string) error {
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
