package main

import (
	"bytes"
	"regexp"
	"strings"
)

var (
	buffer bytes.Buffer
)

func Str(slice []string) (concatenated string) {
	for _, s := range slice {
		buffer.WriteString(s)
	}
	concatenated = buffer.String()
	buffer.Reset()
	return
}

func Slc(args ...string) []string {
	return args
}

func concat(args ...string) string {
	return Str(args)
}

func IsFirstLetter(s string, args ...string) bool {
	firstLetter := string(s[0])
	for _, a := range args {
		if firstLetter == a {
			return true
		}
	}
	return false
}

func IsLastLetter(s string, args ...string) bool {
	lastLetter := string(s[len(s)-1])
	for _, z := range args {
		if lastLetter == z {
			return true
		}
	}
	return false
}

func FmtSvg(svg string) string {

	if isEmpty := (svg == ""); isEmpty {
		return ""
	}

	trimmedExt := strings.TrimSuffix(svg, ".svg")
	if hasExt := (trimmedExt != svg); hasExt {
		return svg
	}

	return concat(svg, ".svg")
}

func fmtDir(dir string) string {

	if dir == "" {
		return ""
	}

	//Pwd := getPwd()
	formatted := dir

	if IsFirstLetter(dir, "/", "~") == false {
		formatted = concat(Pwd, "/", dir)
	}

	if IsLastLetter(dir, "/") == false {
		formatted = concat(formatted, "/")
	}

	if dir == "." {
		formatted = Pwd
	}

	return formatted
}

func fmtCopy(s string) string {
	if strings.HasSuffix(s, ".svg") {
		return strings.Replace(s, ".svg", "-copy.svg", 1)
	}
	return concat(s, "-copy.svg")
}

func fmtDest(path string) (out string) {
	out = strings.Replace(path, SrcDir, DstDir, 1)
	out = strings.Replace(out, "//", "/", -1)
	return
}

func replace(s string) (replaced string) {
	re, replacement := findReplacements(s)
	if replacement == "" {
		return
	}
	replaced = re.ReplaceAllString(s, replacement)
	return
}

func findReplacements(s string) (re *regexp.Regexp, replacement string) {
	if DoAddNew && !strings.Contains(s, ToFind) {
		re = regexp.MustCompile("(<svg )")
		replacement = concat(`${1}fill="`, ToReplace, `" `)
		return
	}

	if ToFind == "" {
		return
	}

	re = regexp.MustCompile(ToFind)
	replacement = ToReplace
	return
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
