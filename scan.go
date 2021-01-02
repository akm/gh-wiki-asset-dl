package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/akm/gh-wiki-asset-dl/primitives"
)

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

type MatchResult struct {
	Type   PatternType
	Result string
}

func (mr *MatchResult) Replace(line string, replaced string) string {
	return strings.ReplaceAll(line, mr.Type.Quote(mr.Result), mr.Type.Quote(replaced))
}

type Scanner struct {
	assetUrlBase string
}

func NewScanner(assetUrlBase string) *Scanner {
	return &Scanner{assetUrlBase: assetUrlBase}
}

func (s *Scanner) buildPatterns() []*Pattern {
	escapedBase := regexp.QuoteMeta(s.assetUrlBase)
	escapedExts := assetExts.Map(regexp.QuoteMeta).Join(`|`)
	imgSrcReStr := `"(` + escapedBase + `[^\"]*(?:` + escapedExts + `))"`
	mdImgReStr := `\((` + escapedBase + `[^\)]*(?:` + escapedExts + `))\)`
	return []*Pattern{
		{Type: ImgSrc, Regexp: regexp.MustCompile(imgSrcReStr)},
		{Type: MdImg, Regexp: regexp.MustCompile(mdImgReStr)},
	}
}

func (s *Scanner) Do(str string) []*MatchResult {
	r := []*MatchResult{}
	patterns := s.buildPatterns()
	for _, ptn := range patterns {
		for _, m := range ptn.Regexp.FindAllStringSubmatch(str, -1) {
			r = append(r, &MatchResult{
				Type:   ptn.Type,
				Result: m[1],
			})
		}
	}
	return r
}
