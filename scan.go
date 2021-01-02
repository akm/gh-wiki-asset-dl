package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/akm/gh-wiki-asset-dl/primitives"
)

const assetUrlBase = "https://user-images.githubusercontent.com/"

var assetExts = primitives.Strings{
	".gif", ".jpeg", ".jpg", ".mov", ".mp4", ".png",
	".docx", ".gz", ".log", ".pdf", ".pptx", ".txt",
	".xlsx", ".zip",
}

type PatternType int

const (
	ImgSrc PatternType = iota + 1
	MdImg
)

func (t PatternType) Quote(s string) string {
	switch t {
	case ImgSrc:
		return fmt.Sprintf(`"%s"`, s)
	case MdImg:
		return fmt.Sprintf(`(%s)`, s)
	default:
		panic(fmt.Errorf("Invalid PatternType: %v", t))
	}
}

type Pattern struct {
	Type   PatternType
	Regexp *regexp.Regexp
}

var patterns = func() []*Pattern {
	escapedBase := regexp.QuoteMeta(assetUrlBase)
	escapedExts := assetExts.Map(regexp.QuoteMeta).Join(`|`)
	imgSrcReStr := `"(` + escapedBase + `[^\"]*(?:` + escapedExts + `))"`
	mdImgReStr := `\((` + escapedBase + `[^\)]*(?:` + escapedExts + `))\)`
	return []*Pattern{
		{Type: ImgSrc, Regexp: regexp.MustCompile(imgSrcReStr)},
		{Type: MdImg, Regexp: regexp.MustCompile(mdImgReStr)},
	}
}()

type MatchResult struct {
	Type   PatternType
	Result string
}

func (mr *MatchResult) Replace(line string, replaced string) string {
	return strings.ReplaceAll(line, mr.Type.Quote(mr.Result), mr.Type.Quote(replaced))
}

func Scan(s string) []*MatchResult {
	r := []*MatchResult{}
	for _, ptn := range patterns {
		for _, m := range ptn.Regexp.FindAllStringSubmatch(s, -1) {
			r = append(r, &MatchResult{
				Type:   ptn.Type,
				Result: m[1],
			})
		}
	}
	return r
}
