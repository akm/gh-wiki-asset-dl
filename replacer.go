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
	Scanner    *Scanner
	Downloader *Downloader
	fileNo     int
}

func NewReplacer(scanner *Scanner, downloader *Downloader) *Replacer {
	return &Replacer{Scanner: scanner, Downloader: downloader}
}

func (rep *Replacer) Do(filepath string) error {
	rr, err := rep.LoadAndReplace(filepath)
	if err != nil {
		return err
	}

	for _, dl := range rr.downloads {
		if err := rep.Downloader.Execute(dl); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(filepath, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close file %s because of %v\n", filepath, err)
		}
	}()

	if _, err := f.Write(rr.output.Bytes()); err != nil {
		return err
	}

	return nil
}

type ReplaceResult struct {
	output    *bytes.Buffer
	downloads []*Download
}

func (rep *Replacer) LoadAndReplace(filepath string) (*ReplaceResult, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	re, err := regexp.Compile(regexp.QuoteMeta(path.Ext(filepath)) + `\z`)
	if err != nil {
		return nil, err
	}
	dest := re.ReplaceAllString(filepath, "")

	return rep.Replace(f, dest)
}

func (rep *Replacer) Replace(input io.Reader, baseDir string) (*ReplaceResult, error) {
	downloads := []*Download{}

	output := bytes.NewBufferString("")

	sc := bufio.NewScanner(input)
	for sc.Scan() {
		line := sc.Text()
		mrs := rep.Scanner.Do(line)
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
