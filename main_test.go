package main

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

var (
	TestFileContent string = `<svg fill="#4CAF50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48"><path d="m40 10h-32c-2.21 0-3.98 1.79-3.98 4l-.02 20c0 2.21 1.79 4 4 4h32c2.21 0 4-1.79 4-4v-20c0-2.21-1.79-4-4-4m-18 6h4v4h-4v-4m0 6h4v4h-4v-4m-6-6h4v4h-4v-4m0 6h4v4h-4v-4m-2 4h-4v-4h4v4m0-6h-4v-4h4v4m18 14h-16v-4h16v4m0-8h-4v-4h4v4m0-6h-4v-4h4v4m6 6h-4v-4h4v4m0-6h-4v-4h4v4"/></svg>`

	TestFileContentBytes []byte = []byte(TestFileContent)

	ToFindBytes    []byte
	ToReplaceBytes []byte

	ReToFindBytes *regexp.Regexp
)

func TestInit(t *testing.T) {
	ToFind = "green"
	ToReplace = "cyan"
	_setFindReplace()

	ToFindBytes = []byte(ToFind)
	ToReplaceBytes = []byte(ToReplace)

	ReToFindBytes = regexp.MustCompile(ToFind)
}

func BenchmarkStringsReplaceAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.Replace(TestFileContent, ToFind, ToReplace, -1)
	}
}

func BenchmarkStringsReplaceOnce(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.Replace(TestFileContent, ToFind, ToReplace, 1)
	}
}

func BenchmarkRegexpReplace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReToFind.ReplaceAllString(TestFileContent, ToReplace)
	}
}

func BenchmarkBytesReplaceAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes.Replace(TestFileContentBytes, ToFindBytes, ToReplaceBytes, -1)
	}
}

func BenchmarkBytesReplaceOnce(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes.Replace(TestFileContentBytes, ToFindBytes, ToReplaceBytes, 1)
	}
}

func BenchmarkRegexpReplaceBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReToFindBytes.ReplaceAll(TestFileContentBytes, ToReplaceBytes)
	}
}

func BenchmarkRegexpReplaceAddFill(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReAddFill.ReplaceAllString(TestFileContent, ToFill)
	}
}

// func BenchmarkStringsTrimLeft(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		strings.TrimLeft("--abc123doremi", "-")
// 	}
// }

// func BenchmarkRegexpTrimLeft(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		re := regexp.MustCompile("^(-+)(.*)$")
// 		re.ReplaceAllString("--abc123doremi", "${2}")
// 	}
// }

// var re *regexp.Regexp = regexp.MustCompile("^(-+)(.*)$")

// func BenchmarkRegexpTrimLeftGlobal(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		re.ReplaceAllString("--abc123doremi", "${2}")
// 	}
// }

// func BenchmarkRegexpTrimLeftLocal(b *testing.B) {
// 	re := regexp.MustCompile("^(-+)(.*)$")
// 	for i := 0; i < b.N; i++ {
// 		re.ReplaceAllString("--abc123doremi", "${2}")
// 	}
// }
