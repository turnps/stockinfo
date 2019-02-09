package tools

import (
	"crypto/tls"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

//寫CSV檔，參數傳入[][]string
func WriteCsv(rows [][]string) {
	csvFile, err := os.Create("finish.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	csvFile.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM，防止中文乱码
	csvWriter := csv.NewWriter(csvFile)
	//rows := f(data)
	err = csvWriter.WriteAll(rows)
	if err != nil {
		fmt.Printf("error (%v)", err)
		return
	}
}

//讀取url原始碼，參數傳入string，傳回*goquery.Document物件
func Fetch(url string) *goquery.Document {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}
