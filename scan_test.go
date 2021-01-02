package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TesScan(t *testing.T) {
	t.Run("imgTag(1)", func(t *testing.T) {
		r := Scan(imgTag(1))
		assert.Len(t, r, 1)
		assert.Equal(t, ImgSrc, r[0].Type)
		assert.Equal(t, imgUrl(1), r[0].Result)
	})

	t.Run("imgTag(1) + imgTag(2)", func(t *testing.T) {
		r := Scan(imgTag(1) + " " + imgTag(2))
		assert.Len(t, r, 2)
		assert.Equal(t, ImgSrc, r[0].Type)
		assert.Equal(t, ImgSrc, r[1].Type)
		assert.Equal(t, imgUrl(1), r[0].Result)
		assert.Equal(t, imgUrl(2), r[1].Result)
	})

	t.Run("mdImg(3)", func(t *testing.T) {
		r := Scan(mdImg(3))
		assert.Len(t, r, 1)
		assert.Equal(t, MdImg, r[0].Type)
		assert.Equal(t, imgUrl(3), r[0].Result)
	})

	t.Run("imgTag(1)+mdImg(3)", func(t *testing.T) {
		r := Scan(imgTag(1) + " " + mdImg(3))
		assert.Len(t, r, 2)
		assert.Equal(t, ImgSrc, r[0].Type)
		assert.Equal(t, MdImg, r[1].Type)
		assert.Equal(t, imgUrl(1), r[0].Result)
		assert.Equal(t, imgUrl(3), r[1].Result)
	})
}
