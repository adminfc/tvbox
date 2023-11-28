package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var reList = []string{
	"https://ghproxy.com/https://raw.githubusercontent.com",
	"https://raw.fgit.cf",
	"https://gcore.jsdelivr.net/gh",
	"https://raw.iqiq.io",
	"https://github.moeyy.xyz/https://raw.githubusercontent.com",
	"https://fastly.jsdelivr.net/gh",
}
var reRawList = []bool{
	false, false, true, false, false, true,
}

func main() {
	urlJson, err := ioutil.ReadFile("./url.json")
	if err != nil {
		panic(err)
	}

	var data []map[string]string
	if err := json.Unmarshal(urlJson, &data); err != nil {
		panic(err)
	}

	var urls []map[string]string

	for _, item := range data {
		for reI := range reList {
			urlName := item["name"]
			filePath := "./tv/" + strconv.Itoa(reI) + "/" + urlName + ".json"
			fmt.Println(filePath)

			resp, err := http.Get(item["url"])
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			reqText := string(body)

			if urlName != "gaotianliuyun_0707" {
				reqText = strings.Replace(reqText, "'./", "'"+item["path"], -1)
				reqText = strings.Replace(reqText, `"./`, `"`+item["path"], -1)
			}

			if reRawList[reI] {
				reqText = strings.Replace(reqText, "/raw/", "@", -1)
			} else {
				reqText = strings.Replace(reqText, "/raw/", "/", -1)
			}

			replacements := map[string]string{
				"'https://github.com":                "'" + reList[reI],
				`"https://github.com`:                `"` + reList[reI],
				"'https://raw.githubusercontent.com": "'" + reList[reI],
				`"https://raw.githubusercontent.com`: `"` + reList[reI],
			}
			for old, new := range replacements {
				reqText = strings.Replace(reqText, old, new, -1)
			}

			err = ioutil.WriteFile(filePath, []byte(reqText), 0644)
			if err != nil {
				panic(err)
			}

			// 收集文件信息
			fileInfo := map[string]string{
				"url":  filePath,
				"name": urlName,
			}
			urls = append(urls, fileInfo)
		}
	}

	// 生成包含所有文件信息的结构体
	urlsData := map[string][]map[string]string{"urls": urls}

	// 将文件信息结构体转换为 JSON
	jsonData, err := json.Marshal(urlsData)
	if err != nil {
		panic(err)
	}

	// 将 JSON 写入到文件
	outputFilePath := "./output.json"
	err = ioutil.WriteFile(outputFilePath, jsonData, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("JSON 文件已生成！")
}
