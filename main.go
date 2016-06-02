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

	ReAddFill *regexp.Regexp = regexp.MustCompile("(<svg)")
	ReToFind  *regexp.Regexp

	ToFill string

	ToFindBytes    []byte
	ToReplaceBytes []byte

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
		Log(err)
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
		SrcSvg = addExt(SrcFileOrDir, ".svg")
		DstSvg = addExt(DstFileOrDir, ".svg")
		SrcDstSame = (SrcSvg == DstSvg)

		SrcDir = filepath.Dir(SrcSvg)
		DstDir = filepath.Dir(DstSvg)
	}
}

func _setFindReplace() {

	ToFind = toMaterial(ToFind)
	ToReplace = toMaterial(ToReplace)

	ReToFind = regexp.MustCompile(ToFind)
	ToFill = concat(`${1} fill="`, ToReplace, `"`)

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

	if ReAddFill == nil {
		return fmt.Errorf("Fatal program bug! ReAddFill not set")
	}

	if ReToFind == nil {
		return fmt.Errorf("Fatal program bug! ReToFind not set")
	}

	return nil
}

func _edit() error {

	switch DoRecursive {

	case true:
		return filepath.Walk(SrcDir, walkReplace)

	case false:
		return editFileFromPath(DstSvg, SrcSvg)
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

// func argsAnalyse() {
// 	switch os.Args[1] {
// 	case "h", "-h", "help", "-help", "--help":
// 		printHelp()
// 	}

// 	flag.StringVar(&ToFind, "o", "", "(old) string in svg file to replace")
// 	flag.StringVar(&ToReplace, "n", "", "(new) string in svg file to replace with")
// 	flag.BoolVar(&DoAddFill, "a", false, "add new string if the old one does not exist")
// 	flag.BoolVar(&DoColor, "c", false, "Make a copy instead of editing file")
// 	flag.BoolVar(&DoRecursive, "r", false, "walk recursively down to the bottom of the directory")
// 	flag.BoolVar(&DoQuiet, "q", false, "don't list edited files")
// 	flag.BoolVar(&DoShutUp, "Q", false, "don't show any output")
// 	flag.Parse()

// 	args := Filter(os.Args,
// 		"-o", ToFind,
// 		"-n", ToReplace,
// 		"-a",
// 		"-c",
// 		"-r",
// 		"-q",
// 		"-Q",
// 	)

// 	fmt.Println("args:", args)
// 	switch DoRecursive {
// 	case true:
// 		argsAnalyseRecursive(args)
// 	case false:
// 		argsAnalyseSingle(args)
// 	}
// 	setFindReplace()
// 	return
// }

// func argsAnalyseSingle(args []string) {
// 	numArgs := len(args)

// 	o := args[numArgs-1]
// 	c := o
// 	if DoColor && numArgs > 2 {
// 		o = args[numArgs-2]
// 		c = args[numArgs-1]
// 	}
// 	if DoColor && numArgs < 3 {
// 		o = args[numArgs-1]
// 		c = fmtCopy(o)
// 	}

// 	SrcSvg = addExt(o, ".svg")
// 	DstSvg = addExt(c, ".svg")
// 	SrcDir = filepath.Dir(SrcSvg)
// 	DstDir = filepath.Dir(DstSvg)
// }

// func argsAnalyseRecursive(args []string) {
// 	numArgs := len(args)

// 	root := args[numArgs-2]
// 	if numArgs < 3 {
// 		root = Pwd
// 	}
// 	dest := args[numArgs-1]

// 	SrcDir = fmtDir(root)
// 	DstDir = fmtDir(dest)
// }
