package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
)

type Replacer struct {
	fileNo int
}

func NewReplacer() *Replacer {
	return &Replacer{}
}

func (rep *Replacer) Do(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	re, err := regexp.Compile(regexp.QuoteMeta(path.Ext(filepath)) + `\z`)
	if err != nil {
		return err
	}
	dest := re.ReplaceAllString(filepath, "")

	rr, err := rep.Replace(f, dest)
	if err != nil {
		return err
	}

	fmt.Printf("%v", rr)

	return nil
}

type Download struct {
	Url  string
	Path string
}

type ReplaceResult struct {
	output    *bytes.Buffer
	downloads []*Download
}

func (rep *Replacer) Replace(input io.Reader, baseDir string) (*ReplaceResult, error) {
	downloads := []*Download{}

	output := bytes.NewBufferString("")

	sc := bufio.NewScanner(input)
	for sc.Scan() {
		line := sc.Text()
		mrs := Scan(line)
		for _, mr := range mrs {
			rep.fileNo++
			dest := fmt.Sprintf("%s/%04d%s", baseDir, rep.fileNo, path.Ext(mr.Result))
			downloads = append(downloads, &Download{Url: mr.Result, Path: dest})
			line = mr.Replace(line, dest)
		}
		output.WriteString(line)
		output.WriteString("\n")
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	return &ReplaceResult{output: output, downloads: downloads}, nil
}

// func (rep *Replacer) makeText(filepath string) (*bytes.Buffer, error) {
// 	var r bytes.Buffer

// 	f, err := os.Open(filepath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()

// 	sc := bufio.NewScanner(f)
// 	for sc.Scan() {
// 		matches := sc.Text()

// 	}
// 	if err := sc.Err(); err != nil {
// 		return nil, err
// 	}

// 	return &r, nil
// }
