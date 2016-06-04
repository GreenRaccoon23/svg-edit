package main

import (
	"fmt"
	"io"
	"io/ioutil"
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

func getPwd() (pwd string) {
	var err error
	pwd, err = os.Getwd()
	LogErr(err)
	return
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

func IsStringSymlink(filename string) bool {
	file, _ := os.Lstat(filename)
	return IsSymlink(file)
}

func IsSymlink(file os.FileInfo) bool {
	if file.Mode()&os.ModeSymlink == os.ModeSymlink {
		return true
	}
	return false
}

func Create(fileName string) (file *os.File) {
	var err error
	file, err = os.Create(fileName)
	LogErr(err)
	return
}

func MakeDir(dir string) {
	if _, err := os.Stat(dir); err != nil {
		err := os.MkdirAll(dir, 0777)
		LogErr(err)
	}
}

func RemoveIfExists(fileName string) {
	if _, err := os.Stat(fileName); err == nil {
		os.Remove(fileName)
	}
}

func WalkReplace(path string, file os.FileInfo, err error) error {
	ext := filepath.Ext(path)
	if ext != ".svg" {
		return nil
	}

	in := path
	out := fmtDest(path)
	genDest(out)

	if IsSymlink(file) {
		Copy(in, out)
		return nil
	}
	edit(in, out)
	return nil
}

func genDest(path string) {
	dir := filepath.Dir(path)
	MakeDir(dir)
	return
}

func edit(in string, out string) {
	content, err := FileToString(in)
	if err != nil {
		Copy(in, out)
		return
	}

	edited := replace(content)
	if edited == "" {
		return
	}

	newFile, err := os.Create(out)
	if err != nil {
		Log(err)
		return
	}
	defer newFile.Close()

	StringToFile(edited, newFile)
	TotalEdited += 1
	Progress(out)
}

func FileToString(fileName string) (fileString string, err error) {
	var file []byte
	file, err = ioutil.ReadFile(fileName)
	if err != nil {
		Log(err)
		return
	}
	fileString = string(file)
	return
}

func StringToFile(s string, file *os.File) {
	b := []byte(s)
	_, err := file.Write(b)
	LogErr(err)

	err = file.Sync()
	LogErr(err)
}

func Copy(source, destination string) {
	if destination == source {
		return
	}

	toRead, err := os.Open(source)
	if err != nil {
		Log(err)
		return
	}
	defer toRead.Close()

	toWrite, err := os.Create(destination)
	if err != nil {
		Log(err)
		return
	}
	defer toWrite.Close()

	_, err = io.Copy(toWrite, toRead)
	if err != nil {
		Log(err)
		return
	}

	err = toWrite.Sync()
	LogErr(err)
}

/*func Copy(source, destination string) error {
	toRead, err := os.Open(source)
	if err != nil {
		return err
	}
	defer toRead.Close()

	toWrite, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer toWrite.Close()

	_, err = io.Copy(toWrite, toRead)
	if err != nil {
		return err
	}
	return
}*/
