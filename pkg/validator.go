package scraper

import (
	"database/sql"
	"encoding/binary"
	"html"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

const (
	shlPath   = "data/shl.db"
	countPath = "data/count.dat"
)

type content struct {
	Title string
	Tags  []string
	Cmds  []string
}

type validator struct {
	Scp    scraper
	Db     *sql.DB
	Stream *os.File
}

func (this *validator) open() {
	this.Scp.open()
	db, err := sql.Open("sqlite3", shlPath)
	if err != nil {
		log.Fatalln("Validator: SQLite connection failed.")
	}
	this.Db = db
	stream, err := os.OpenFile(countPath,
		os.O_CREATE|os.O_SYNC|os.O_RDWR, 0777)
	if err != nil {
		log.Fatalln("Validator: Log connection failed.")
	}
	this.Stream = stream
}

func (this *validator) close() {
	defer this.Db.Close()
	defer this.Stream.Close()
}

func (this *validator) findExec(exec string) bool {
	rows := this.Db.QueryRow(`
		SELECT * FROM execl
		WHERE execs = ?`,
		exec,
	)
	var val string
	err := rows.Scan(&val)
	if err != nil {
		return false
	}
	return true
}

func (this *validator) getCmds(body string) []string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil
	}
	var cmds []string
	doc.Find("code").Each(func(i int, val *goquery.Selection) {
		unescp := html.UnescapeString(val.Text())
		matches := regexp.MustCompile(`\r?\n`).FindAllStringIndex(unescp, -1)
		if matches != nil && len(matches) > 1 {
			return
		}
		cmd := regexp.MustCompile(`^[$|#|>] ?`).ReplaceAllString(unescp, "")
		cmd = regexp.MustCompile(`\r?\n`).ReplaceAllString(cmd, "")
		exec := regexp.MustCompile(`^\S*`).FindString(cmd)
		if this.findExec(exec) {
			cmds = append(cmds, cmd)
		}
	})
	return cmds
}

func (this *validator) getContent(id uint64) (*content, bool) {
	qst, miss := this.Scp.getQuestion(id)
	if miss {
		return nil, true
	}
	ans, miss := this.Scp.getAnswer(id)
	if miss {
		return nil, true
	}
	if qst == nil || ans == nil ||
		qst.Items == nil || ans.Items == nil ||
		len(qst.Items) == 0 || len(ans.Items) == 0 {
		return nil, false
	}
	cnt := &content{}
	cnt.Title = html.UnescapeString(qst.Items[0].Title)
	for _, tag := range qst.Items[0].Tags {
		cnt.Tags = append(cnt.Tags, html.UnescapeString(tag))
	}
	for _, item := range ans.Items {
		if item.Score != 0 {
			cnt.Cmds = append(cnt.Cmds, this.getCmds(item.Body)...)
		}
	}
	return cnt, false
}

func (this *validator) saveContent(cnt *content) {
	if cnt == nil {
		return
	}
	this.Db.Exec(`
		INSERT INTO titlel(titles)
		VALUES(?)`,
		cnt.Title,
	)
	for _, tag := range cnt.Tags {
		this.Db.Exec(`
			INSERT INTO tagl(titles, tags)
			VALUES(?, ?)`,
			cnt.Title, tag,
		)
	}
	for _, cmd := range cnt.Cmds {
		this.Db.Exec(`
			INSERT INTO cmdl(titles, cmds)
			VALUES(?, ?)`,
			cnt.Title, cmd,
		)
	}
}

func (this *validator) getCount() uint64 {
	var cast [8]byte
	_, err := this.Stream.ReadAt(cast[:], 0)
	if err != nil {
		err = this.Stream.Truncate(8)
		if err != nil {
			log.Fatalln("Validator: Log reading failed.")
		}
		return 0
	}
	return binary.LittleEndian.Uint64(cast[:])
}

func (this *validator) saveCount(count uint64) {
	var cast [8]byte
	binary.LittleEndian.PutUint64(cast[:], count)
	this.Stream.WriteAt(cast[:], 0)
}

func Scrap(max uint64, onSucceed func(), onFailed func()) {
	var valdr validator
	valdr.open()
	defer valdr.close()
	for {
		count := valdr.getCount()
		if count == max {
			break
		}
		cnt, miss := valdr.getContent(count)
		if miss {
			log.Printf("Scraping (%v) failed.", count)
			onFailed()
			continue
		}
		valdr.saveContent(cnt)
		valdr.saveCount(count + 1)
		log.Printf("Scraped (%v) succeed", count)
		onSucceed()
	}
	log.Println("Scraping successfully completed.")
}
