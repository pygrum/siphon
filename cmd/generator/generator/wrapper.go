package generator

import (
	"bytes"
	"fmt"
	"github.com/pygrum/siphon/internal/logger"
	"os"
	"os/exec"
	"strings"
)

type Builder struct {
	CC       string
	GOOS     string
	GOARCH   string
	SrcPaths []string
	outFile  string
	flags    []Flag
}

type Flag struct {
	Name  string
	Value string
}

const (
	outFileOption = "-o"
	flagsOption   = "-ldflags"
	flagPrefix    = "-X"
)

func NewBuilder(cc, goos, goarch string) *Builder {
	return &Builder{
		CC:     cc,
		GOOS:   goos,
		GOARCH: goarch,
	}
}

func (b *Builder) AddSrcPath(path string) {
	b.SrcPaths = append(b.SrcPaths, path)
}

func (b *Builder) SetFlags(flags ...Flag) {
	b.flags = flags
}

func (b *Builder) SetOutFile(name string) {
	b.outFile = name
}

func (b *Builder) Build() {
	var buildCmd []string
	buildCmd = append(buildCmd, "build")
	buildCmd = append(buildCmd, outFileOption)
	buildCmd = append(buildCmd, b.outFile)
	buildCmd = append(buildCmd, flagsOption)
	var flags []string
	for _, f := range b.flags {
		flags = append(flags, flagPrefix)
		formatString := "'%s=%s'"
		flags = append(flags, fmt.Sprintf(formatString, f.Name, f.Value))
	}
	buildCmd = append(buildCmd, strings.Join(flags, " "))
	for _, s := range b.SrcPaths {
		buildCmd = append(buildCmd, s)
	}
	fmt.Println(b.CC, strings.Join(buildCmd, " "))
	var cerr bytes.Buffer
	cmd := exec.Command(b.CC, buildCmd...)
	// Set arch and os environment vars
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", b.GOOS), fmt.Sprintf("GOARCH=%s", b.GOARCH)) // go-sqlite3 requires cgo
	cmd.Stderr = &cerr
	if err := cmd.Run(); err != nil {
		logger.Fatalf("%v: %s", err, cerr.String())
	}
}
