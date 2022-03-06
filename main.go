package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	// url := "https://raw.githubusercontent.com/jobbole/awesome-go-cn/master/README.md"
	urlCN := "https://gitee.com/muyegit/awesome-go-cn/raw/master/README.md"
	get, err := http.Get(urlCN)

	if err != nil {
		return
	}
	defer get.Body.Close()

	// 结果写入新的 README 文件
	f, err := os.Create("./README.md")
	defer f.Close()
	w := bufio.NewWriter(f)

	br := bufio.NewReader(get.Body)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		lineStr := string(line)
		ss := regexp.MustCompile(`\([\s\S]*?\)`)
		lineS := ss.FindString(lineStr)
		lineS = strings.TrimLeft(lineS, "(")
		lineS = strings.TrimRight(lineS, ") :")
		// 以 github.com 开头的包
		IfGitHub := strings.HasPrefix(lineS, "https://github.com/")
		if IfGitHub {
			star := FindStar(lineS)
			fmt.Println(lineS, star)
			lineStr += "(" + star + ") \n"
		}

		_, err = w.WriteString(lineStr)
		w.Flush()
	}

}

// FindStar 获取 github 项目的 star 数
func FindStar(link string) string {
	result := ""
	c := colly.NewCollector()
	c.OnHTML("#repo-stars-counter-star", func(e *colly.HTMLElement) {
		result = e.Text
	})
	_ = c.Visit(link)
	return result
}
