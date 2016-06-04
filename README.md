# svg-edit
## Install
    go get github.com/GreenRaccoon23/svg-edit
## Clone
    git clone https://github.com/GreenRaccoon23/svg-edit.git
## Description
CLI tool to edit strings in svg images. Created for my [archdroid-icon-theme](https://github.com/GreenRaccoon23/archdroid-icon-theme) repo.

    [chuck@norris ~]$ svg-edit --help
    Usage: svg-edit [options] <original file/directory> <new file/directory>
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
                 (QUIET) don't show any output
 
## Features
#### Material Design Color Shortcuts
As an extra feature for my [archdroid-icon-theme](https://github.com/GreenRaccoon23/archdroid-icon-theme) repo, the `-o` and `-n` options have shortcuts to the hex colors in the [Material Design Color Palette](http://www.google.com/design/spec/style/color.html#color-color-palette "Material Design Color Palette"). The shortcuts are in the format `<color-group><shade>`.  

Examples:

    green => #4CAF50
    green500 => #4CAF50
    green-500 => #4CAF50
    green:500 => #4CAF50

    green900 => #1B5E20
    greendark => #1B5E20
    green-dark => #1B5E20
    green:dark => #1B5E20

    green100 => #C8E6C9
    greenlight => #C8E6C9
