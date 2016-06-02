package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	var fileContent string
	if err := _fileToString(srcPath, &fileContent); err != nil {

		if fileContent == "" {
			return err
		}

		Log(err)
		return _copyFromPath(dstPath, srcPath)
	}

	var editedFileContent string
	_replace(fileContent, &editedFileContent)
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
	Log(dstPath)

	return nil
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

func _fileToString(path string, fileContent *string) error {

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	*fileContent = string(fileBytes)
	return nil
}

func _replace(fileContent string, fileContentEdited *string) {

	*fileContentEdited = strings.Replace(fileContent, ToFind, ToReplace, -1)

	if wasEdited := (*fileContentEdited != fileContent); wasEdited {
		return
	}

	if shouldAddFill := (DoAddFill && !_hasFill(fileContent)); !shouldAddFill {
		return
	}

	*fileContentEdited = ReAddNew.ReplaceAllString(fileContent, ToAdd)
}

func _hasFill(fileContent string) bool {
	return strings.Contains(fileContent, "fill=") ||
		strings.Contains(fileContent, "fill:")
}

// func _replace(fileContent string) string {

//  re, replacement, needsToBeEdited := _getFindAndReplace(fileContent)
//  if !needsToBeEdited {
//      return fileContent
//  }

//  return re.ReplaceAllString(fileContent, replacement)
// }

// func _getFindAndReplace(fileContent string) (*regexp.Regexp, string, bool) {

//  if nothingToReplace := (!_containsToFind(fileContent)); nothingToReplace {

//      if shouldAddFill := (DoAddFill && !_hasFill(fileContent)); shouldAddFill {
//          return ReAddNew, ToAdd, true
//      }

//      return nil, "", false
//  }

//  return ReToFind, ToReplace, true
// }

// func _getFindAndReplace(fileContent string) (*regexp.Regexp, string, bool) {

// 	if nothingToReplace := (!_containsToFind(fileContent)); nothingToReplace {

// 		if shouldAddFill := (DoAddFill && !_hasFill(fileContent)); shouldAddFill {
// 			return ReAddNew, ToAdd, true
// 		}

// 		return nil, "", false
// 	}

// 	return ReToFind, ToReplace, true
// }

// func _containsToFind(fileContent string) bool {
// 	return strings.Contains(fileContent, ToFind)
// }

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
