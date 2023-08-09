package godoc2md

import (
	"bytes"
	"fmt"
	"go/ast"
	"net/url"
	"path"
	"strings"
	"time"

	"golang.org/x/tools/godoc"
)

const (
	URLScheme = "https"
)

var (
	TimeFormat = "2-Jan-2006 15:04:05 -0700"
)

// TemplateUtils contains a collection of functions that can be used by the
// provided text template.
//
// TemplateUtils most likely cannot be created directly, and a new instance
// should be created by calling `NewTemplateUtils(config)`.
type TemplateUtils struct {
	sourceID          string
	basePrefix        string
	urlPrefix         string
	timeFormat        string
	srcLinkHashFormat string
}

// NewTemplateUtils returns a new TemplateUtils object configured from the
// provided CLI instance.
func NewTemplateUtils(cfg *Cli) TemplateUtils {
	return TemplateUtils{
		sourceID:          *cfg.SourceID,
		basePrefix:        *cfg.BasePrefix,
		urlPrefix:         *cfg.UrlPrefix,
		srcLinkHashFormat: *cfg.SrcLinkHashFormat,
		timeFormat:        TimeFormat,
	}
}

// Methods returns a map of name to func of all the methods of this struct. It's
// provided to the presenter and the keys are made available as functions to the
// template.
func (t TemplateUtils) Methods() map[string]interface{} {
	return map[string]interface{}{
		"comment_md":   t.CommentToMD,
		"srcfile_url":  t.GetSourceFileURL,
		"base":         t.StripBasePrefix,
		"md":           t.MDEscapeInline,
		"goCode":       t.MDEscapeGo,
		"kebab":        t.kebabFunc,
		"bitscape":     t.bitscapeFunc, //Escape [] for bitbucket confusion
		"trim_prefix":  strings.TrimPrefix,
		"last_item":    t.isLastItem,
		"current_time": t.GetCurrentTime,
		"get_full_url": t.GetFullURL,
	}
}

// CommentToMD converts the provided text, from Go source comment, into markdown.
func (t TemplateUtils) CommentToMD(comment string) string {
	var buf bytes.Buffer
	ToMD(&buf, comment)
	return buf.String()
}

// GetFullURL returns the URL, including line number, of the provided source
// code declaration.
func (t TemplateUtils) GetFullURL(pkg *godoc.PageInfo, decl ast.Decl) string {
	sourceURL, err := url.Parse(t.urlPrefix)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	if sourceURL.Host == "" {
		sourceURL, err = url.Parse(pkg.PDoc.ImportPath)
		if err != nil {
			return fmt.Sprintf("%v", err)
		}
	}
	sourceURL.Scheme = URLScheme

	// Gather the fragments of the intended file path/location.
	pathFragments := []string{t.getPathPrefix()}
	pathFragments = append(pathFragments, t.StripBasePrefix(pkg.PDoc.ImportPath))

	// find source file/position and generate string.
	sourceLoc := pkg.FSet.Position(decl.Pos())
	raw, err := url.Parse(fmt.Sprintf(t.srcLinkHashFormat, sourceLoc.Line))
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	pathFragments = append(pathFragments, raw.Path)
	sourceURL.Fragment = raw.Fragment
	sourceURL.RawQuery = raw.RawQuery

	// strings.Split(sourceLoc.Filename, "/")
	filename := path.Clean("/" + strings.TrimPrefix(sourceLoc.Filename, "/target"))
	pathFragments = append(pathFragments, filename)

	sourceURL.Path = path.Join(pathFragments...)

	return sourceURL.String()
}

// GetSourceFileURL reads the provided string and converts it into a URL.
func (t TemplateUtils) GetSourceFileURL(s string) string {
	sourceURL, err := url.Parse(t.urlPrefix)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	if sourceURL.Host == "" {
		sourceURL, err = url.Parse(s)
		if err != nil {
			return fmt.Sprintf("%v", err)
		}
	}
	sourceURL.Scheme = URLScheme

	repoPath := t.StripBasePrefix(sourceURL.Path)
	filename := path.Clean("/" + strings.TrimPrefix(repoPath, "/target"))
	sourceURL.Path = path.Join(t.getPathPrefix(), filename)

	return sourceURL.String()
}

// StripBasePrefix removes the configured basePrefix from the provided string.
func (t TemplateUtils) StripBasePrefix(path string) string {
	return strings.TrimPrefix(path, t.basePrefix)
}

// MDEscapeInline escapes inline emphasis and bold marks.
func (t TemplateUtils) MDEscapeInline(text string) string {
	text = strings.Replace(text, "*", "\\*", -1)
	text = strings.Replace(text, "_", "\\_", -1)

	return text
}

// MDEscapeGo fences a string of text as Go Code.
func (t TemplateUtils) MDEscapeGo(text string) string {
	return "```go\n" + strings.TrimRight(text, " \n") + "\n```\n"
}

// GetCurrentTime returns the current time in UTC using the configured format.
func (t TemplateUtils) GetCurrentTime() string {
	return time.Now().UTC().Format(t.timeFormat)
}

func (t TemplateUtils) kebabFunc(text string) string {
	s := strings.Replace(strings.ToLower(text), " ", "-", -1)
	s = strings.Replace(s, ".", "-", -1)
	s = strings.Replace(s, "\\*", "42", -1)
	return s
}

func (t TemplateUtils) bitscapeFunc(text string) string {
	s := strings.Replace(text, "[", "\\[", -1)
	s = strings.Replace(s, "]", "\\]", -1)
	return s
}

func (t TemplateUtils) getPathPrefix() string {
	branchPath := fmt.Sprintf("blob/%s", t.sourceID)

	return path.Join(t.basePrefix, branchPath)
}

func (t TemplateUtils) isLastItem(idx int, list []string) bool {
	return idx+1 >= len(list)
}
