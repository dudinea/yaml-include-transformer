package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dudinea/yaml-include-transformer/pkg/config"
	"github.com/dudinea/yaml-include-transformer/pkg/kustomize"
	"github.com/dudinea/yaml-include-transformer/pkg/transform"
)

// these two are set from ldflags
var version string
var dockertag string

func main() {
	config.Dockertag = dockertag
	err, conf := config.ReadArgs(os.Args[1:])
	if nil != err {
		transform.Errexit(1, "Invalid command line arguments: %v", err)
	}
	if conf.Debug {
		log.Printf("run with args %v\n", os.Args)
	}

	if conf.Version {
		fmt.Println(version)
		os.Exit(0)
	}
	if conf.PrintUsage {
		config.Help()
		os.Exit(1)
	}
	if conf.Legacy && conf.Krm {
		transform.Errexit(1, "Invalid arguments given: %v",
			fmt.Errorf("Cannot select both --legacy (-L) and --krm (-K)"))
	}
	if !conf.Legacy && !conf.Krm {
		conf.Legacy = true
	}

	if conf.Legacy {
		conf.Exec = true
	}
	config.Conf = &conf

	if conf.ExecInstall {
		err = kustomize.PluginInstall()
		if nil != err {
			transform.Errexit(2, "Kustomize exec plugin installation failed: %v", err)
		} else {
			transform.Errexit(0, "Kustomize exec plugin Installation complete")
		}
	}

	if conf.PluginConf {
		kustomize.PluginConf()
		os.Exit(0)
	}

	if len(conf.Files) == 0 {
		if conf.Debug {
			log.Println("using stdin as input")
		}
		reader := os.Stdin
		transform.Transform(reader)
	} else {
		for idx := 0; idx < len(conf.Files); idx++ {
			if idx > 0 {
				fmt.Fprintln(os.Stdout, "---")
			}
			transform.TransformFileOrDir(conf.Files[idx], true)
		}
	}
}
