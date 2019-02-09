package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/turnps/stockinfo/tools"
)

type StockData struct {
	str0 string //股票名稱
	str1 string //產業別
	str2 string //上市/上櫃
	str3 string //成立日期
	str4 string //上市日期
	str5 string //資本額
	str6 string //目前市值
	str7 string //發行股數
	str8 string //股票代號
}

func main() {
	var url string
	var codes []string
	var webContent *goquery.Document
	var stockData []StockData

	codes = csvCodes("stock.csv")
	for _, stockid := range codes {
		url = "https://goodinfo.tw/StockInfo/BasicInfo.asp?STOCK_ID=" + stockid
		webContent = tools.Fetch(url)
		stockData = append(stockData, setStockData(webContent))
	}
	tools.WriteCsv(setCsvData(stockData))
	fmt.Println("Press the Enter Key to terminate the console screen!")
	fmt.Scanln() // wait for Enter Key
}

func setStockData(doc *goquery.Document) StockData {
	data := StockData{}

	topicsSelection := doc.Find(".solid_1_padding_2_6_tbl tr")
	topicNode := goquery.NewDocumentFromNode(topicsSelection.Get(1)).Find("td").Get(3)
	data.str0 = goquery.NewDocumentFromNode(topicNode).Text()
	topicNode = goquery.NewDocumentFromNode(topicsSelection.Get(1)).Find("td").Get(1)
	data.str8 = goquery.NewDocumentFromNode(topicNode).Text()
	topicNode = goquery.NewDocumentFromNode(topicsSelection.Get(2)).Find("td").Get(1)
	data.str1 = goquery.NewDocumentFromNode(topicNode).Text()
	topicNode = goquery.NewDocumentFromNode(topicsSelection.Get(2)).Find("td").Get(3)
	data.str2 = goquery.NewDocumentFromNode(topicNode).Text()
	topicNode = goquery.NewDocumentFromNode(topicsSelection.Get(5)).Find("td").Get(1)
	data.str3 = goquery.NewDocumentFromNode(topicNode).Text()
	topicNode = goquery.NewDocumentFromNode(topicsSelection.Get(6)).Find("td").Get(1)
	data.str4 = goquery.NewDocumentFromNode(topicNode).Text()
	topicNode = goquery.NewDocumentFromNode(topicsSelection.Get(7)).Find("td").Get(1)
	data.str5 = goquery.NewDocumentFromNode(topicNode).Text()
	topicNode = goquery.NewDocumentFromNode(topicsSelection.Get(8)).Find("td").Get(1)
	data.str6 = goquery.NewDocumentFromNode(topicNode).Text()
	topicNode = goquery.NewDocumentFromNode(topicsSelection.Get(9)).Find("td").Get(1)
	data.str7 = goquery.NewDocumentFromNode(topicNode).Text()
	fmt.Println(data.str8, data)
	return data
}

func csvCodes(name string) []string {
	csvFile, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()
	csvReader := csv.NewReader(csvFile)
	//csvReader.Read()
	codes := []string{}
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err) // or handle it another way
		}
		// use the `row` here
		if row[0] == "" {
			continue
		}
		codes = append(codes, row[0])
	}
	return codes
}

func setCsvData(data []StockData) [][]string {
	var high52 string
	var low52 string
	var price string
	var rows = [][]string{
		{"", "名稱", "產業", "52週最高", "52週最低", "最新價", "上市/上櫃", "成立日期", "上市日期", "資本額", "目前市值", "發行股數"},
	}
	for _, s := range data {
		if s.str2 == "上市" {
			high52 = "=GoogleFinance(CONCATENATE(\"TPE:\"," + s.str8 + "),\"high52\")"
			low52 = "=GoogleFinance(CONCATENATE(\"TPE:\"," + s.str8 + "),\"low52\")"
			price = "=GoogleFinance(CONCATENATE(\"TPE:\"," + s.str8 + "),\"price\")"
		} else {
			high52 = ""
			low52 = ""
			price = ""
		}
		stockData := []string{s.str8, s.str0, s.str1, high52, low52, price, s.str2, s.str3, s.str4, s.str5, s.str6, s.str7}
		rows = append(rows, stockData)
	}
	return rows
}
