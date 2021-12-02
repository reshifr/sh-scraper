package scraper

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

const execlPath = "../data/execl.db"

type content struct {
	Title string
	Tags  []string
	Cmds  []string
}

func findExec(exec string) bool {
	db, err := sql.Open("sqlite3", execlPath)
	if err != nil {
		log.Fatalln("Validator: SQLite connection failed.")
	}
	defer db.Close()
	rows := db.QueryRow("SELECT * FROM execl WHERE execs = ?", exec)
	var val string
	err = rows.Scan(&val)
	if err != nil {
		return false
	}
	return true
}

func getCmds(body string) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil
	}
	var cmds []string
	doc.Find("code").Each(func(i int, val *goquery.Selection) {
		matches := regexp.MustCompile(`\r?\n`).
			FindAllStringIndex(val.Text(), -1)
		if matches != nil && len(matches) > 1 {
			return
		}
		cmd := regexp.MustCompile(`^[$|#|>] ?`).
			ReplaceAllString(val.Text(), "")
		cmd = regexp.MustCompile(`\r?\n`).ReplaceAllString(cmd, "")
		exec := regexp.MustCompile(`^\S*`).FindString(cmd)
		if findExec(exec) {
			cmds = append(cmds, cmd)
		}
	})
	return cmds
}

func getContent(id uint) (*content, bool) {
	var scrp scraper
	scrp.open()
	qst, miss := scrp.getQuestion(id)
	if miss {
		return nil, true
	}
	ans, miss := scrp.getAnswer(id)
	if miss {
		return nil, true
	}
	if qst == nil || ans == nil ||
		qst.Items == nil || ans.Items == nil ||
		len(qst.Items) == 0 || len(ans.Items) == 0 {
		return nil, false
	}
	cnt := &content{}
	cnt.Title = qst.Items[0].Title
	cnt.Tags = append(qst.Items[0].Tags)
	for _, item := range ans.Items {
		if item.Score != 0 {
			cnt.Cmds = append(cnt.Cmds, getCmds(item.Body)...)
		}
	}
	return cnt, false
}

func printContent(cnt *content) {
	if cnt == nil {
		return
	}
	fmt.Printf("%v\n", cnt.Title)
	for _, tag := range cnt.Tags {
		fmt.Printf("(%v) ", tag)
	}
	fmt.Println()
	for _, cmd := range cnt.Cmds {
		fmt.Printf("%v\n", cmd)
	}
	fmt.Println()
}

func Load() {
	result, miss := getContent(1)
	printContent(result)
	fmt.Println(":::", miss)
}
