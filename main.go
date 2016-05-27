package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var (
	DoRecursive bool
	DoAddNew    bool
	DoCopy      bool
	DoQuiet     bool
	DoShutUp    bool

	OldString string
	NewString string

	SrcFileOrDir string
	DstFileOrDir string

	SrcSvg      string
	DstSvg      string
	Pwd         string = getPwd()
	SrcDir      string
	DstDir      string
	TotalEdited int
)

func init() {

	boolFlagVars := map[string]*bool{
		"r": &DoRecursive,
		"a": &DoAddNew,
		"c": &DoCopy,
		"q": &DoQuiet,
		"Q": &DoShutUp,
	}

	stringFlagVars := map[string]*string{
		"o": &OldString,
		"n": &NewString,
		"d": &SrcDir,
	}

	noFlagVars := []*string{
		&SrcFileOrDir,
		&DstFileOrDir,
	}

	parseArgs(boolFlagVars, stringFlagVars, noFlagVars)
	_formatGlobalVars()
	// argsAnalyse()

	_printFlags()
	os.Exit(0)
}

func main() {
	defer color.Unset()

	mkDir(DstDir)
	checkMethod()
	report()
}

func _formatGlobalVars() {

	if DstFileOrDir == "" {
		DstFileOrDir = SrcFileOrDir
	}

	switch DoRecursive {

	case true:
		SrcDir = FmtDir(SrcFileOrDir)
		DstDir = FmtDir(DstFileOrDir)

	case false:
		SrcSvg = FmtSvg(SrcFileOrDir)
		DstSvg = FmtSvg(DstFileOrDir)

		SrcDir = filepath.Dir(SrcSvg)
		DstDir = filepath.Dir(DstSvg)
	}
}

// func argsAnalyse() {
// 	switch os.Args[1] {
// 	case "h", "-h", "help", "-help", "--help":
// 		printHelp()
// 	}

// 	flag.StringVar(&OldString, "o", "", "(old) string in svg file to replace")
// 	flag.StringVar(&NewString, "n", "", "(new) string in svg file to replace with")
// 	flag.BoolVar(&DoAddNew, "a", false, "add new string if the old one does not exist")
// 	flag.BoolVar(&DoCopy, "c", false, "Make a copy instead of editing file")
// 	flag.BoolVar(&DoRecursive, "r", false, "walk recursively down to the bottom of the directory")
// 	flag.BoolVar(&DoQuiet, "q", false, "don't list edited files")
// 	flag.BoolVar(&DoShutUp, "Q", false, "don't show any output")
// 	flag.Parse()

// 	args := Filter(os.Args,
// 		"-o", OldString,
// 		"-n", NewString,
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
// 	analyseColor()
// 	return
// }

// func argsAnalyseSingle(args []string) {
// 	numArgs := len(args)

// 	o := args[numArgs-1]
// 	c := o
// 	if DoCopy && numArgs > 2 {
// 		o = args[numArgs-2]
// 		c = args[numArgs-1]
// 	}
// 	if DoCopy && numArgs < 3 {
// 		o = args[numArgs-1]
// 		c = fmtCopy(o)
// 	}

// 	SrcSvg = FmtSvg(o)
// 	DstSvg = FmtSvg(c)
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

// 	SrcDir = FmtDir(root)
// 	DstDir = FmtDir(dest)
// }

func analyseColor() {
	checkOld := strings.ToLower(OldString)
	checkNew := strings.ToLower(NewString)

	if IsKeyInMap(MaterialDesign, checkOld) {
		OldString = MaterialDesign[checkOld]
	}
	if IsKeyInMap(MaterialDesign, checkNew) {
		NewString = MaterialDesign[checkNew]
	}
}

func checkMethod() {
	if DoRecursive == false {
		editSingle()
		return
	}
	Progress("Editing all svg files recursively...")
	editRecursive()
}

/*func editLollipop() {
	origDestination := DstDir
	for k, v := range MaterialDesign {
		NewString = v
		DstDir = Concat(origDestination, k, v, "/")
		if DoRecursive {
			editRecursive()
		} else {
			editSingle()
		}
	}
}
*/
func editSingle() {
	in := SrcSvg
	out := DstSvg
	edit(in, out)
}

func editRecursive() {
	err := filepath.Walk(SrcDir, WalkReplace)
	LogErr(err)
}

func _printFlags() {
	fmt.Println("r:", "DoRecursive:", DoRecursive)
	fmt.Println("a:", "DoAddNew:", DoAddNew)
	fmt.Println("c:", "DoCopy:", DoCopy)
	fmt.Println("q:", "DoQuiet:", DoQuiet)
	fmt.Println("Q:", "DoShutUp:", DoShutUp)

	fmt.Println("o:", "OldString:", OldString)
	fmt.Println("n:", "NewString:", NewString)
	fmt.Println("d:", "SrcDir:", SrcDir)

	fmt.Println("_:", "SrcFileOrDir:", SrcFileOrDir)
	fmt.Println("_:", "DstFileOrDir:", DstFileOrDir)
}
