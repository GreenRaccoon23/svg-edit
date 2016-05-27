package main

import (
	"bytes"
	"regexp"
	"strings"
)

func concat(args ...string) string {

	var b bytes.Buffer
	defer b.Reset()

	for _, s := range args {
		b.WriteString(s)
	}

	return b.String()
}

func addExt(path string, ext string) string {

	if path == "" {
		return ""
	}

	if hasExt := strings.HasSuffix(path, ext); hasExt {
		return path
	}

	return concat(path, ext)
}

func fmtDir(dir string) string {

	if dir == "" {
		return ""
	}

	formatted := dir

	if isFirstChar(dir, "/", "~") == false {
		formatted = concat(Pwd, "/", dir)
	}

	if isLastChar(dir, "/") == false {
		formatted = concat(formatted, "/")
	}

	if dir == "." || dir == "./" {
		formatted = Pwd
	}

	return formatted
}

func isFirstChar(s string, args ...string) bool {
	firstChar := string(s[0])
	for _, a := range args {
		if firstChar == a {
			return true
		}
	}
	return false
}

func isLastChar(s string, args ...string) bool {
	lastChar := string(s[len(s)-1])
	for _, z := range args {
		if lastChar == z {
			return true
		}
	}
	return false
}

// func fmtCopy(s string) string {
// 	if strings.HasSuffix(s, ".svg") {
// 		return strings.Replace(s, ".svg", "-copy.svg", 1)
// 	}
// 	return concat(s, "-copy.svg")
// }

func fmtDst(path string) (dstPath string) {
	dstPath = strings.Replace(path, SrcDir, DstDir, 1)
	dstPath = strings.Replace(dstPath, "//", "/", -1)
	return
}

func replace(fileContent string) string {

	re, replacement, needsToBeEdited := _getFindAndReplace(fileContent)
	if !needsToBeEdited {
		return fileContent
	}

	return re.ReplaceAllString(fileContent, replacement)
}

func _getFindAndReplace(fileContent string) (*regexp.Regexp, string, bool) {

	if nothingToReplace := (!_containsToFind(fileContent)); nothingToReplace {

		if shouldAddFill := (DoAddNew && !_hasFill(fileContent)); shouldAddFill {
			return ReAddNew, ToAdd, true
		}

		return nil, "", false
	}

	return ReToFind, ToReplace, true
}

func _containsToFind(fileContent string) bool {
	return strings.Contains(fileContent, ToFind)
}

func _hasFill(fileContent string) bool {
	return strings.Contains(fileContent, "fill=") ||
		strings.Contains(fileContent, "fill:")
}

func pop(slc []string) (string, []string) {
	iEnd := len(slc) - 1
	return slc[iEnd], slc[:iEnd]
}

func cut(slc []string, i int, j int) []string {

	if copyAll := (i == 0 && j == -1); copyAll {
		return copySlc(slc)
	}

	if goToEnd := (j == -1); goToEnd {
		return slc[i:]
	}

	return slc[i:j]
}

func copySlc(slc []string) []string {
	newSlc := make([]string, len(slc))
	copy(newSlc, slc)
	return newSlc
}

func shift(slc []string) (string, []string) {
	return slc[0], slc[1:]
}

func unshift(slc []string, s string) []string {
	return append([]string{s}, slc...)
}
