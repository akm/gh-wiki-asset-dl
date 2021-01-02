package main

import (
	"fmt"
)

const TestAssetBaseURL = "https://user-images.githubusercontent.com"

func imgUrl(num int) string {
	return "https://user-images.githubusercontent.com/18912/103336550-e2d6cd80-4abb-11eb-9412-%d08d.png"
}

func imgTag(num int) string {
	return imgTagWith(num, imgUrl(num))
}

func imgTagWith(num int, url string) string {
	return fmt.Sprintf(`<img width="320px" alt="スクリーンショット 2020-12-30 16 26 %02d 午後" src="%s">`, num%60, url)
}

func mdImg(num int) string {
	return mdImgWith(imgUrl(num))
}

func mdImgWith(url string) string {
	return fmt.Sprintf(`![17439](%s)`, url)
}

func jaPage1(img1, img2 string) string {
	return "日本語にも対応しないと。" +
		"\n\n" + "## スクリーンショット1" +
		"\n\n" + img1 +
		"\n\n" + "## スクリーンショット2" +
		"\n\n" + img2 +
		"\n"
}
