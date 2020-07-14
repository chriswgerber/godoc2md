package godoc2md

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
)

var (
	cmdName = "godoc2md"

	defaultURLPrefix = "github.com"

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

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s package [more-packages ...]\n", cmdName)
	flag.PrintDefaults()
	os.Exit(2)
}

func getBasePkgPrefix(potench string) *string {
	cwd, _ := os.Getwd()
	modfilePath := path.Join(cwd, "go.mod")

	// Check if we have a go.mod
	if _, err := os.Stat(modfilePath); err == nil || !os.IsNotExist(err) {
		file, err := os.Open(modfilePath)
		if err != nil {
			log.Fatalf("failed to open %s: %v", modfilePath, err)
		}

		fileBuf := bufio.NewReader(file)
		nlByte := []byte("\n")
		l1, err := fileBuf.ReadString(nlByte[0])
		if err == nil && l1 != "" {
			tmpString := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(l1, "module "), "\n"))
			return &tmpString
		}
	}

	// Try and guess the package path
	p := os.Getenv("GOPATH")
	if p != "" {
		if newPath := strings.TrimPrefix(cwd, path.Join(p, "/src/")); newPath != "" {
			return &newPath
		}
	}

	return &potench
}

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

func Parse() ([]string, *Cli) {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		usage()
	}

	if *Config.Goroot == "" {
		root := runtime.GOROOT()
		Config.Goroot = &root
	}

	if *Config.BasePrefix == "" {
		Config.BasePrefix = getBasePkgPrefix(args[0])
	}

	return args, Config
}
