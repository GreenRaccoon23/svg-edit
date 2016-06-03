package main

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"

	"github.com/fatih/color"
)

var (
	Pwd string = getPwd()

	DoRecursive bool
	DoAddFill   bool
	DoQuiet     bool
	DoShutUp    bool

	ToFind    string
	ToReplace string

	SrcFileOrDir string
	DstFileOrDir string

	SrcSvg string
	DstSvg string
	SrcDir string
	DstDir string

	SrcDstSame bool

	ToFindBytes    []byte
	ToReplaceBytes []byte

	ReAddFill *regexp.Regexp = regexp.MustCompile("(<svg)")

	ToFillBytes []byte

	TotalEdited int
)

func init() {

	boolFlagVars := map[string]*bool{
		"r": &DoRecursive,
		"a": &DoAddFill,
		"c": &DoAddFill,
		"q": &DoQuiet,
		"Q": &DoShutUp,
	}

	stringFlagVars := map[string]*string{
		"o": &ToFind,
		"n": &ToReplace,
		"d": &SrcDir,
	}

	noFlagVars := []*string{
		&SrcFileOrDir,
		&DstFileOrDir,
	}

	parseArgs(boolFlagVars, stringFlagVars, noFlagVars)
	_setLogger()
	_setSrcDst()
	_setFindReplace()

	_verifyGlobalVars()

	// _printFlags()
	// os.Exit(0)
}

func main() {
	defer color.Unset()

	if err := mkDir(DstDir); err != nil {
		log.Fatal(err)
	}

	if err := _edit(); err != nil {
		LogErr(err)
	}

	report()
}

func _setLogger() {
	if DoQuiet {
		Log = LogNoop
	}

	if DoShutUp {
		Log = LogNoop
		LogErr = LogNoop
	}
}

func _setSrcDst() {

	if DstFileOrDir == "" {
		DstFileOrDir = SrcFileOrDir
	}

	switch DoRecursive {

	case true:
		SrcDir = fmtDir(SrcFileOrDir)
		DstDir = fmtDir(DstFileOrDir)
		SrcDstSame = (SrcDir == DstDir)

	case false:
		SrcSvg = SrcFileOrDir
		DstSvg = DstFileOrDir

		SrcDstSame = (SrcSvg == DstSvg)

		SrcDir = filepath.Dir(SrcSvg)
		DstDir = filepath.Dir(DstSvg)
	}
}

func _setFindReplace() {

	if hex := getHex(ToFind); hex != "" {
		ToFind = hex
	}

	if hex := getHex(ToReplace); hex != "" {
		ToReplace = hex
	}

	ToFindBytes = []byte(ToFind)
	ToReplaceBytes = []byte(ToReplace)

	ToFillBytes = toBytes(`${1} fill="`, ToReplace, `"`)
}

func _verifyGlobalVars() error {

	if ToFind == "" {
		return fmt.Errorf("-o paramater required")
	}

	if ToReplace == "" {
		return fmt.Errorf("-n paramater required")
	}

	if SrcFileOrDir == "" {
		return fmt.Errorf("<original file/directory> required")
	}

	if DstFileOrDir == "" {
		// return fmt.Errorf("<new file/directory> required")
		return fmt.Errorf("Fatal program bug! DstFileOrDir not set")
	}

	if SrcDir == "" {
		return fmt.Errorf("Fatal program bug! SrcDir not set")
	}

	if DstDir == "" {
		return fmt.Errorf("Fatal program bug! DstDir not set")
	}

	if SrcSvg == "" && !DoRecursive {
		return fmt.Errorf("Fatal program bug! SrcSvg not set")
	}

	if DstSvg == "" && !DoRecursive {
		return fmt.Errorf("Fatal program bug! DstSvg not set")
	}

	if len(ToFindBytes) == 0 {
		return fmt.Errorf("Fatal program bug! ToFindBytes not set")
	}

	if len(ToReplaceBytes) == 0 {
		return fmt.Errorf("Fatal program bug! ToReplaceBytes not set")
	}

	if len(ToFillBytes) <= 12 && ToReplace != "" {
		return fmt.Errorf("Fatal program bug! ToFillBytes not set")
	}

	if ReAddFill == nil {
		return fmt.Errorf("Fatal program bug! ReAddFill not set")
	}

	return nil
}

func _edit() error {

	switch DoRecursive {
	case true:
		return editRecursive()
	case false:
		return editOne()
	}

	return nil //why does Go require this?
}

func _printFlags() {
	fmt.Println("r:", "DoRecursive:", DoRecursive)
	fmt.Println("a:", "DoAddFill:", DoAddFill)
	fmt.Println("q:", "DoQuiet:", DoQuiet)
	fmt.Println("Q:", "DoShutUp:", DoShutUp)

	fmt.Println("o:", "ToFind:", ToFind)
	fmt.Println("n:", "ToReplace:", ToReplace)
	fmt.Println("d:", "SrcDir:", SrcDir)

	fmt.Println("_:", "SrcFileOrDir:", SrcFileOrDir)
	fmt.Println("_:", "DstFileOrDir:", DstFileOrDir)

	fmt.Println("_:", "DstDir:", DstDir)
	fmt.Println("_:", "SrcSvg:", SrcSvg)
	fmt.Println("_:", "DstSvg:", DstSvg)
	fmt.Println("_:", "ToFindBytes:", string(ToFindBytes))
	fmt.Println("_:", "ToReplaceBytes:", string(ToReplaceBytes))
}

/*func editLollipop() {
	origDestination := DstDir
	for k, v := range MaterialDesign {
		ToReplace = v
		DstDir = concat(origDestination, k, v, "/")
		if DoRecursive {
			editRecursive()
		} else {
			editSingle()
		}
	}
}
*/
