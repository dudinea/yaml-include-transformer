package main

import (
	"fmt"
	//	"flag"
	"sigs.k8s.io/yaml"
	"os"
	//	"io"
)

func main() {

	args := os.Args
	progname := args[0]
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "%v: Input file not provided!\nUsage: \n\t%v input_file\n", progname, progname);
		os.Exit(1);
	}
	filename := args[1]
	//fmt.Printf("opening %s\n", args[1]);
	
	bytes,err  := os.ReadFile(filename)

	if nil != err {
		fmt.Fprintf(os.Stderr, "%v: Failed to read input file %s: %v\n", progname, filename, err.Error());
		os.Exit(2)
	}
	//var m map[string]interface{}
	var m interface{}
	err = yaml.Unmarshal(bytes, &m)

	if nil != err {
		fmt.Fprintf(os.Stderr, "%v: Failed to parse yaml file %s: %v\n", progname, filename, err.Error());
		os.Exit(3)
	}


	outBytes, err := yaml.Marshal(m)
	if nil != err {
		fmt.Fprintf(os.Stderr, "%v: Failed to convert to yaml: %v\n", progname, err.Error());
		os.Exit(5)
	}

	_, err = os.Stdout.Write(outBytes)
	if nil != err {
		fmt.Fprintf(os.Stderr, "%v: Failed to write output: %v\n", progname, err.Error());
		os.Exit(5)
	}
	
	
}
