package httpRequest

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type HttpRequest struct {
	Timeout time.Duration
	Url     string
}

func EncodeUrl(url string) string {
	if !strings.HasPrefix(url, "http") {
		return "http://" + url
	}

	return url
}

func (r *HttpRequest) ApiGet(params map[string]string) (*http.Response, error) {
	client := &http.Client{Timeout: r.Timeout}
	Url := EncodeUrl(r.Url)

	var queryArray []string
	for key, value := range params {
		queryArray = append(queryArray, key+"="+value)
	}

	//no use url.QueryEscape(strings.Join(queryArray, "&"))
	if len(params) != 0 {
		Url += "?" + strings.Join(queryArray, "&")
	}
	log.Printf("Req url:%s by get method\n", Url)

	return client.Get(Url)
}

func (r *HttpRequest) ApiPost(params map[string]string) (*http.Response, error) {
	client := &http.Client{Timeout: r.Timeout}

	log.Printf("Req url:%s, body:%v by post method\n", r.Url, params)
	body, err := json.Marshal(params)
	if err != nil {
		log.Printf("Json marshal failed %v\n", params)
		return nil, err
	}

	return client.Post(EncodeUrl(r.Url),
		"application/json; charset=utf-8",
		strings.NewReader(bytes.NewBuffer(body).String()))
}
