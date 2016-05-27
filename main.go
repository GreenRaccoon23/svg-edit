package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var (
	doRecursive bool
	doAddNew    bool
	doCopy      bool
	doQuiet     bool
	doShutUp    bool

	oldString string
	newString string

	srcFileOrDir string
	dstFileOrDir string

	SrcSvg     string
	DstSvg     string
	Pwd        string = getPwd()
	SrcDir     string
	DstDir     string
	numChanged int
	svgFile    *os.File
)

func init() {

	boolFlagVars := map[string]*bool{
		"r": &doRecursive,
		"a": &doAddNew,
		"c": &doCopy,
		"q": &doQuiet,
		"Q": &doShutUp,
	}

	stringFlagVars := map[string]*string{
		"o": &oldString,
		"n": &newString,
		"d": &SrcDir,
	}

	noFlagVars := []*string{
		&srcFileOrDir,
		&dstFileOrDir,
	}

	parseArgs(boolFlagVars, stringFlagVars, noFlagVars)
	_formatGlobalVars()
	// argsAnalyse()

	_printFlags()
	os.Exit(0)
}

func main() {
	defer color.Unset()

	MakeDir(DstDir)
	checkMethod()
	report()
}

func _formatGlobalVars() {

	if dstFileOrDir == "" {
		dstFileOrDir = srcFileOrDir
	}

	switch doRecursive {

	case true:
		SrcDir = FmtDir(srcFileOrDir)
		DstDir = FmtDir(dstFileOrDir)

	case false:
		SrcSvg = FmtSvg(srcFileOrDir)
		DstSvg = FmtSvg(dstFileOrDir)

		SrcDir = filepath.Dir(SrcSvg)
		DstDir = filepath.Dir(DstSvg)
	}
}

func argsAnalyse() {
	switch os.Args[1] {
	case "h", "-h", "help", "-help", "--help":
		printHelp()
	}

	flag.StringVar(&oldString, "o", "", "(old) string in svg file to replace")
	flag.StringVar(&newString, "n", "", "(new) string in svg file to replace with")
	flag.BoolVar(&doAddNew, "a", false, "add new string if the old one does not exist")
	flag.BoolVar(&doCopy, "c", false, "Make a copy instead of editing file")
	flag.BoolVar(&doRecursive, "r", false, "walk recursively down to the bottom of the directory")
	flag.BoolVar(&doQuiet, "q", false, "don't list edited files")
	flag.BoolVar(&doShutUp, "Q", false, "don't show any output")
	flag.Parse()

	args := FilterOut(os.Args,
		"-o", oldString,
		"-n", newString,
		"-a",
		"-c",
		"-r",
		"-q",
		"-Q",
	)

	fmt.Println("args:", args)
	switch doRecursive {
	case true:
		argsAnalyseRecursive(args)
	case false:
		argsAnalyseSingle(args)
	}
	analyseColor()
	return
}

func FilterOut(slc []string, args ...string) (filtered []string) {
	for _, s := range slc {
		if !SlcContains(args, s) {
			filtered = append(filtered, s)
		}
	}
	return
}

func FilterOut2(slc []string, args ...string) (filtered []string) {

	lenSlc := len(slc)
	var wg sync.WaitGroup
	wg.Add(lenSlc)

	for _, s := range slc {
		go func(s string) {
			defer wg.Done()

			if !SlcContains(args, s) {
				filtered = append(filtered, s)
			}
		}(s)
	}

	wg.Wait()

	return
}

func argsAnalyseSingle(args []string) {
	numArgs := len(args)

	o := args[numArgs-1]
	c := o
	if doCopy && numArgs > 2 {
		o = args[numArgs-2]
		c = args[numArgs-1]
	}
	if doCopy && numArgs < 3 {
		o = args[numArgs-1]
		c = fmtCopy(o)
	}

	SrcSvg = FmtSvg(o)
	DstSvg = FmtSvg(c)
	SrcDir = filepath.Dir(SrcSvg)
	DstDir = filepath.Dir(DstSvg)
}

func argsAnalyseRecursive(args []string) {
	numArgs := len(args)

	root := args[numArgs-2]
	if numArgs < 3 {
		root = Pwd
	}
	dest := args[numArgs-1]

	SrcDir = FmtDir(root)
	DstDir = FmtDir(dest)
}

func analyseColor() {
	checkOld := strings.ToLower(oldString)
	checkNew := strings.ToLower(newString)

	if IsKeyInMap(MaterialDesign, checkOld) {
		oldString = MaterialDesign[checkOld]
	}
	if IsKeyInMap(MaterialDesign, checkNew) {
		newString = MaterialDesign[checkNew]
	}
}

func checkMethod() {
	if doRecursive == false {
		editSingle()
		return
	}
	Progress("Editing all svg files recursively...")
	editRecursive()
}

/*func editLollipop() {
	origDestination := DstDir
	for k, v := range MaterialDesign {
		newString = v
		DstDir = Concat(origDestination, k, v, "/")
		if doRecursive {
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
	fmt.Println("r:", "doRecursive:", doRecursive)
	fmt.Println("a:", "doAddNew:", doAddNew)
	fmt.Println("c:", "doCopy:", doCopy)
	fmt.Println("q:", "doQuiet:", doQuiet)
	fmt.Println("Q:", "doShutUp:", doShutUp)

	fmt.Println("o:", "oldString:", oldString)
	fmt.Println("n:", "newString:", newString)
	fmt.Println("d:", "SrcDir:", SrcDir)

	fmt.Println("_:", "srcFileOrDir:", srcFileOrDir)
	fmt.Println("_:", "dstFileOrDir:", dstFileOrDir)
}
