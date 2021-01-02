package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type Download struct {
	Url  string
	Path string
}

type Downloader struct {
	bufferSize int
}

func NewDownloader(bufferSize int) *Downloader {
	return &Downloader{bufferSize: bufferSize}
}

func (x *Downloader) Execute(dl *Download) error {
	cl := &http.Client{}
	resp, err := cl.Get(dl.Url)
	if err != nil {
		return err
	}
	if resp.StatusCode < 100 || resp.StatusCode >= 300 {
		return fmt.Errorf("Failed to download %s because of %s", dl.Url, resp.Status)
	}

	f, err := os.OpenFile(dl.Path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close file %s because of %v\n", dl.Path, err)
		}
	}()

	buf := make([]byte, x.bufferSize)
	if _, err := io.CopyBuffer(f, resp.Body, buf); err != nil {
		return err
	}

	return nil
}
