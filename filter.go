package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/akm/gh-wiki-asset-dl/primitives"
)

type Filter struct {
	TargetExts primitives.Strings
}

func NewFilter(targetExts []string) *Filter {
	return &Filter{
		TargetExts: primitives.Strings(targetExts),
	}
}

func (ff *Filter) Glob(dir string, callback func(string) error) error {
	f, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, fi := range fileInfos {
		ext := path.Ext(fi.Name())
		if ff.TargetExts.Contains(ext) {
			if err := callback(fi.Name()); err != nil {
				return err
			}
		}
	}

	return nil
}
