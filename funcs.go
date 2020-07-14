package godoc2md

import (
	"bytes"
	"fmt"
	"go/doc"
	"net/url"
	"path"
	"strings"
	"time"
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

func (t TemplateUtils) getFullURL(pkg *doc.Package, s string) string {
	sourceURL, err := url.Parse(t.urlPrefix)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	sourceURL.Scheme = defaultScheme

	repoPath := strings.TrimPrefix(pkg.ImportPath, sourceURL.String())
	raw, err := url.Parse(s)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	sourceURL.Fragment = raw.Fragment
	sourceURL.RawQuery = raw.RawQuery
	sourceURL.Path = path.Join(sourceURL.Path, fileBranchPath, repoPath, raw.Path)

	return sourceURL.String()
}
