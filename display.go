package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	Red      = color.New(color.FgRed)
	Blue     = color.New(color.FgBlue)
	Green    = color.New(color.FgGreen)
	Magenta  = color.New(color.FgMagenta)
	White    = color.New(color.FgWhite)
	Black    = color.New(color.FgBlack)
	BRed     = color.New(color.FgRed, color.Bold)
	BBlue    = color.New(color.FgBlue, color.Bold)
	BGreen   = color.New(color.FgGreen, color.Bold)
	BMagenta = color.New(color.FgMagenta, color.Bold)
	BWhite   = color.New(color.Bold, color.FgWhite)
	BBlack   = color.New(color.Bold, color.FgBlack)

	Log    = fmt.Println
	LogErr = fmt.Println
)

func LogNoop(x ...interface{}) (int, error) {
	return 0, nil
}

func report() {

	if DoShutUp {
		return
	}

	if DoRecursive {
		fmt.Printf("Edited %d svg images in %v\n", TotalEdited, DstDir)
		return
	}

	if TotalEdited == 0 {
		fmt.Printf("Failed to edit %v\n", DstSvg)
		return
	}

	fmt.Printf("Successfully edited %v\n", DstSvg)
}

func printHelp() {
	defer os.Exit(0)
	fmt.Printf(
		`svg-edit [options] <original file/directory> <new file/directory>
    -o="":
             (old) string in svg file to replace
    -n="":
             (new) string to replace old string with
    -a=false:
             (add) add fill color of 'new string' for files without one
    -c=false:
             (color) [same as '-a']
    -r=false:
             (recursive) edit svg files beneath the specified folder
    -q=false:
             (quiet) don't list edited files
    -Q=false:
             (QUIET) don't show any output%v`, "\n",
	)
}
