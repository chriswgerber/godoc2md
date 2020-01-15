// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// godoc2md converts godoc formatted package documentation into Markdown format.
//
//
// Usage
//
//    godoc2md $PACKAGE > $GOPATH/src/$PACKAGE/README.md
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"golang.org/x/tools/godoc"
	"golang.org/x/tools/godoc/vfs"
)

var (
	verbose = flag.Bool("v", false, "verbose mode")

	// file system roots
	// TODO(gri) consider the invariant that goroot always end in '/'
	goroot = flag.String("goroot", runtime.GOROOT(), "Go root directory")

	// layout control
	tabWidth       = flag.Int("tabwidth", 4, "tab width")
	showTimestamps = flag.Bool("timestamps", true, "show timestamps with directory listings")
	basePrefix     = flag.String("basePrefix", "gitlab.com/welllabs/devops/", "path prefix of go files")
	urlPrefix      = flag.String("urlPrefix", "gitlab.com", "path prefix of go files")
	altPkgTemplate = flag.String("template", "", "path to an alternate template file")
	showPlayground = flag.Bool("play", true, "enable playground in web interface")
	showExamples   = flag.Bool("ex", false, "show examples in command line mode")
	declLinks      = flag.Bool("links", true, "link identifiers to their declarations")

	// The hash format for Github is the default `#L%d`; but other source control platforms do not
	// use the same format. For example Bitbucket Enterprise uses `#%d`. This option provides the
	// user the option to switch the format as needed and still remain backwards compatible.
	srcLinkHashFormat = flag.String("hashformat", "#L%d", "source link URL hash format")
	srcLinkFormat     = flag.String("srclink", "", "if set, format for filename of source link")
)

func usage() {
	fmt.Fprintf(os.Stderr,
		"usage: godoc2md package [name ...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	pres *godoc.Presentation
	fs   = vfs.NameSpace{}
)

// Original Source https://github.com/golang/tools/blob/master/godoc/godoc.go#L562
func srcLinkFunc(s string) string {
	return path.Clean("/" + strings.TrimPrefix(s, "/target"))
}

// Removed code line that always substracted 10 from the value of `line`.
// Made format for the source link hash configurable to support source control platforms other than Github.
// Original Source https://github.com/golang/tools/blob/master/godoc/godoc.go#L540
func srcPosLinkFunc(s string, line, low, high int) string {
	s = srcLinkFunc(s)

	var buf bytes.Buffer
	template.HTMLEscape(&buf, []byte(s))

	// line id's in html-printed source are of the
	// form "L%d" (on Github) where %d stands for the line number
	if line > 0 {
		fmt.Fprintf(&buf, *srcLinkHashFormat, line) // no need for URL escaping
	}
	return buf.String()
}

func readTemplate(name, data string) *template.Template {
	// be explicit with errors (for app engine use)
	t, err := template.New(name).Funcs(pres.FuncMap()).Funcs(funcs).Parse(data)
	if err != nil {
		log.Fatal("readTemplate: ", err)
	}
	return t
}

func main() {
	flag.Usage = usage
	flag.Parse()

	// Check usage
	if flag.NArg() == 0 {
		usage()
	}

	// use file system of underlying OS
	fs.Bind("/", vfs.OS(*goroot), "/", vfs.BindReplace)

	// Bind $GOPATH trees into Go root.
	for _, p := range filepath.SplitList(build.Default.GOPATH) {
		fs.Bind("/src/pkg", vfs.OS(p), "/src", vfs.BindAfter)
	}

	corpus := godoc.NewCorpus(fs)
	corpus.Verbose = *verbose

	pres = godoc.NewPresentation(corpus)
	pres.TabWidth = *tabWidth
	pres.ShowTimestamps = *showTimestamps
	pres.ShowPlayground = *showPlayground
	pres.ShowExamples = *showExamples
	pres.DeclLinks = *declLinks
	pres.URLForSrc = srcLinkFunc
	pres.URLForSrcPos = srcPosLinkFunc
	pres.SrcMode = false
	pres.HTMLMode = false

	if *altPkgTemplate != "" {
		buf, err := ioutil.ReadFile(*altPkgTemplate)
		if err != nil {
			log.Fatal(err)
		}
		pres.PackageText = readTemplate("package.txt", string(buf))
	} else {
		pres.PackageText = readTemplate("package.txt", pkgTemplate)
	}

	if err := godoc.CommandLine(os.Stdout, fs, pres, flag.Args()); err != nil {
		log.Print(err)
	}
}
