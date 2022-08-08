package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const _fileName = "/home/meng/文档/myshare/colly/olx/record_phone.csv"

// const fileName = "record_01.csv"

var _urlField = 1

func getHighlight(ch chan []string) {
	f, err := os.Open(_fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	i := 0
	for {
		i++
		row, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("read row error:%v\n", err)
		}
		if i <= 1 {
			continue
		}

		row = append(row, reqHighlight(row[1]))

		ch <- row
	}
}

func reqHighlight(_url string) string {
	client := http.Client{Timeout: time.Second * 20}
	req, err := http.NewRequest("GET", _url, nil) //GET大写
	if err != nil {
		log.Printf("req header err:%+v\n", err)
		return ""
	}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("req header err:%+v\n", err)
		return ""
	}
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return ""
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("read body error:%+v\n", err)
		return ""
	}
	fmt.Println("url", _url)
	// Find the review items
	return doc.Find("label._10wN3 span").Text()
}
