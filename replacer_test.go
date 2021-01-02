package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplacerReplace(t *testing.T) {
	t.Run("jaPage1", func(t *testing.T) {
		input := bytes.NewBufferString(jaPage1(imgTag(1), mdImg(3)))
		rep := NewReplacer(NewScanner(TestAssetBaseURL))
		rr, err := rep.Replace(input, "日本語ページ1")

		path1 := "日本語ページ1/0001.png"
		path2 := "日本語ページ1/0002.png"
		assert.NoError(t, err)
		assert.Len(t, rr.downloads, 2)
		assert.Equal(t, imgUrl(1), rr.downloads[0].Url)
		assert.Equal(t, path1, rr.downloads[0].Path)
		assert.Equal(t, imgUrl(3), rr.downloads[1].Url)
		assert.Equal(t, path2, rr.downloads[1].Path)

		assert.Equal(t, jaPage1(imgTagWith(1, path1), mdImgWith(path2)), rr.output.String())
	})
}
