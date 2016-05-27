package main

import (
	"fmt"
	"log"
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

	ToFind    string
	ToReplace string

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
		"o": &ToFind,
		"n": &ToReplace,
		"d": &SrcDir,
	}

	noFlagVars := []*string{
		&SrcFileOrDir,
		&DstFileOrDir,
	}

	parseArgs(boolFlagVars, stringFlagVars, noFlagVars)
	_setSrcDst()
	_setFindReplace()

	_printFlags()
	os.Exit(0)
}

func main() {
	defer color.Unset()

	if err := mkDir(DstDir); err != nil {
		log.Fatal(err)
	}
	_edit()
	report()
}

func _setFindReplace() {

	oldString := strings.ToLower(ToFind)
	if MaterialDesign[oldString] != "" {
		ToFind = MaterialDesign[oldString]
	}

	newString := strings.ToLower(ToReplace)
	if MaterialDesign[oldString] != "" {
		ToReplace = MaterialDesign[newString]
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

	case false:
		SrcSvg = addExt(SrcFileOrDir, ".svg")
		DstSvg = addExt(DstFileOrDir, ".svg")

		SrcDir = filepath.Dir(SrcSvg)
		DstDir = filepath.Dir(DstSvg)
	}
}

func _edit() {

	switch DoRecursive {

	case true:
		Progress("Editing all svg files recursively...")
		err := filepath.Walk(SrcDir, walkReplace)
		LogErr(err)

	case false:
		editFileFromPath(DstSvg, SrcSvg)
	}
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
// 	flag.BoolVar(&DoAddNew, "a", false, "add new string if the old one does not exist")
// 	flag.BoolVar(&DoCopy, "c", false, "Make a copy instead of editing file")
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
// 	if DoCopy && numArgs > 2 {
// 		o = args[numArgs-2]
// 		c = args[numArgs-1]
// 	}
// 	if DoCopy && numArgs < 3 {
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

func _printFlags() {
	fmt.Println("r:", "DoRecursive:", DoRecursive)
	fmt.Println("a:", "DoAddNew:", DoAddNew)
	fmt.Println("c:", "DoCopy:", DoCopy)
	fmt.Println("q:", "DoQuiet:", DoQuiet)
	fmt.Println("Q:", "DoShutUp:", DoShutUp)

	fmt.Println("o:", "ToFind:", ToFind)
	fmt.Println("n:", "ToReplace:", ToReplace)
	fmt.Println("d:", "SrcDir:", SrcDir)

	fmt.Println("_:", "SrcFileOrDir:", SrcFileOrDir)
	fmt.Println("_:", "DstFileOrDir:", DstFileOrDir)
}
