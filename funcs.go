package godoc2md

import (
	"bytes"
	"fmt"
	"go/ast"
	"net/url"
	"path"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/godoc"
)

const (
	defaultScheme  = "https"
	fileBranchPath = "blob/master"
)

var (
	TimeFormat = "2-Jan-2006 15:04:05 -0700"
)

type TemplateUtils struct {
	basePrefix string
	urlPrefix  string
	timeFormat string
}

func NewTemplateUtils(cfg *Cli) TemplateUtils {
	return TemplateUtils{
		basePrefix: *cfg.BasePrefix,
		urlPrefix:  *cfg.UrlPrefix,
		timeFormat: TimeFormat,
	}
}

func (t TemplateUtils) Methods() map[string]interface{} {
	return map[string]interface{}{
		"comment_md":   t.commentMdFunc,
		"srcfile_url":  t.getSourceFileURL,
		"base":         t.baseFunc,
		"md":           t.mdFunc,
		"pre":          t.preFunc,
		"kebab":        t.kebabFunc,
		"bitscape":     t.bitscapeFunc, //Escape [] for bitbucket confusion
		"trim_prefix":  strings.TrimPrefix,
		"last_item":    t.isLastItem,
		"current_time": t.getCurrentTime,
		"get_full_url": t.getFullURL,
	}
}

func startsWithUppercase(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsUpper(r)
}

func splitExampleName(s string) (name, suffix string) {
	i := strings.LastIndex(s, "_")
	if 0 <= i && i < len(s)-1 && !startsWithUppercase(s[i+1:]) {
		name = s[:i]
		suffix = " (" + strings.Title(s[i+1:]) + ")"
		return
	}
	name = s
	return
}

func (t TemplateUtils) commentMdFunc(comment string) string {
	var buf bytes.Buffer
	ToMD(&buf, comment)
	return buf.String()
}

func (t TemplateUtils) baseFunc(path string) string {
	return strings.TrimPrefix(path, t.basePrefix)
}

func (t TemplateUtils) mdFunc(text string) string {
	text = strings.Replace(text, "*", "\\*", -1)
	text = strings.Replace(text, "_", "\\_", -1)
	return text
}

func (t TemplateUtils) preFunc(text string) string {
	return "```go\n" + strings.TrimRight(text, " \n") + "\n```\n"
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

func (t TemplateUtils) getCurrentTime() string {
	return time.Now().Format(t.timeFormat)
}

func (t TemplateUtils) isLastItem(idx int, list []string) bool {
	return idx+1 >= len(list)
}

func (t TemplateUtils) getFullURL(pkg *godoc.PageInfo, decl ast.Decl) string {
	sourceURL, err := url.Parse(t.urlPrefix)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	sourceURL.Scheme = defaultScheme

	repoPath := t.baseFunc(pkg.PDoc.ImportPath)

	sourceLoc := pkg.FSet.Position(decl.Pos())
	raw, err := url.Parse(fmt.Sprintf(*Config.SrcLinkHashFormat, sourceLoc.Line))
	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	filename := strings.Split(sourceLoc.Filename, "/")

	sourceURL.Fragment = raw.Fragment
	sourceURL.RawQuery = raw.RawQuery
	sourceURL.Path = path.Join(t.basePrefix, fileBranchPath, repoPath, raw.Path, filename[len(filename)-1])

	return sourceURL.String()
}

func (t TemplateUtils) getSourceFileURL(s string) string {
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
	repoPath := t.baseFunc(sourceURL.Path)
	sourceURL.Scheme = defaultScheme
	sourceURL.Path = path.Join(t.basePrefix, fileBranchPath, repoPath)

	return sourceURL.String()
}
