package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	doRecursive bool
	doAddNew    bool
	doCopy      bool
	doQuiet     bool
	doShutUp    bool

	oldString string
	newString string

	origSvg    string
	copySvg    string
	pwd        string = Pwd()
	Root       string
	destDir    string
	numChanged int
	svgFile    *os.File
)

func init() {

	boolFlags := map[string]*bool{
		"r": &doRecursive,
		"a": &doAddNew,
		"c": &doCopy,
		"q": &doQuiet,
		"Q": &doShutUp,
	}

	stringFlags := map[string]*string{
		"o": &oldString,
		"n": &newString,
		"d": &Root,
	}

	extraArgs := parseArgs(boolFlags, stringFlags)
	fmt.Println("extraArgs:", extraArgs)
	// argsAnalyse()

	_printFlags()
	os.Exit(0)
}

func main() {
	defer ColorUnset()

	MakeDir(destDir)
	checkMethod()
	report()
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

	origSvg = FmtSvg(o)
	copySvg = FmtSvg(c)
	Root = filepath.Dir(origSvg)
	destDir = filepath.Dir(copySvg)
}

func argsAnalyseRecursive(args []string) {
	numArgs := len(args)

	root := args[numArgs-2]
	if numArgs < 3 {
		root = pwd
	}
	dest := args[numArgs-1]

	Root = FmtDir(root)
	destDir = FmtDir(dest)
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
	origDestination := destDir
	for k, v := range MaterialDesign {
		newString = v
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
	err := filepath.Walk(Root, WalkReplace)
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
	fmt.Println("d:", "Root:", Root)
}
