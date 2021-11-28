package scraper

import (
	"fmt"
	"html"
	"regexp"
)

var (
	bodyReplaceNL   *regexp.Regexp = regexp.MustCompile(`[\r|]*\n`)
	cmdFindRegex    *regexp.Regexp = regexp.MustCompile(`\<code\>.*\<\/code\>`)
	cmdReplaceRegex *regexp.Regexp = regexp.MustCompile(`\<[\/|]*code>`)
)

type Validator struct {
	Tags     []string
	Title    string
	Commands []string
}

func getRawCommand(body string) []string {
	body = regexp.MustCompile(`\r?\n`).ReplaceAllLiteralString(body, "")
	rawCmdList := regexp.MustCompile(`<code>(.*?)<\/code>`).FindAllString(body, -1)
	for i, rawCmd := range rawCmdList {
		rawCmdList[i] = html.UnescapeString(
			regexp.MustCompile(
				`<[\/|]*code>[\$|\#]*[ |]*`).
				ReplaceAllLiteralString(rawCmd, ""))
	}
	return rawCmdList
}

func (this *Validator) Load() {
	var scrp scraper
	ans := scrp.getAnswer(254335, unixSource)
	var cmdList []string
	for _, item := range ans.Items {
		cmdList = append(cmdList, getRawCommand(item.Body)...)
	}

	for _, item := range cmdList {
		fmt.Println(item)
	}
}
