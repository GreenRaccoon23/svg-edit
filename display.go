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
)

func Progress(current ...string) {
	if doShutUp || doQuiet {
		return
	}
	for _, c := range current {
		fmt.Printf("%v ", c)
		fmt.Println()
	}
}

func report() {
	if doShutUp {
		return
	}
	fmt.Printf("Edited %d svg images in %v\n", numChanged, DstDir)
}

func printHelp() {
	defer os.Exit(0)
	fmt.Printf(
		"%v\n  %v\n%v\n  %v\n%v\n  %v\n%v\n  %v\n%v\n  %v\n%v\n  %v\n%v\n  %v\n%v\n",
		"svg-edit <options> <original file/directory> <new file/directory>",
		`-o="": (old)`,
		"      string in svg file to replace",
		`-n="": (new)`,
		"      string to replace old string with",
		"-a=false: (add)",
		"      add 'new string' even if 'old string' does not exist",
		"-c=false: (copy)",
		"      create 'new string' even if 'old string' does not exist",
		"-r=false: (recursive)",
		"      edit svg files beneath the specified folder",
		"-q=false: (doQuiet)",
		"      don't list edited files",
		"-Q=false: (Quiet)",
		"      don't show any output",
	)
}
