// Package godoc2md contains the code used to perform the CLI command
// `godoc2md`.
//
// [![GoDoc](https://godoc.org/github.com/ThatGerber/godoc2md?status.svg)](https://godoc.org/github.com/ThatGerber/godoc2md)
//
// This package is forked from https://github.com/davecheney/godoc2md
// which is no longer updated.
//
// godoc2md converts godoc formatted package documentation into Markdown format.
//
// 	# Generate Package Readme
// 	$ godoc2md $PACKAGE > $GOPATH/src/$PACKAGE/README.md
//
// 	# See all Options
// 	$ godoc2md
// 	  usage: godoc2md package [more-packages ...]
// 	  -basePrefix go.mod
// 		  	path prefix of go files. If not set, cli will attempt to set it by checking go.mod, current directory, and the 1st position argument
// 	  -ex
// 		  	show examples in command line mode
// 	  -goroot GOROOT
// 		  	directory of Go Root. Will attempt to lookup from GOROOT
// 	  -hashformat string
// 		  	source link URL hash format (default "#L%d")
// 	  -links
// 		  	link identifiers to their declarations (default true)
// 	  -play
// 		  	enable playground in web interface (default true)
// 	  -srclink string
// 		  	if set, format for filename of source link
// 	  -tabwidth int
// 		  	tab width (default 4)
// 	  -template string
// 		  	path to an alternate template file
// 	  -timestamps
// 		  	show timestamps with directory listings (default true)
// 	  -urlPrefix string
// 		  	URL for generated URLs. (default "github.com")
// 	  -v	verbose mode
//
package godoc2md
