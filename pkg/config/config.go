package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Debug       bool
	PrintUsage  bool
	ExecInstall bool
	PluginConf  bool
	File        string
	Updir       bool
	Links       bool
	Abs         bool
	Version     bool
	Exec        bool
	Krm         bool
	Legacy      bool
	Dockertag   string
}

var Conf *Config

const Progname = "YamlIncludeTransformer"
const ApiVersion = "kustomize-utils.dudinea.org/v1"

var Dockertag string

var UsageFunc func()

func ReadArgs(args []string) (error, Config) {
	conf := Config{}
	fs := flag.NewFlagSet(Progname, flag.ContinueOnError)
	UsageFunc = func() {
		fmt.Fprintf(os.Stderr, usagestr, os.Args[0])
	}
	fs.Usage = UsageFunc
	fs.BoolVar(&conf.PrintUsage, "help", false, "Print help message")
	fs.BoolVar(&conf.PrintUsage, "h", false, "Print help message")
	fs.BoolVar(&conf.ExecInstall, "install", false, "Install as kustomize exec plugin")
	fs.BoolVar(&conf.ExecInstall, "i", false, "Install as kustomize exec plugin")
	fs.BoolVar(&conf.PluginConf, "p", false, "Print kustomize plugin configuration")
	fs.BoolVar(&conf.PluginConf, "plugin-conf", false, "Print kustomize plugin configuration")
	fs.StringVar(&conf.File, "f", "", "Input file")
	fs.StringVar(&conf.File, "file", "", "Input file")
	fs.BoolVar(&conf.Updir, "u", false, "Allow specifying .. in file paths")
	fs.BoolVar(&conf.Updir, "updir", false, "Allow specifying .. in file paths")
	fs.BoolVar(&conf.Links, "l", false, "Allow following symlinks in file paths")
	fs.BoolVar(&conf.Links, "links", false, "Allow following symlinks in file paths")
	fs.BoolVar(&conf.Abs, "a", false, "Allow absolute file paths")
	fs.BoolVar(&conf.Abs, "abs", false, "Allow absolute file paths")
	fs.BoolVar(&conf.Version, "v", false, "Print program version")
	fs.BoolVar(&conf.Version, "version", false, "Print program version")
	fs.BoolVar(&conf.Debug, "d", false, "Print debug messages on stderr")
	fs.BoolVar(&conf.Debug, "debug", false, "Print debug messages on stderr")
	fs.BoolVar(&conf.Exec, "E", false, "Exec style plugin")
	fs.BoolVar(&conf.Exec, "exec", false, "Exec style plugin")
	fs.BoolVar(&conf.Krm, "K", false, "KRM-function style plugin")
	fs.BoolVar(&conf.Krm, "krm", false, "KRM-function style plugin")
	fs.BoolVar(&conf.Legacy, "L", false, "Legacy style plugin")
	fs.BoolVar(&conf.Legacy, "legacy", false, "Legacy style plugin")
	fs.StringVar(&conf.Dockertag, "dockertag", Dockertag, "Docker tag of the KRM function")
	fs.StringVar(&conf.Dockertag, "D", Dockertag, "Docker tag of the KRM function")
	err := fs.Parse(args)
	return err, conf
}

func Help() {
	fmt.Fprintf(os.Stderr, descstr)
	UsageFunc()
}

const descstr = "A Simple Include Transformer for YAML files --\n" +
	"Reads YAML resources from stdin or an input file, and performs\n" +
	"include substitutions. Please see\n" +
	"https://github.com/dudinea/yaml-include-transformer"

const usagestr = "\nUsage: \n" +
	"  %s [configfile] [options ...]\n" +
	"\n" +
	"Options:\n" +
	"  -h --help	       Print this usage message\n" +
	"  -i --install        Install as kustomize exec plugin\n" +
	"  -p --plugin-conf    Print kustomize plugin configuration file\n" +
	"  -E --exec           Exec plugin (for -p and -i)\n" +
	"  -L --legacy         Legacy  plugin (for -p and -i), default\n" +
	"  -K --krm            KRM-function plugin (for -p and -i)\n" +
	"  -D --dockertag      KRM-function docker tag\n" +
	"  -f --file file.yaml Input file\n" +
	"  -u --up-dir         Allow specifying .. in file paths\n" +
	"  -l --links          Allow following symlinks in file paths\n" +
	"  -a --abs            Allow absolute paths in file paths\n" +
	"  -v --version        Print program version\n" +
	"  -d --debug          Print debug messages on stderr\n" +
	"\n" +
	"Supported YAML include directives:\n" +
	"  foo!textfile: file.txt    -- include file.txt as text field\n" +
	"  foo!base64file: file.bin  -- include file.bin as base64 text\n" +
	"  foo!jsonfile: file.json   -- deserialize and include file.json\n" +
	"  foo!yamlfile: file.yaml   -- deserialize and include file.yaml\n"
