# godoc2md

`import "github.com/thatgerber/godoc2md"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>

Package godoc2md contains the code used to perform the CLI command
`godoc2md`.

[![GoDoc](<a href="https://godoc.org/github.com/ThatGerber/godoc2md?status.svg">https://godoc.org/github.com/ThatGerber/godoc2md?status.svg</a>)](<a href="https://godoc.org/github.com/ThatGerber/godoc2md">https://godoc.org/github.com/ThatGerber/godoc2md</a>)

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

* [Variables](#pkg-variables)
* [func NewPresentation(corpus *godoc.Corpus, config *Cli) *godoc.Presentation](#NewPresentation)
* [func ToMD(w io.Writer, text string)](#ToMD)
* [type Cli](#Cli)
  * [func Parse() ([]string, *Cli)](#Parse)
* [type TemplateUtils](#TemplateUtils)
  * [func NewTemplateUtils(cfg *Cli) TemplateUtils](#NewTemplateUtils)
  * [func (t TemplateUtils) Methods() map[string]interface{}](#TemplateUtils.Methods)

#### <a name="pkg-files">Package files</a>

[comment.go](https://github.com/thatgerber/godoc2md/blob/master/comment.go) [config.go](https://github.com/thatgerber/godoc2md/blob/master/config.go) [doc.go](https://github.com/thatgerber/godoc2md/blob/master/doc.go) [funcs.go](https://github.com/thatgerber/godoc2md/blob/master/funcs.go) [presentation.go](https://github.com/thatgerber/godoc2md/blob/master/presentation.go) [template.go](https://github.com/thatgerber/godoc2md/blob/master/template.go) 

## <a name="pkg-variables">Variables</a>

```go
var (

    // Config contains the configuration for the CLI. To populate config, call
    // `Parse()`.
    Config = &Cli{
        Verbose:           flag.Bool("v", false, "verbose mode"),
        Goroot:            flag.String("goroot", "", "directory of Go Root. Will attempt to lookup from `GOROOT`"),
        TabWidth:          flag.Int("tabwidth", 4, "tab width"),
        ShowTimestamps:    flag.Bool("timestamps", true, "show timestamps with directory listings"),
        BasePrefix:        flag.String("basePrefix", "", "path prefix of go files. If not set, cli will attempt to set it by checking `go.mod`, current directory, and the 1st position argument"),
        UrlPrefix:         flag.String("urlPrefix", defaultURLPrefix, "URL for generated URLs."),
        AltPkgTemplate:    flag.String("template", "", "path to an alternate template file"),
        ShowPlayground:    flag.Bool("play", true, "enable playground in web interface"),
        ShowExamples:      flag.Bool("ex", false, "show examples in command line mode"),
        DeclLinks:         flag.Bool("links", true, "link identifiers to their declarations"),
        SrcLinkHashFormat: flag.String("hashformat", "#L%d", "source link URL hash format"),
        SrcLinkFormat:     flag.String("srclink", "", "if set, format for filename of source link"),
    }
)
```

```go
var (
    TimeFormat = "2-Jan-2006 15:04:05 -0700"
)
```

## <a name="NewPresentation">func</a> [NewPresentation](https://github.com/thatgerber/godoc2md/blob/master/presentation.go#L46)

```go
func NewPresentation(corpus *godoc.Corpus, config *Cli) *godoc.Presentation
```

## <a name="ToMD">func</a> [ToMD](https://github.com/thatgerber/godoc2md/blob/master/comment.go#L58)

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

## <a name="Cli">type</a> [Cli](https://github.com/thatgerber/godoc2md/blob/master/config.go#L74)

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
    AltPkgTemplate *string
    ShowPlayground *bool
    ShowExamples   *bool
    DeclLinks      *bool

    // The hash format for Github is the default `#L%d`; but other source control platforms do not
    // use the same format. For example Bitbucket Enterprise uses `#%d`. This option provides the
    // user the option to switch the format as needed and still remain backwards compatible.
    SrcLinkHashFormat *string
    SrcLinkFormat     *string
}
```

### <a name="Parse">func</a> [Parse](https://github.com/thatgerber/godoc2md/blob/master/config.go#L96)

```go
func Parse() ([]string, *Cli)
```

## <a name="TemplateUtils">type</a> [TemplateUtils](https://github.com/thatgerber/godoc2md/blob/master/funcs.go#L26)

```go
type TemplateUtils struct {
    // contains filtered or unexported fields
}
```

### <a name="NewTemplateUtils">func</a> [NewTemplateUtils](https://github.com/thatgerber/godoc2md/blob/master/funcs.go#L32)

```go
func NewTemplateUtils(cfg *Cli) TemplateUtils
```

### <a name="TemplateUtils.Methods">func</a> (TemplateUtils) [Methods](https://github.com/thatgerber/godoc2md/blob/master/funcs.go#L40)

```go
func (t TemplateUtils) Methods() map[string]interface{}
```

- - -
Generated by [godoc2md](http://github.com/thatgerber/godoc2md)
