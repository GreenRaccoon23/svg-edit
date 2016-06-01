package main

import (
	"regexp"
	"strings"
)

// func replace(fileContent string) string {

//  re, replacement, needsToBeEdited := _getFindAndReplace(fileContent)
//  if !needsToBeEdited {
//      return fileContent
//  }

//  return re.ReplaceAllString(fileContent, replacement)
// }

// func _getFindAndReplace(fileContent string) (*regexp.Regexp, string, bool) {

//  if nothingToReplace := (!_containsToFind(fileContent)); nothingToReplace {

//      if shouldAddFill := (DoAddFill && !_hasFill(fileContent)); shouldAddFill {
//          return ReAddNew, ToAdd, true
//      }

//      return nil, "", false
//  }

//  return ReToFind, ToReplace, true
// }

func replace(fileContent string) string {

	replaced := strings.Replace(fileContent, ToFind, ToReplace, -1)

	if wasEdited := (replaced != fileContent); wasEdited {
		return replaced
	}

	if shouldAddFill := (DoAddFill && !_hasFill(fileContent)); !shouldAddFill {
		return replaced
	}

	return ReAddNew.ReplaceAllString(fileContent, ToAdd)
}

func _getFindAndReplace(fileContent string) (*regexp.Regexp, string, bool) {

	if nothingToReplace := (!_containsToFind(fileContent)); nothingToReplace {

		if shouldAddFill := (DoAddFill && !_hasFill(fileContent)); shouldAddFill {
			return ReAddNew, ToAdd, true
		}

		return nil, "", false
	}

	return ReToFind, ToReplace, true
}

func _containsToFind(fileContent string) bool {
	return strings.Contains(fileContent, ToFind)
}

func _hasFill(fileContent string) bool {
	return strings.Contains(fileContent, "fill=") ||
		strings.Contains(fileContent, "fill:")
}
