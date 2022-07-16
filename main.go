package main

import (
	"os"

	"github.com/dudinea/kustomize-field-include/pkg/config"
	"github.com/dudinea/kustomize-field-include/pkg/kustomize"
	"github.com/dudinea/kustomize-field-include/pkg/transform"
)

func main() {
	var err error
	err, conf := config.ReadArgs(os.Args[1:])
	if nil != err {
		os.Exit(1)
	}
	if conf.PrintUsage {
		config.Help()
		os.Exit(1)
	}
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

	reader := os.Stdin
	transform.Transform(reader)
}
