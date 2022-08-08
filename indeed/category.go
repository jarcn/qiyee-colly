package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// 获取分类

func getCategory() []string {
	cateGoryURL := `https://autocomplete.indeed.com/api/v0/suggestions/what?country=ID&language=in&query=&count=3000`
	resp, err := http.Get(cateGoryURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return strings.Split(strings.Trim(strings.Trim(string(data), "["), "]"), ",")
}
