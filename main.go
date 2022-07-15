package main

import (
	b64 "encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

var header = []byte("---\n")

// Include suffixes
const TEXTFILE = "!textfile"
const BASE64FILE = "!base64file"

const helpstr = ":\nAn Simple Include Transformer for YAML files --\n" +
	"Reads YAML resources from stdin or input files and performs\n" +
	"include substitutions.\n" +
	"see https://github.com/dudinea/yaml-include\n" +
	"\n" +
	"Usage: \n" +
	"  %s [configfile] [options ...]\n" +
	"\n" +
	"Options:\n" +
	"  -h --help	     Print this usage message\n" +
	"  -i --install      Install as kustomize exec plugin\n" +
	"  -p --plugin-conf  Print kustomize plugin configuration file\n" +
	"  -f --file file    Input file\n" +
	"  -u --up-dir       Allow specifying .. in file paths\n" +
	"  -l --links        Allow following symlinks .. in file paths\n" +
	"\n" +
	"Supported YAML include directives:\n" +
	"  foo!textfile: file.txt    -- include file.txt as text field\n" +
	"  bar:base64file: file.bin  -- include file.bin as base64 text\n"

func main() {
	var err error
	args := os.Args
	//progname := args[0]

	fs := flag.NewFlagSet("FieldIncludePlugin", flag.ExitOnError)
	fs.SetOutput(os.Stderr)

	printUsage := false
	execInstall := false
	flag.BoolVar(&printUsage, "help", false, "Print usage")
	flag.BoolVar(&printUsage, "h", false, "Print usage")
	flag.BoolVar(&execInstall, "install", false, "Install exec plugin")

	flag.Parse()

	if printUsage {
		errexit(1, helpstr, os.Args[0])
	}
	if execInstall {
		pluginDir := getPluginDir()
		fmt.Fprintf(os.Stderr, "Installing kustomize exec plugin %v\n", pluginDir)
		err := os.MkdirAll(pluginDir, os.ModePerm)
		if nil != err {
			errexit(2, "Failed to create plugin directory: %v", err.Error())
		}
		err = copyfile(args[0], pluginDir+string(filepath.Separator)+"FieldIncludeTransformer")
		if nil != err {
			errexit(2, "Failed to copy plugin: %v", err.Error())
		}
		errexit(0, "Installation complete")
	}
	reader := os.Stdin
	err = nil

	decoder := yaml.NewDecoder(reader)
	var m interface{}

	for err == nil {
		err = decoder.Decode(&m)
		if nil == err {
			err = processAny(m)
			if nil != err {
				errexit(5, "Failed to process yaml: %v", err.Error())
			}
			outBytes, err := yaml.Marshal(m)
			if nil != err {
				errexit(5, "Failed to convert to yaml: %v", err.Error())
			}
			writeBytes(&header)
			writeBytes(&outBytes)
		}
	}

	if err != io.EOF {
		errexit(3, "Error reading input stream: %v", err.Error())
	}

}

func usage() {

}

// return include type and original key
func isInclude(k string) (include_type string, new_key string) {
	//fmt.Fprintf(os.Stderr, "%v: isInclude: %v\n", os.Args[0], k)
	if strings.HasSuffix(k, TEXTFILE) {
		return TEXTFILE, k[0 : len(k)-len(TEXTFILE)]
	} else if strings.HasSuffix(k, BASE64FILE) {
		return BASE64FILE, k[0 : len(k)-len(BASE64FILE)]
	}
	return "", k
}

// FIXME: security: disallow absolute paths and other tricks
func readFile(path string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(path)
	if nil != err {
		return nil, err
	}
	return bytes, nil
}

func includeTextfile(v interface{}) (string, error) {
	data, err := includeFile(v)
	if nil != err {
		return "", err
	}
	return string(data), err
}

func includeFile(v interface{}) ([]byte, error) {
	var data []byte
	var err error
	switch v.(type) {
	case string:
		data, err = readFile(v.(string))
	default:
		err = fmt.Errorf("Invalid value for include specification: %v", reflect.TypeOf(v))
	}
	if nil != err {
		return nil, err
	}
	return data, err
}

func include(incl_type string, k string, v interface{}) (interface{}, error) {
	switch incl_type {
	case TEXTFILE:
		return includeTextfile(v)
	case BASE64FILE:
		data, err := includeFile(v)
		if nil != err {
			return "", err
		}
		encoded := make([]byte, b64.StdEncoding.EncodedLen(len(data)))
		b64.StdEncoding.Encode(encoded, data)
		return string(encoded), nil
	default:
		return v, fmt.Errorf("Internal error: invalid include type %s", incl_type)
	}
}

func processMap(m map[string]interface{}) error {
	for k, v := range m {
		//fmt.Fprintf(os.Stderr, "%v: process map: %v: %v\n", os.Args[0], k, v)
		incl_type, new_key := isInclude(k)
		if incl_type != "" {
			//fmt.Fprintf(os.Stderr, "%v: include: %v: %v\n", os.Args[0], k, v)
			newval, err := include(incl_type, k, v)
			if err != nil {
				return err
			}
			m[new_key] = newval
			delete(m, k)
		} else {
			processAny(v)
		}
	}
	return nil
}

func processArray(a []interface{}) error {
	for _, k := range a {
		err := processAny(k)
		if nil != err {
			return err
		}
	}
	return nil
}

func processAny(data interface{}) error {
	switch data.(type) {
	case map[string]interface{}:
		return processMap(data.(map[string]interface{}))
	case []interface{}:
		return processArray(data.([]interface{}))
	}
	return nil
}

func getPluginDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		errexit(2, "Failed to get user's home directory")
	}
	return filepath.FromSlash(homeDir + "/.config/kustomize/plugin/kustomize-utils.dudinea.org/v1/fieldincludetransformer")
}

func copyfile(src string, dst string) error {
	fmt.Fprintf(os.Stderr, "copy '%v' to '%v'\n", src, dst)
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}

	if runtime.GOOS != "windows" {
		err = os.Chmod(dst, os.FileMode(0755))
		if err != nil {
			return fmt.Errorf("WARNING: Failed to make file executable: %v\n", err)
		}
	}

	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func errexit(code int, msg string, args ...interface{}) {
	fmt.Fprint(os.Stderr, os.Args[0])
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(code)
}

func writeBytes(bytes *[]byte) {
	_, err := os.Stdout.Write(*bytes)
	if nil != err {
		errexit(5, "Failed to write output: %v", err.Error())
	}
}
