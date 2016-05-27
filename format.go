package main

import (
	"bytes"
	"regexp"
	"strings"
)

var (
	buffer bytes.Buffer
)

func IsTrue(args ...bool) bool {
	for _, a := range args {
		if a {
			return true
		}
	}
	return false
}

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

func Concat(args ...string) string {
	return Str(args)
}

func Filter(slc []string, args ...string) (filtered []string) {
	sediment := Strain(slc, args...)
	for _, s := range slc {
		if s == "" {
			continue
		}
		if SlcContains(sediment, s) {
			continue
		}
		filtered = append(filtered, s)
	}
	return
}

func Strain(slc []string, args ...string) (sediment []string) {
	for _, s := range args {
		if s == "" {
			continue
		}
		if SlcContains(slc, s) == false {
			continue
		}
		sediment = append(sediment, s)
	}
	return
}

func SlcContainsLoose(slc []string, args ...string) bool {
	for _, s := range slc {
		for _, a := range args {
			//if s == a {
			if strings.Contains(s, a) {
				return true
			}
		}
	}
	return false
}

func SlcContains(slc []string, args ...string) bool {
	for _, s := range slc {
		for _, a := range args {
			//if s == a {
			if strings.Contains(s, a) {
				return true
			}
		}
	}
	return false
}

func IsMatch(s string, q string) bool {
	if s == q {
		return true
	}
	return false
}

func IsMatchAny(s string, args ...string) bool {
	for _, a := range args {
		if a == s {
			return true
		}
	}
	return false
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

func IsKeyInMap(m map[string]string, s string) bool {
	for k, _ := range m {
		if k == s {
			return true
		}
	}
	return false
}

func IsValueInMap(m map[string]string, s string) bool {
	for _, v := range m {
		if v == s {
			return true
		}
	}
	return false
}

func EndsWithAny(s string, args ...string) bool {
	for _, a := range args {
		if EndsWith(s, a) {
			return true
		}
	}
	return false
}

func EndsWith(s string, sub string) bool {
	subZ := sub[len(sub)-1]
	sZ := s[len(s)-1]
	if sZ != subZ {
		return false
	}
	subA := sub[0]
	target, exists := WhereIsByteInString(s, subA)
	if exists == false {
		return false
	}
	cutS := s[target:]
	for i := 0; i < len(cutS); i++ {
		if cutS[i] != sub[i] {
			return false
		}
	}
	return true
}

func IsByteInString(s string, b byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return true
		}
	}
	return false
}

func WhereIsByteInString(s string, b byte) (int, bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i, true
		}
	}
	return 0, false
}

func FmtSvg(svg string) string {

	if isEmpty := (svg == ""); isEmpty {
		return ""
	}

	trimmedExt := strings.trimSuffix(svg, ".svg")
	if hasExt := (trimmedExt != svg); hasExt {
		return svg
	}

	return Concat(svg, ".svg")
}

func FmtDir(dir string) string {

	if dir == "" {
		return ""
	}

	//Pwd := getPwd()
	formatted := dir

	if IsFirstLetter(dir, "/", "~") == false {
		formatted = Concat(Pwd, "/", dir)
	}

	if IsLastLetter(dir, "/") == false {
		formatted = Concat(formatted, "/")
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
	return Concat(s, "-copy.svg")
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
	if DoAddNew && !strings.Contains(s, OldString) {
		re = regexp.MustCompile("(<svg )")
		replacement = Concat(`${1}fill="`, NewString, `" `)
		return
	}

	if OldString == "" {
		return
	}

	re = regexp.MustCompile(OldString)
	replacement = NewString
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
