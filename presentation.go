package godoc2md

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
	"text/template"

	"golang.org/x/tools/godoc"
)

var (
	templateName = "template.go"
)

type sourceLinker struct {
	HashFormat string
}

// Original Source https://github.com/golang/tools/blob/master/godoc/godoc.go#L562
func (l *sourceLinker) source(s string) string {
	return path.Clean("/" + strings.TrimPrefix(s, "/target"))
}

// Removed code line that always substracted 10 from the value of `line`.
// Made format for the source link hash configurable to support source control platforms other than Github.
// Original Source https://github.com/golang/tools/blob/master/godoc/godoc.go#L540
func (l *sourceLinker) sourcePosition(p string, line, low, high int) string {
	f := l.source(p)

	var buf bytes.Buffer
	template.HTMLEscape(&buf, []byte(f))

	// line id's in html-printed source are of the
	// form "L%d" (on Github) where %d stands for the line number
	if line > 0 {
		fmt.Fprintf(&buf, l.HashFormat, line) // no need for URL escaping
	}

	return buf.String()
}

func NewPresentation(corpus *godoc.Corpus, config *Cli) *godoc.Presentation {
	pres := godoc.NewPresentation(corpus)

	pres.TabWidth = *config.TabWidth
	pres.ShowTimestamps = *config.ShowTimestamps
	pres.ShowPlayground = *config.ShowPlayground
	pres.ShowExamples = *config.ShowExamples
	pres.DeclLinks = *config.DeclLinks
	pres.SrcMode = false
	pres.HTMLMode = false

	sl := &sourceLinker{HashFormat: *config.SrcLinkHashFormat}
	pres.URLForSrc = sl.source
	pres.URLForSrcPos = sl.sourcePosition

	templateText := pkgTemplate
	if *config.AltPkgTemplate != "" {
		templateName = *config.AltPkgTemplate
		buf, err := ioutil.ReadFile(templateName)
		if err != nil {
			log.Fatal(err)
		}
		templateText = string(buf)
	}
	docTemplate := template.New(templateName)
	docTemplate.Funcs(pres.FuncMap())

	utilFuncs := NewTemplateUtils(config)
	docTemplate.Funcs(utilFuncs.Methods())

	var err error
	pres.PackageText, err = docTemplate.Parse(templateText)

	if err != nil {
		log.Fatal("error parsing template: %v", err)
	}

	return pres
}
