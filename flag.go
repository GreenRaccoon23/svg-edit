package main

import (
	"os"
	"strings"
)

func parseArgs(boolFlags map[string]*bool, stringFlags map[string]*string) (extraArgs []string) {

	if _helpRequested() {
		printHelp()
	}

	a := _argParser{
		args:        os.Args,
		boolFlags:   boolFlags,
		stringFlags: stringFlags,
	}

	extraArgs = a._parseArgs()
	return
}

func _helpRequested() bool {

	if len(os.Args) < 2 {
		return true
	}

	switch os.Args[1] {
	case "-h", "h", "help", "--help", "-H", "H", "HELP", "--HELP", "-help", "--h", "--H":
		return true
	}

	return false
}

type _argParser struct {
	args   []string
	_iLast int

	boolFlags   map[string]*bool
	stringFlags map[string]*string
}

func (a *_argParser) _parseArgs() (extraArgs []string) {

	args := a.args
	iLast := len(args) - 1
	a._iLast = iLast

	for i := 1; i <= iLast; i++ {
		arg := args[i]

		if isFlag := a._parseArg(arg, &i); !isFlag {
			extraArgs = append(extraArgs, arg)
		}
	}

	return
}

func (a *_argParser) _parseArg(arg string, i *int) bool {

	if beginsWithHyphen := (string(arg[0]) == "-"); !beginsWithHyphen {
		return false
	}

	argTrimmed := strings.TrimLeft(arg, "-")

	if hasBoolFlags := a._checkBoolFlags(argTrimmed); hasBoolFlags {
		return true
	}

	if isLastArg := (*i == a._iLast); isLastArg {
		return false
	}

	if isStringFlag := a._checkStringFlags(argTrimmed, i); isStringFlag {
		return true
	}

	return false
}

func (a *_argParser) _checkBoolFlags(argTrimmed string) (hasBoolFlags bool) {

	iEnd := len(argTrimmed) - 1
	for i := 0; i <= iEnd; i++ {
		c := string(argTrimmed[i])

		if isBoolFlag := (a.boolFlags[c] != nil); isBoolFlag {
			*(a.boolFlags[c]) = true
			hasBoolFlags = true
		}
	}

	return
}

func (a *_argParser) _checkStringFlags(argTrimmed string, i *int) (isStringFlag bool) {

	if isStringFlag = (a.stringFlags[argTrimmed] != nil); isStringFlag {
		*i++
		nextArg := a.args[*i]
		*(a.stringFlags[argTrimmed]) = nextArg
	}

	return
}
