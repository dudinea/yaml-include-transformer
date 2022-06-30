package main

import (
	"fmt"
	//	"flag"
	//"sigs.k8s.io/yaml"
	"gopkg.in/yaml.v3"
	"os"
	"io"
)


func writeBytes(bytes *[]byte) {
	_, err := os.Stdout.Write(*bytes)
	if nil != err {
		fmt.Fprintf(os.Stderr, "%v: Failed to write output: %v\n", os.Args[0], err.Error());
		os.Exit(5)
	}
}

var header = []byte("---\n");

func main() {

	args := os.Args
	progname := args[0]
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%v: Input file not provided!\nUsage: \n\t%v input_file\n", progname, progname);
		os.Exit(1);
	}
	filename := args[1]
	//fmt.Printf("opening %s\n", args[1]);

	reader, err := os.Open(filename)


	if nil != err {
		fmt.Fprintf(os.Stderr, "%v: Failed to open input file %s: %v\n", progname, filename, err.Error());
		os.Exit(2)
	}

	decoder := yaml.NewDecoder(reader)
	//var m map[string]interface{}
	var m interface{}

	
	for err == nil {
		err = decoder.Decode(&m)
		if nil == err {
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
