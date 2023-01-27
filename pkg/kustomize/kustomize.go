package kustomize

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dudinea/yaml-include-transformer/pkg/config"
	"github.com/dudinea/yaml-include-transformer/pkg/transform"
)

func getPluginDir() (string, error) {
	if transform.Conf.Legacy {
		return getPluginDirLegacy()
	} else if transform.Conf.Exec {
		return getPluginDirKrm()
	} else {
		return "", fmt.Errorf("There is no need to install the binary when using a containerized KRM plugin.")
	}
}

func getPluginDirLegacy() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Failed to get user's home directory: %v", err)
	}
	return filepath.FromSlash(homeDir +
		"/.config/kustomize/plugin/" + config.ApiVersion + "/" +
		strings.ToLower(config.Progname)), nil
}

func getPluginDirKrm() (string, error) {
	return filepath.FromSlash("plugins/"), nil
}

func PluginInstall() error {
	pluginDir, err := getPluginDir()
	if nil != err {
		return err
	}
	fmt.Fprintf(os.Stderr, "Installing kustomize exec plugin %v\n", pluginDir)
	err = os.MkdirAll(pluginDir, os.ModePerm)
	if nil != err {
		return fmt.Errorf("Failed to create plugin directory: %v", err.Error())
	}
	targetPath := pluginDir + string(filepath.Separator) + config.Progname
	err = copyfile(os.Args[0], targetPath)
	// workaround for a Windows specific bug of some kustomize versions:
	//  it needs both files (with and without .exe extension)
	if nil == err && isWindows() {
		err = copyfile(os.Args[0], targetPath+".exe")
	}
	if nil != err {
		return fmt.Errorf("Failed to copy plugin: %v", err.Error())
	}
	return nil
}

func isWindows() bool {
	return runtime.GOOS == "windows"
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
	// can we assume that all not-windows means need chmod works?
	if !isWindows() {
		err = os.Chmod(dst, os.FileMode(0755))
		if err != nil {
			return fmt.Errorf("WARNING: Failed to make file executable: %v\n", err)
		}
	}

	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func PluginConf() {
	fmt.Printf(`---
# put this into a file in your kustomize directory
# and add filename to the list of transformers in
# the kustomization.yaml
apiVersion: %s
kind: %s
metadata:
  name: notImportantHere
`,
		config.ApiVersion, config.Progname)
	// FIXME: ugly - go templates? deserialize?
	if transform.Conf.Krm {
		fmt.Printf("  annotations:\n")
		if transform.Conf.Exec {
			fmt.Printf(`    config.kubernetes.io/function: |
      exec:
        path: ./plugins/%s
`, config.Progname)
		} else {
			fmt.Printf(`    config.kubernetes.io/function: |
      container:
        image: %s
`, transform.Conf.Dockertag)
		}

	}
}
