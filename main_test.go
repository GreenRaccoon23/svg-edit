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

	ReToFind *regexp.Regexp

	ToFill string
)

func TestInit(t *testing.T) {

	ToFind = "green"
	ToReplace = "cyan"

	_setFindReplace()

	ReToFind = regexp.MustCompile(ToFind)

	ToFill = string(ToFillBytes)
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
		ReToFind.ReplaceAll(TestFileContentBytes, ToReplaceBytes)
	}
}

func BenchmarkRegexpReplaceAddFill(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReAddFill.ReplaceAllString(TestFileContent, ToFill)
	}
}

func BenchmarkRegexpReplaceAddFillBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReAddFill.ReplaceAll(TestFileContentBytes, ToFillBytes)
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

// func TestPrintMaterialPalette(t *testing.T) {
// 	// showGroupNames := true
// 	showGroupNames := false
// 	// showGroups := true
// 	showGroups := false
// 	for groupName, group := range MaterialPalette {
// 		if showGroupNames {
// 			t.Log(groupName)
// 		}
// 		if showGroups {
// 			t.Log(group)
// 		}
// 	}
// }

func TestGetGroupNameShade(t *testing.T) {

	slc := []string{
		"deeporange50",
		"deeporange500",
		"deeporange900",
		"deeporangeA100",
		"deeporangeA200",
		"deeporangeA400",
		"deeporangeA700",
		"deeporange-50",
		"deeporange:500",
		"deeporange_900",
		"deeporange A100",
		"deeporange-A200",
		"deeporange:A400",
		"deeporange_A700",
		"deeporange000",
		"deeporange1000",
		"deeporangeA000",
		"deeporangeA300",
		"deeporangeA1000",
		"deeporange:000",
		"deeporange-1000",
		"deeporange_A000",
		"deeporange A300",
		"deeporange-A1000",
		"deeporange",
	}

	for _, s := range slc {

		groupName, shade, wasFound := getGroupNameShade(s)
		if !wasFound {
			t.Log("false:", s)
		} else {
			t.Log("true:", groupName, " ", shade)
		}
	}
}
