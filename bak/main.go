package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
)

var (
	doRecursive bool
	doAddNew    bool
	doCopy      bool
	quiet       bool
	shutUp      bool
	sOld        string
	sNew        string
	origSvg     string
	copySvg     string
	pwd         string
	rootDir     string
	destDir     string
	numChanged  int
	svgFile     *os.File
)

func init() {
	pwd = Pwd()
	argsAnalyse()
}

func argsAnalyse() {
	switch os.Args[1] {
	case "h", "-h", "help", "-help", "--help":
		printHelp()
	}

	flag.StringVar(&sOld, "o", "", "(old) string in svg file to replace")
	flag.StringVar(&sNew, "n", "", "(new) string in svg file to replace with")
	flag.BoolVar(&doAddNew, "a", false, "add new string if the old one does not exist")
	flag.BoolVar(&doCopy, "c", false, "Make a copy instead of editing file")
	flag.BoolVar(&doRecursive, "r", false, "walk recursively down to the bottom of the directory")
	flag.BoolVar(&quiet, "q", false, "don't list edited files")
	flag.BoolVar(&shutUp, "Q", false, "don't show any output")
	flag.Parse()

	args := Filter(os.Args,
		"-o", sOld,
		"-n", sNew,
		"-a",
		"-c",
		"-r",
		"-q",
		"-Q",
	)
	switch doRecursive {
	case true:
		argsAnalyseRecursive(args)
	case false:
		argsAnalyseSingle(args)
	}
	analyseColor()
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

	origSvg = FmtSvg(o)
	copySvg = FmtSvg(c)
	rootDir = filepath.Dir(origSvg)
	destDir = filepath.Dir(copySvg)
}

func argsAnalyseRecursive(args []string) {
	numArgs := len(args)

	root := args[numArgs-2]
	if numArgs < 3 {
		root = pwd
	}
	dest := args[numArgs-1]

	rootDir = FmtDir(root)
	destDir = FmtDir(dest)
}

func analyseColor() {
	checkOld := strings.ToLower(sOld)
	checkNew := strings.ToLower(sNew)

	if IsKeyInMap(MaterialDesign, checkOld) {
		sOld = MaterialDesign[checkOld]
	}
	if IsKeyInMap(MaterialDesign, checkNew) {
		sNew = MaterialDesign[checkNew]
	}
}

func main() {
	defer ColUn()
	MakeDir(destDir)
	checkMethod()
	report()
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
	origDestination := destDir
	for k, v := range MaterialDesign {
		sNew = v
		destDir = Concat(origDestination, k, v, "/")
		if doRecursive {
			editRecursive()
		} else {
			editSingle()
		}
	}
}
*/
func editSingle() {
	in := origSvg
	out := copySvg
	edit(in, out)
}

func editRecursive() {
	err := filepath.Walk(rootDir, WalkReplace)
	LogErr(err)
}
