package main

import (
	"fmt"
	//	"flag"
	//"sigs.k8s.io/yaml"
	"gopkg.in/yaml.v3"
	"os"
	"io"
	"strings"
	"reflect"
)


func writeBytes(bytes *[]byte) {
	_, err := os.Stdout.Write(*bytes)
	if nil != err {
		fmt.Fprintf(os.Stderr, "%v: Failed to write output: %v\n", os.Args[0], err.Error());
		os.Exit(5)
	}
}

var header = []byte("---\n");

const TEXTFILE = "!textfile"

func isInclude(k string) string {
	fmt.Fprintf(os.Stderr, "%v: isInclude: %v\n", os.Args[0],  k);
	if strings.HasSuffix(k, TEXTFILE) {
		return TEXTFILE;
	}
	return "";
}


func includeTextfile(v interface{}) (interface{}, error) {
	switch v.(type) {
	case string:
		return "<FILE>", nil
	default:
		return nil, fmt.Errorf("Invalid value for include specification: %v", reflect.TypeOf(v));
	}
	
}


func include(incl_type string, k string, v interface{}) (interface{}, error) {
	switch incl_type {
	case TEXTFILE:
		return includeTextfile(v);
	default:
		return v, fmt.Errorf("Internal error: invalid include type %s", incl_type);
	}
	
}

func processMap(m map[string]interface{}) error {
	for k, v := range m {
		fmt.Fprintf(os.Stderr, "%v: process map: %v: %v\n", os.Args[0],  k, v);
		incl_type := isInclude(k)
		if incl_type != "" {
			fmt.Fprintf(os.Stderr, "%v: include: %v: %v\n", os.Args[0],  k, v);
			newval, err := include(incl_type, k, v)
			if err != nil {
				return err;
			}
			m[k]=newval;
		} else {
			processAny(v)
		}
	}
	return nil
}


func processAny(data interface{}) error {
	switch  data.(type) {
	case map[string]interface{}:
		fmt.Fprintf(os.Stderr, "switch: data is MAP\n");
		return processMap(data.(map[string]interface{}))
	case []interface{}:
		fmt.Fprintf(os.Stderr, "switch: data is ARRAY\n");
		
	}
	return nil
}



func main() {
	args := os.Args
	progname := args[0]
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%v: Input file not provided!\nUsage: \n\t%v input_file\n", progname, progname);
		os.Exit(1);
	}
	filename := args[1]
	reader, err := os.Open(filename)

	if nil != err {
		fmt.Fprintf(os.Stderr, "%v: Failed to open input file %s: %v\n", progname, filename, err.Error());
		os.Exit(2)
	}

	decoder :=  yaml.NewDecoder(reader)
	var m interface{}
	
	for err == nil {
		err = decoder.Decode(&m)
		if nil == err {
			err = processAny(m)
			if nil != err {
				fmt.Fprintf(os.Stderr, "%v: Failed to process yaml: %v\n", progname, err.Error());
				os.Exit(5)
			}			
			outBytes, err := yaml.Marshal(m)
			if nil != err {
				fmt.Fprintf(os.Stderr, "%v: Failed to convert to yaml: %v\n", progname, err.Error());
				os.Exit(5)
			}
			writeBytes(&header);
			writeBytes(&outBytes);
		}
	}
	
	if err != io.EOF {
		fmt.Fprintf(os.Stderr, "%v: Failed to read yaml file %s: %v\n", progname, filename, err.Error());
		os.Exit(3)
	}
	
}
