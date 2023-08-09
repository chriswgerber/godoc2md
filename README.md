# godoc2md

`import "github.com/chriswgerber/godoc2md"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>

Package godoc2md contains the code used to perform the CLI command
`godoc2md`.

[![GoDoc](<a href="https://godoc.org/github.com/chriswgerber/godoc2md?status.svg">https://godoc.org/github.com/chriswgerber/godoc2md?status.svg</a>)](<a href="https://godoc.org/github.com/chriswgerber/godoc2md">https://godoc.org/github.com/chriswgerber/godoc2md</a>)

This package is forked from <a href="https://github.com/davecheney/godoc2md">https://github.com/davecheney/godoc2md</a>
which is no longer updated.

godoc2md converts godoc formatted package documentation into Markdown format.

```
# Generate Package Readme
$ godoc2md $PACKAGE > $GOPATH/src/$PACKAGE/README.md

# See all Options
$ godoc2md
  usage: godoc2md package [more-packages ...]
  -basePrefix go.mod
	  	path prefix of go files. If not set, cli will attempt to set it by checking go.mod, current directory, and the 1st position argument
  -ex
	  	show examples in command line mode
  -goroot GOROOT
	  	directory of Go Root. Will attempt to lookup from GOROOT
  -hashformat string
	  	source link URL hash format (default "#L%d")
  -links
	  	link identifiers to their declarations (default true)
  -play
	  	enable playground in web interface (default true)
  -srclink string
	  	if set, format for filename of source link
  -tabwidth int
	  	tab width (default 4)
  -template string
	  	path to an alternate template file
  -timestamps
	  	show timestamps with directory listings (default true)
  -urlPrefix string
	  	URL for generated URLs. (default "github.com")
  -v	verbose mode
```

## <a name="pkg-index">Index</a>

* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func NewPresentation(corpus *godoc.Corpus, config *Cli) *godoc.Presentation](#NewPresentation)
* [func ToMD(w io.Writer, text string)](#ToMD)
* [type Cli](#Cli)
  * [func Parse() ([]string, *Cli)](#Parse)
* [type TemplateUtils](#TemplateUtils)
  * [func NewTemplateUtils(cfg *Cli) TemplateUtils](#NewTemplateUtils)
  * [func (t TemplateUtils) CommentToMD(comment string) string](#TemplateUtils.CommentToMD)
  * [func (t TemplateUtils) GetCurrentTime() string](#TemplateUtils.GetCurrentTime)
  * [func (t TemplateUtils) GetFullURL(pkg *godoc.PageInfo, decl ast.Decl) string](#TemplateUtils.GetFullURL)
  * [func (t TemplateUtils) GetSourceFileURL(s string) string](#TemplateUtils.GetSourceFileURL)
  * [func (t TemplateUtils) MDEscapeGo(text string) string](#TemplateUtils.MDEscapeGo)
  * [func (t TemplateUtils) MDEscapeInline(text string) string](#TemplateUtils.MDEscapeInline)
  * [func (t TemplateUtils) Methods() map[string]interface{}](#TemplateUtils.Methods)
  * [func (t TemplateUtils) StripBasePrefix(path string) string](#TemplateUtils.StripBasePrefix)

#### <a name="pkg-files">Package files</a>

[comment.go](https://github.com/chriswgerber/godoc2md/blob/master/comment.go) [config.go](https://github.com/chriswgerber/godoc2md/blob/master/config.go) [doc.go](https://github.com/chriswgerber/godoc2md/blob/master/doc.go) [funcs.go](https://github.com/chriswgerber/godoc2md/blob/master/funcs.go) [presentation.go](https://github.com/chriswgerber/godoc2md/blob/master/presentation.go) [template.go](https://github.com/chriswgerber/godoc2md/blob/master/template.go) 

## <a name="pkg-constants">Constants</a>

```go
const (
    URLScheme = "https"
)
```

## <a name="pkg-variables">Variables</a>

```go
var (

    // Config contains the configuration for the CLI. To populate config, call
    // `Parse()` and use the provided response.
    Config = &Cli{
        Verbose:           flag.Bool("v", false, "verbose mode"),
        Goroot:            flag.String("goroot", "", "directory of Go Root. Will attempt to lookup from `GOROOT`"),
        TabWidth:          flag.Int("tabwidth", 4, "tab width"),
        ShowTimestamps:    flag.Bool("timestamps", true, "show timestamps with directory listings"),
        BasePrefix:        flag.String("basePrefix", "", "path prefix of go files. If not set, cli will attempt to set it by checking `go.mod`, current directory, and the 1st position argument"),
        UrlPrefix:         flag.String("urlPrefix", defaultURLPrefix, "URL for generated URLs."),
        SourceID:          flag.String("sourceID", defaultSourceID, "URL for generated URLs."),
        AltPkgTemplate:    flag.String("template", "", "path to an alternate template file"),
        ShowPlayground:    flag.Bool("play", true, "enable playground in web interface"),
        ShowExamples:      flag.Bool("ex", false, "show examples in command line mode"),
        DeclLinks:         flag.Bool("links", true, "link identifiers to their declarations"),
        SrcLinkHashFormat: flag.String("hashformat", "#L%d", "source link URL hash format"),
    }
)
```

```go
var (
    TimeFormat = "2-Jan-2006 15:04:05 -0700"
)
```

## <a name="NewPresentation">func</a> [NewPresentation](https://github.com/chriswgerber/godoc2md/blob/master/presentation.go#L46)

```go
func NewPresentation(corpus *godoc.Corpus, config *Cli) *godoc.Presentation
```

## <a name="ToMD">func</a> [ToMD](https://github.com/chriswgerber/godoc2md/blob/master/comment.go#L58)

```go
func ToMD(w io.Writer, text string)
```

ToMD converts comment text to formatted Markdown. The comment was prepared by
DocReader, so it is known not to have leading, trailing blank lines nor to
have trailing spaces at the end of lines. The comment markers have already
been removed.

Each span of unindented non-blank lines is converted into a single paragraph.
There is one exception to the rule: a span that consists of a single line, is
followed by another paragraph span, begins with a capital letter, and
contains no punctuation is formatted as a heading.

A span of indented lines is converted into a `<pre>` block, with the common
indent prefix removed.

URLs in the comment text are converted into links.

## <a name="Cli">type</a> [Cli](https://github.com/chriswgerber/godoc2md/blob/master/config.go#L75)

```go
type Cli struct {
    Verbose *bool
    // Goroot
    Goroot *string

    // layout control
    TabWidth       *int
    ShowTimestamps *bool
    BasePrefix     *string
    UrlPrefix      *string
    SourceID       *string
    AltPkgTemplate *string
    ShowPlayground *bool
    ShowExamples   *bool
    DeclLinks      *bool

    // The hash format for Github is the default `#L%d`; but other source control platforms do not
    // use the same format. For example Bitbucket Enterprise uses `#%d`. This option provides the
    // user the option to switch the format as needed and still remain backwards compatible.
    SrcLinkHashFormat *string
}
```

### <a name="Parse">func</a> [Parse](https://github.com/chriswgerber/godoc2md/blob/master/config.go#L97)

```go
func Parse() ([]string, *Cli)
```

## <a name="TemplateUtils">type</a> [TemplateUtils](https://github.com/chriswgerber/godoc2md/blob/master/funcs.go#L28)

```go
type TemplateUtils struct {
    // contains filtered or unexported fields
}
```

TemplateUtils contains a collection of functions that can be used by the
provided text template.

TemplateUtils most likely cannot be created directly, and a new instance
should be created by calling `NewTemplateUtils(config)`.

### <a name="NewTemplateUtils">func</a> [NewTemplateUtils](https://github.com/chriswgerber/godoc2md/blob/master/funcs.go#L38)

```go
func NewTemplateUtils(cfg *Cli) TemplateUtils
```

NewTemplateUtils returns a new TemplateUtils object configured from the
provided CLI instance.

### <a name="TemplateUtils.CommentToMD">func</a> (TemplateUtils) [CommentToMD](https://github.com/chriswgerber/godoc2md/blob/master/funcs.go#L68)

```go
func (t TemplateUtils) CommentToMD(comment string) string
```

CommentToMD converts the provided text, from Go source comment, into markdown.

### <a name="TemplateUtils.GetCurrentTime">func</a> (TemplateUtils) [GetCurrentTime](https://github.com/chriswgerber/godoc2md/blob/master/funcs.go#L152)

```go
func (t TemplateUtils) GetCurrentTime() string
```

GetCurrentTime returns the current time in UTC using the configured format.

### <a name="TemplateUtils.GetFullURL">func</a> (TemplateUtils) [GetFullURL](https://github.com/chriswgerber/godoc2md/blob/master/funcs.go#L76)

```go
func (t TemplateUtils) GetFullURL(pkg *godoc.PageInfo, decl ast.Decl) string
```

GetFullURL returns the URL, including line number, of the provided source
code declaration.

### <a name="TemplateUtils.GetSourceFileURL">func</a> (TemplateUtils) [GetSourceFileURL](https://github.com/chriswgerber/godoc2md/blob/master/funcs.go#L113)

```go
func (t TemplateUtils) GetSourceFileURL(s string) string
```

GetSourceFileURL reads the provided string and converts it into a URL.

### <a name="TemplateUtils.MDEscapeGo">func</a> (TemplateUtils) [MDEscapeGo](https://github.com/chriswgerber/godoc2md/blob/master/funcs.go#L147)

```go
func (t TemplateUtils) MDEscapeGo(text string) string
```

MDEscapeGo fences a string of text as Go Code.

### <a name="TemplateUtils.MDEscapeInline">func</a> (TemplateUtils) [MDEscapeInline](https://github.com/chriswgerber/godoc2md/blob/master/funcs.go#L139)

```go
func (t TemplateUtils) MDEscapeInline(text string) string
```

MDEscapeInline escapes inline emphasis and bold marks.

### <a name="TemplateUtils.Methods">func</a> (TemplateUtils) [Methods](https://github.com/chriswgerber/godoc2md/blob/master/funcs.go#L51)

```go
func (t TemplateUtils) Methods() map[string]interface{}
```

Methods returns a map of name to func of all the methods of this struct. It's
provided to the presenter and the keys are made available as functions to the
template.

### <a name="TemplateUtils.StripBasePrefix">func</a> (TemplateUtils) [StripBasePrefix](https://github.com/chriswgerber/godoc2md/blob/master/funcs.go#L134)

```go
func (t TemplateUtils) StripBasePrefix(path string) string
```

StripBasePrefix removes the configured basePrefix from the provided string.

- - -
Created: 9-Aug-2023 20:02:57 +0000
Generated by [godoc2md](http://github.com/chriswgerber/godoc2md)
