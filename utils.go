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

func slcContains(slc []string, s string) bool {
	lenSlc := len(slc)
	for i := 0; i < lenSlc; i++ {
		if slc[i] == s {
			return true
		}
	}
	return false
}

func slcIsEmpty(slc []string) bool {
	lenSlc := len(slc)
	for i := 0; i < lenSlc; i++ {
		if slc[i] != "" {
			return false
		}
	}
	return true
}

func compact(args ...string) (compacted []string) {
	lenArgs := len(args)
	for i := 0; i < lenArgs; i++ {
		if s := args[i]; s != "" {
			compacted = append(compacted, s)
		}
	}
	return
}

func isEmpty(args ...string) bool {
	lenArgs := len(args)
	for i := 0; i < lenArgs; i++ {
		if notEmpty := args[i] != ""; notEmpty {
			return false
		}
	}
	return true
}

func filter(unfiltered []string, unwanted ...string) (filtered []string) {

	lenUnfiltered := len(unfiltered)
	for i := 0; i < lenUnfiltered; i++ {
		s := unfiltered[i]

		if isUnwanted := slcContains(unwanted, s); isUnwanted {
			continue
		}

		filtered = append(filtered, s)
	}

	return
}

func extract(excess []string, wanted ...string) (extracted []string) {

	lenExcess := len(excess)
	for i := 0; i < lenExcess; i++ {
		s := excess[i]

		if isWanted := slcContains(wanted, s); isWanted {
			extracted = append(extracted, s)
		}
	}

	return
}
