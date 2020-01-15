package main

import (
	"bytes"
	"fmt"
	"go/doc"
	"net/url"
	"path"
	"strings"
	"time"
)

var timeFormat = "2-Jan-2006 15:04:05 -0700"

var funcs = map[string]interface{}{
	"comment_md":   commentMdFunc,
	"base":         baseFunc,
	"md":           mdFunc,
	"pre":          preFunc,
	"kebab":        kebabFunc,
	"bitscape":     bitscapeFunc, //Escape [] for bitbucket confusion
	"trim_prefix":  strings.TrimPrefix,
	"last_item":    isLastItem,
	"current_time": getCurrentTime,
	"get_full_url": getFullURL,
}

func commentMdFunc(comment string) string {
	var buf bytes.Buffer
	ToMD(&buf, comment)
	return buf.String()
}

func baseFunc(path string) string {
	return strings.TrimPrefix(path, *basePrefix)
}

func mdFunc(text string) string {
	text = strings.Replace(text, "*", "\\*", -1)
	text = strings.Replace(text, "_", "\\_", -1)
	return text
}

func preFunc(text string) string {
	return "```go\n" + strings.TrimRight(text, " \n") + "\n```\n"
}

func kebabFunc(text string) string {
	s := strings.Replace(strings.ToLower(text), " ", "-", -1)
	s = strings.Replace(s, ".", "-", -1)
	s = strings.Replace(s, "\\*", "42", -1)
	return s
}

func bitscapeFunc(text string) string {
	s := strings.Replace(text, "[", "\\[", -1)
	s = strings.Replace(s, "]", "\\]", -1)
	return s
}

func getCurrentTime() string {
	return time.Now().Format(timeFormat)
}

func isLastItem(idx int, list []string) bool {
	return idx+1 >= len(list)
}

const (
	defaultScheme  = "https"
	fileBranchPath = "blob/master"
)

func getFullURL(pkg *doc.Package, s string) string {
	sourceURL, err := url.Parse(*urlPrefix)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	repoPath := strings.TrimPrefix(pkg.ImportPath, sourceURL.String())
	raw, err := url.Parse(s)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	sourceURL.Fragment = raw.Fragment
	sourceURL.RawQuery = raw.RawQuery
	sourceURL.Path = path.Join(sourceURL.Path, fileBranchPath, repoPath, raw.Path)
	sourceURL.Scheme = defaultScheme
	return sourceURL.String()
}
