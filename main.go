package main

import (
	"fmt"
	//	"flag"
	//"sigs.k8s.io/yaml"
	b64 "encoding/base64"
	"io"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

func writeBytes(bytes *[]byte) {
	_, err := os.Stdout.Write(*bytes)
	if nil != err {
		fmt.Fprintf(os.Stderr, "%v: Failed to write output: %v\n", os.Args[0], err.Error())
		os.Exit(5)
	}
}

var header = []byte("---\n")

// Include types
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
	bytes, err := os.ReadFile(path)
	if nil != err {
		return nil, err
	}
	return bytes, nil
}

func includeTextfile(v interface{}) ([]byte, error) {
	switch v.(type) {
	case string:
		return readFile(v.(string))
	default:
		return nil, fmt.Errorf("Invalid value for include specification: %v", reflect.TypeOf(v))
	}
}

func include(incl_type string, k string, v interface{}) (interface{}, error) {
	switch incl_type {
	case TEXTFILE:
		return includeTextfile(v)
	case BASE64FILE:
		data, err := includeTextfile(v)
		if nil != err {
			return "", err
		}
		return make([]byte, b64.StdEncoding.EncodedLen(len(data))), nil
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

func main() {
	args := os.Args
	progname := args[0]
	/*if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "%v: config file not provided!\nUsage: \n\t%v input_file\n", progname, progname)
		os.Exit(1)
	}
	filename := args[1]
	*/
	//reader, err := os.Open(filename)
	reader := os.Stdin

	// if nil != err {
	// 	fmt.Fprintf(os.Stderr, "%v: Failed to open input file %s: %v\n", progname, filename, err.Error())
	// 	os.Exit(2)
	// }
	var err error
	err = nil

	decoder := yaml.NewDecoder(reader)
	var m interface{}

	for err == nil {
		err = decoder.Decode(&m)
		if nil == err {
			err = processAny(m)
			if nil != err {
				//fmt.Fprintf(os.Stderr, "%v: Failed to process yaml: %v\n", progname, err.Error())
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
