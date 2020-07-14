// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// godoc2md
package main

import (
	"go/build"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/tools/godoc"
	"golang.org/x/tools/godoc/vfs"

	"github.com/thatgerber/godoc2md"
)

func main() {
	args, config := godoc2md.Parse()

	fs := newFilesystem(*config.Goroot)
	corpus := godoc.NewCorpus(fs)
	corpus.Verbose = *config.Verbose
	pres := godoc2md.NewPresentation(corpus, config)
	output := os.Stdout

	if err := godoc.CommandLine(output, fs, pres, args); err != nil {
		log.Print(err)
	}
}

func newFilesystem(root string) vfs.NameSpace {
	fs := vfs.NameSpace{}

	// use file system of underlying OS
	fs.Bind("/", vfs.OS(root), "/", vfs.BindReplace)

	// Bind $GOPATH trees into Go root.
	for _, p := range filepath.SplitList(build.Default.GOPATH) {
		fs.Bind("/src/pkg", vfs.OS(p), "/src", vfs.BindAfter)
	}

	return fs
}
