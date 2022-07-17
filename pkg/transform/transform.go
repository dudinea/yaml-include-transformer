package transform

import (
	b64 "encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dudinea/yaml-include-transformer/pkg/config"
	"gopkg.in/yaml.v3"
)

var header = []byte("---\n")

const TEXTFILE = "!textfile"
const BASE64FILE = "!base64file"

var Conf *config.Config

func Transform(reader *os.File) {
	var err error = nil

	decoder := yaml.NewDecoder(reader)
	var m interface{}

	for err == nil {
		err = decoder.Decode(&m)
		if nil == err {
			err = processAny(m)
			if nil != err {
				Errexit(5, "Failed to data: %v", err.Error())
			}
			outBytes, err := yaml.Marshal(m)
			if nil != err {
				Errexit(5, "Failed to convert to yaml: %v", err.Error())
			}
			writeBytes(&header)
			writeBytes(&outBytes)
		}
	}
	if err != io.EOF {
		Errexit(3, "Error decoding input stream: %v", err.Error())
	}
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
	checkPath(path)
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

func checkPath(path string) {
	// 1. check if path looks like absolute
	checkAbsPath(path)
	checkUpDir(path)
	// Resolve synlinks if any
	resolved := resolveAndCheckLinks(path)
	// Check if resolved path became absolute
	checkAbsPath(resolved)
	checkUpDir(resolved)
}

func resolveAndCheckLinks(path string) string {
	resolved, err := filepath.EvalSymlinks(path)
	if nil != err {
		Errexit(6, "Error: invalid path '%s'", path)
	}
	// test if not equal to original
	// (EvalSymLinks calls Clean() before return)
	if !Conf.Links && resolved != filepath.Clean(path) {
		Errexit(6, "Error: path '%s' contains symlinks", path)
	}
	return resolved
}

func checkAbsPath(path string) {
	if !Conf.Abs {
		if filepath.IsAbs(path) {
			Errexit(6, "Error: absolute file path '%s' is not allowed")
		}
	}
}

func checkUpDir(path string) {
	if !Conf.Updir {
		platformPath := filepath.FromSlash(path)
		parts := strings.Split(platformPath, string(os.PathSeparator))
		for _, v := range parts {
			if v == ".." {
				Errexit(6, "Error: absolute file path '%s' is not allowed")
			}
		}
	}
}

func processMap(m map[string]interface{}) error {
	for k, v := range m {
		incl_type, new_key := isInclude(k)
		if incl_type != "" {
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

func Errexit(code int, msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: ", os.Args[0])
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(code)
}

func writeBytes(bytes *[]byte) {
	_, err := os.Stdout.Write(*bytes)
	if nil != err {
		Errexit(5, "Failed to write output: %v", err.Error())
	}
}