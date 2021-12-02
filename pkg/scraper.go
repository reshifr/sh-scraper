package scraper

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type scraper struct {
	client *http.Client
}

type question struct {
	Items []struct {
		Tags  []string `json:"tags"`
		Owner struct {
			AccountID    int    `json:"account_id"`
			Reputation   int    `json:"reputation"`
			UserID       int    `json:"user_id"`
			UserType     string `json:"user_type"`
			AcceptRate   int    `json:"accept_rate"`
			ProfileImage string `json:"profile_image"`
			DisplayName  string `json:"display_name"`
			Link         string `json:"link"`
		} `json:"owner"`
		IsAnswered       bool   `json:"is_answered"`
		ViewCount        int    `json:"view_count"`
		AcceptedAnswerID int    `json:"accepted_answer_id"`
		AnswerCount      int    `json:"answer_count"`
		Score            int    `json:"score"`
		LastActivityDate int    `json:"last_activity_date"`
		CreationDate     int    `json:"creation_date"`
		LastEditDate     int    `json:"last_edit_date"`
		QuestionID       int    `json:"question_id"`
		ContentLicense   string `json:"content_license"`
		Link             string `json:"link"`
		Title            string `json:"title"`
		Body             string `json:"body"`
	} `json:"items"`
	HasMore        bool `json:"has_more"`
	QuotaMax       int  `json:"quota_max"`
	QuotaRemaining int  `json:"quota_remaining"`
}

type answer struct {
	Items []struct {
		Owner struct {
			AccountID    int    `json:"account_id"`
			Reputation   int    `json:"reputation"`
			UserID       int    `json:"user_id"`
			UserType     string `json:"user_type"`
			ProfileImage string `json:"profile_image"`
			DisplayName  string `json:"display_name"`
			Link         string `json:"link"`
		} `json:"owner"`
		IsAccepted       bool   `json:"is_accepted"`
		Score            int    `json:"score"`
		LastActivityDate int    `json:"last_activity_date"`
		CreationDate     int    `json:"creation_date"`
		AnswerID         int    `json:"answer_id"`
		QuestionID       int    `json:"question_id"`
		ContentLicense   string `json:"content_license"`
		Body             string `json:"body"`
		LastEditDate     int    `json:"last_edit_date,omitempty"`
	} `json:"items"`
	HasMore        bool `json:"has_more"`
	QuotaMax       int  `json:"quota_max"`
	QuotaRemaining int  `json:"quota_remaining"`
}

func (this *scraper) open() {
	proxy, err := url.Parse("socks5://127.0.0.1:9050")
	if err != nil {
		log.Fatalln(err)
	}
	this.client = &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxy)},
	}
}

func (this *scraper) get(endPoint, source string) ([]byte, bool) {
	req, err := http.NewRequest("GET", endPoint, nil)
	if err != nil {
		log.Println("Scraper: HTTP request failed.")
		return nil, true
	}
	query := req.URL.Query()
	query.Set("order", "desc")
	query.Set("sort", "votes")
	query.Set("site", source)
	query.Set("filter", "withbody")
	req.URL.RawQuery = query.Encode()
	resp, err := this.client.Do(req)
	if err != nil {
		log.Println("Scraper: Connection timeout.")
		return nil, true
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, true
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Scraper: Buffer reading failed.")
		return nil, false
	}
	return []byte(html.UnescapeString(string(respBody))), false
}

func (this *scraper) getQuestion(id uint) (*question, bool) {
	endPoint := fmt.Sprintf(
		"https://api.stackexchange.com/2.3/questions/%v", id)
	data, miss := this.get(endPoint, "unix.stackexchange.com")
	if miss || data == nil {
		return nil, true
	}
	qst := &question{}
	err := json.Unmarshal(data, qst)
	if err != nil {
		log.Println("Scraper: JSON parsing failed.")
		return nil, false
	}
	return qst, false
}

func (this *scraper) getAnswer(id uint) (*answer, bool) {
	endPoint := fmt.Sprintf(
		"https://api.stackexchange.com/2.3/questions/%v/answers", id)
	data, miss := this.get(endPoint, "unix.stackexchange.com")
	if miss || data == nil {
		return nil, true
	}
	ans := &answer{}
	err := json.Unmarshal(data, ans)
	if err != nil {
		log.Println("Scraper: JSON parsing failed.")
		return nil, false
	}
	return ans, false
}
