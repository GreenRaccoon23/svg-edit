package main

import (
	"os"
	"strings"
)

func parseArgs(boolFlagVars map[string]*bool, stringFlagVars map[string]*string, noFlagVars []*string) (extraArgs []string) {

	if _helpRequested() {
		printHelp()
	}

	a := _argParser{
		args:           os.Args,
		boolFlagVars:   boolFlagVars,
		stringFlagVars: stringFlagVars,
		noFlagVars:     noFlagVars,
	}
	defer a._reset()

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
	args     []string
	_iEndArg int

	boolFlagVars   map[string]*bool
	stringFlagVars map[string]*string
	noFlagVars     []*string

	_argsNotFlagged []string
}

func (a *_argParser) _reset() {
	go func() { a.args = nil }()
	go func() { a.boolFlagVars = nil }()
	go func() { a.stringFlagVars = nil }()
	go func() { a.noFlagVars = nil }()
	go func() { a._argsNotFlagged = nil }()
}

func (a *_argParser) _parseArgs() []string {

	args := a.args
	iEnd := len(args) - 1
	a._iEndArg = iEnd

	for i := 1; i <= iEnd; i++ {
		arg := args[i]

		if isFlag := a._parseArg(arg, &i); !isFlag {
			a._argsNotFlagged = append(a._argsNotFlagged, arg)
		}
	}

	return a._setNoFlags()
}

func (a *_argParser) _setNoFlags() (extraArgs []string) {

	argsNotFlagged := a._argsNotFlagged
	lenArgsNotFlagged := len(argsNotFlagged)

	noFlagVars := a.noFlagVars
	lenNoFlagVars := len(noFlagVars)

	iEnd := lenNoFlagVars - 1
	if enoughArgs := (lenArgsNotFlagged > lenNoFlagVars); !enoughArgs {
		iEnd = lenArgsNotFlagged - 1
	}

	for i := 0; i <= iEnd; i++ {
		*noFlagVars[i] = argsNotFlagged[i]
	}

	extraArgs = cut(argsNotFlagged, iEnd+1, -1)
	return
}

func (a *_argParser) _parseArg(arg string, i *int) bool {

	if beginsWithHyphen := (string(arg[0]) == "-"); !beginsWithHyphen {
		return false
	}

	argTrimmed := strings.TrimLeft(arg, "-")

	if hasBoolFlags := a._setBoolFlags(argTrimmed); hasBoolFlags {
		return true
	}

	if isLastArg := (*i == a._iEndArg); isLastArg {
		return false
	}

	if isStringFlag := a._setStringFlag(argTrimmed, i); isStringFlag {
		return true
	}

	return false
}

func (a *_argParser) _setBoolFlags(argTrimmed string) (hasBoolFlags bool) {

	iEnd := len(argTrimmed) - 1
	for i := 0; i <= iEnd; i++ {
		c := string(argTrimmed[i])

		if isBoolFlag := (a.boolFlagVars[c] != nil); isBoolFlag {
			*(a.boolFlagVars[c]) = true
			hasBoolFlags = true
		}
	}

	return
}

func (a *_argParser) _setStringFlag(argTrimmed string, i *int) (isStringFlag bool) {

	if isStringFlag = (a.stringFlagVars[argTrimmed] != nil); isStringFlag {
		*i++
		nextArg := a.args[*i]
		*(a.stringFlagVars[argTrimmed]) = nextArg
	}

	return
}
