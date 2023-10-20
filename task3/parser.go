package main

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
)

func main() {

	file, err := os.Create("output.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Rank", "Influencer", "Category", "Followers", "Country", "Eng. (Auth.)", "Eng. (Avg.)"}
	writer.Write(headers)

	url := "https://hypeauditor.com/top-instagram-all-russia/"

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatalf("HTTP запрос вернул статус: %d", response.StatusCode)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	table := document.Find(".table")

	table.Find(".row .row__top").Each(func(index int, row *goquery.Selection) {

		rank := row.Find(".row-cell.rank span[data-v-2e6a30b8]").Text()
		influencer := row.Find(".row-cell.contributor .contributor__title").Text()
		category := row.Find(".row-cell.category").Text()
		followers := row.Find(".row-cell.subscribers").Text()
		country := row.Find(".row-cell.audience").Text()
		authEng := row.Find(".row-cell.authentic").Text()
		avgEng := row.Find(".row-cell.engagement").Text()

		data := []string{rank, influencer, category, followers, country, authEng, avgEng}
		writer.Write(data)
	})
	fmt.Println("Парсинг завершен. Результат сохранен в output.csv")
}
