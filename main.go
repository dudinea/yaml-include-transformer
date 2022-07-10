package main

import (
	"fmt"
	//	"flag"
	//"sigs.k8s.io/yaml"
	b64 "encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"io"
	"gopkg.in/yaml.v3"
	"flag"
	"runtime"
)

func writeBytes(bytes *[]byte) {
	_, err := os.Stdout.Write(*bytes)
	if nil != err {
		fmt.Fprintf(os.Stderr, "%v: Failed to write output: %v\n", os.Args[0], err.Error())
		os.Exit(5)
	}
}

var header = []byte("---\n")

// Include suffixes
const TEXTFILE = "!textfile"
const BASE64FILE = "!base64file"

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
		data, err =  readFile(v.(string))
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
		encoded:=make([]byte, b64.StdEncoding.EncodedLen(len(data)))
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
		//fmt.Fprintf(os.Stderr, "switch: data is MAP\n")
		return processMap(data.(map[string]interface{}))
	case []interface{}:
		return processArray(data.([]interface{}))
		//fmt.Fprintf(os.Stderr, "switch: data is ARRAY\n")

	}
	return nil
}


func getPluginDir() string {
	homeDir, err :=os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get user's home directory\n")
		os.Exit(2)
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
        return 	err
}



func main() {
	var err error
	args := os.Args
	progname := args[0]

	fs := flag.NewFlagSet("KustomizeFieldInclude", flag.ExitOnError)
	fs.SetOutput(os.Stderr)

	printUsage:=false
	execInstall:=false
	flag.BoolVar(&printUsage, "h", false, "Print usage")
	flag.BoolVar(&execInstall, "i", false, "Install exec plugin")
	flag.Parse()
	
	if (printUsage) {
		flag.Usage()
		os.Exit(1)
	}
	if (execInstall) {
		pluginDir := getPluginDir()
		fmt.Fprintf(os.Stderr, "Installing kustomize exec plugin %v\n", pluginDir)
		err := os.MkdirAll(pluginDir,  os.ModePerm)
		if nil != err {
			fmt.Fprintf(os.Stderr, "Failed to create plugin directory: %v\n", err.Error())
		}
		err = copyfile(args[0], pluginDir + string(filepath.Separator) + "FieldIncludeTransformer")
		if nil == err {
			fmt.Fprintf(os.Stderr, "Installation complete\n")
			os.Exit(0)
		} else {
			fmt.Fprintf(os.Stderr, "Failed to copy plugin: %v\n", err.Error())
			os.Exit(2)
		}
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
				fmt.Fprintf(os.Stderr, "%v: Failed to process yaml: %v\n", progname, err.Error())
				os.Exit(5)
			}
			outBytes, err := yaml.Marshal(m)
			if nil != err {
				fmt.Fprintf(os.Stderr, "%v: Failed to convert to yaml: %v\n", progname, err.Error())
				os.Exit(5)
			}
			writeBytes(&header)
			writeBytes(&outBytes)
		}
	}

	if err != io.EOF {
		fmt.Fprintf(os.Stderr, "%v: Error reading input stream: %v\n", progname, err.Error())
		os.Exit(3)
	}

}
