package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	PrintUsage  bool
	ExecInstall bool
	File        string
	Updir       bool
	Links       bool
	Abs         bool
}

var UsageFunc func()

func ReadArgs(args []string) (error, Config) {
	conf := Config{}
	fs := flag.NewFlagSet("YamlIncludeTransformer", flag.ContinueOnError)
	UsageFunc = func() {
		fmt.Fprintf(os.Stderr, usagestr, os.Args[0])
	}
	fs.Usage = UsageFunc
	fs.BoolVar(&conf.PrintUsage, "help", false, "Print help message")
	fs.BoolVar(&conf.PrintUsage, "h", false, "Print help message")
	fs.BoolVar(&conf.ExecInstall, "install", false, "Install as kustomize exec plugin")
	fs.BoolVar(&conf.ExecInstall, "i", false, "Install as kustomize exec plugin")
	fs.BoolVar(&conf.ExecInstall, "p", false, "Print kustomize plugin configuration")
	fs.BoolVar(&conf.ExecInstall, "plugin-conf", false, "Print kustomize plugin configuration")
	fs.StringVar(&conf.File, "f", "", "Input file")
	fs.StringVar(&conf.File, "file", "", "Input file")
	fs.BoolVar(&conf.Updir, "u", false, "Allow specifying .. in file paths")
	fs.BoolVar(&conf.Updir, "updir", false, "Allow specifying .. in file paths")
	fs.BoolVar(&conf.Links, "l", false, "Allow following symlinks in file paths")
	fs.BoolVar(&conf.Links, "links", false, "Allow following symlinks in file paths")
	fs.BoolVar(&conf.Links, "a", false, "Allow following symlinks in file paths")
	fs.BoolVar(&conf.Links, "abs", false, "Allow absolute file paths")
	err := fs.Parse(args)
	return err, conf
}

func Help() {
	fmt.Fprintf(os.Stderr, descstr)
	UsageFunc()
}

const descstr = "A Simple Include Transformer for YAML files --\n" +
	"Reads YAML resources from stdin or input files and performs\n" +
	"include substitutions.\n" +
	"see https://github.com/dudinea/yaml-include\n"

const usagestr = "\nUsage: \n" +
	"  %s [configfile] [options ...]\n" +
	"\n" +
	"Options:\n" +
	"  -h --help	       Print this usage message\n" +
	"  -i --install        Install as kustomize exec plugin\n" +
	"  -p --plugin-conf    Print kustomize plugin configuration file\n" +
	"  -f --file file.yaml Input file\n" +
	"  -u --up-dir         Allow specifying .. in file paths\n" +
	"  -l --links          Allow following symlinks in file paths\n" +
	"  -a --abs            Allow absolute paths in file paths\n" +
	"\n" +
	"Supported YAML include directives:\n" +
	"  foo!textfile: file.txt    -- include file.txt as text field\n" +
	"  bar:base64file: file.bin  -- include file.bin as base64 text\n"
