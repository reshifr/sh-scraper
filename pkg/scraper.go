package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	defaultTimeout = time.Duration(8 * time.Second)
	unixSource     = "unix.stackexchange.com"
)

type scraper struct {
	client http.Client
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

func (this *scraper) init(timeout time.Duration) {
	this.client = http.Client{Timeout: timeout}
}

func (this *scraper) get(endPoint, source string) []byte {
	req, err := http.NewRequest("GET", endPoint, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	query := req.URL.Query()
	query.Set("order", "asc")
	query.Set("sort", "votes")
	query.Set("site", source)
	query.Set("filter", "withbody")
	req.URL.RawQuery = query.Encode()
	resp, err := this.client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	return respBody
}

func (this *scraper) getQuestion(id uint, source string) *question {
	endPoint := fmt.Sprintf(
		"https://api.stackexchange.com/2.3/questions/%v", id)
	resp := this.get(endPoint, source)
	qst := new(question)
	err := json.Unmarshal(resp, qst)
	if err != nil {
		log.Println(err)
		return nil
	}
	return qst
}

func (this *scraper) getAnswer(id uint, source string) *answer {
	endPoint := fmt.Sprintf(
		"https://api.stackexchange.com/2.3/questions/%v/answers", id)
	resp := this.get(endPoint, source)
	ans := new(answer)
	err := json.Unmarshal(resp, ans)
	if err != nil {
		log.Println(err)
		return nil
	}
	return ans
}
