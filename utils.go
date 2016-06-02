package main

import (
	"bytes"
	"log"
	"path/filepath"
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

func toBytes(args ...string) []byte {

	var b bytes.Buffer
	defer b.Reset()

	for _, s := range args {
		b.WriteString(s)
	}

	return b.Bytes()
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

	formatted, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}

	return formatted
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

// func pop(slc []string) (string, []string) {
// 	iEnd := len(slc) - 1
// 	return slc[iEnd], slc[:iEnd]
// }

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

// func shift(slc []string) (string, []string) {
// 	return slc[0], slc[1:]
// }

// func unshift(slc []string, s string) []string {
// 	return append([]string{s}, slc...)
// }
