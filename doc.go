// Package godoc2md contains the code used to perform the CLI command
// `godoc2md`.
//
// This package is forked from https://github.com/davecheney/godoc2md
// which is no longer updated.
//
// godoc2md converts godoc formatted package documentation into Markdown format.
//
// ### Usage
//
// 	# Generate Package Readme
// 	$ godoc2md $PACKAGE > $GOPATH/src/$PACKAGE/README.md
//
// 	# See all Options
// 	$ godoc2md
// 	usage: godoc2md package [name ...]
// 	   -basePrefix string
// 	       path prefix of go files.
// 	   -ex
// 	       show examples in command line mode
// 	   -goroot string
// 	       Go root directory (default $GOROOT)
// 	   -hashformat string
// 	       source link URL hash format (default "#L%d")
// 	   -links
// 	       link identifiers to their declarations (default true)
// 	   -play
// 	       enable playground in web interface (default true)
// 	   -srclink string
// 	       if set, format for filename of source link
// 	   -tabwidth int
// 	       tab width (default 4)
// 	   -template string
// 	       path to an alternate template file
// 	   -timestamps
// 	       show timestamps with directory listings (default true)
// 	   -urlPrefix string
// 	       path prefix of go files (default "github.com")
// 	   -v	verbose mode
//
package godoc2md
