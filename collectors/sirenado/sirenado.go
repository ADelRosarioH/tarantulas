package sirenado

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/adelrosarioh/tarantulas/utils"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

type record struct {
	Source           string
	Description      string
	Category         string
	ImageUrl         string
	CurrencyAndPrice string
	CreatedAt        string
}

func Run() error {

	name := "sirenado"
	rootUrl := "https://sirena.do"

	storage := utils.Storage(name)

	root := colly.NewCollector(colly.AllowURLRevisit())

	root.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	extensions.RandomUserAgent(root)

	if err := root.SetStorage(storage); err != nil {
		return err
	}

	categories := root.Clone()

	page := root.Clone()

	root.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach("ul.navbar-nav .dropdown-menu .dropdown-item", func(_ int, el *colly.HTMLElement) {
			partialUrl := el.Attr("href")
			categoryUrl := fmt.Sprintf("%s%s", rootUrl, partialUrl)
			fmt.Println("Visiting", categoryUrl)
			categories.Visit(categoryUrl)
		})
	})

	categories.OnHTML("div.col-12.col-sm-8.shelf-info-txt span a[href]", func(e *colly.HTMLElement) {
		fmt.Println("Category", e.Attr("href"))
		pageUrl := fmt.Sprintf("%s%s", rootUrl, e.Attr("href"))
		page.Visit(pageUrl)
	})

	page.OnHTML("body", func(e *colly.HTMLElement) {
		records := []record{}
		category := e.ChildText(".row .page-title a")

		e.ForEach(".item-card", func(_ int, el *colly.HTMLElement) {
			r := record{
				CreatedAt: time.Now().UTC().String(),
			}

			r.Source = fmt.Sprintf("%s%s", rootUrl, el.ChildAttr(".item-pic a", "href"))
			r.Description = el.ChildText(".item-title span")
			r.Category = el.ChildText(".item-category")
			r.CurrencyAndPrice = el.ChildText(".item-price")
			r.ImageUrl = el.ChildAttr(".item-pic .lazy", "data-src")

			records = append(records, r)
		})

		tempFile, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("%s-%s-%s.*.json", name, category, time.Now().UTC().Format("2006-01-02")))

		if err != nil {
			panic(err)
		}

		defer tempFile.Close()
		defer os.Remove(tempFile.Name())

		encoder := json.NewEncoder(tempFile)

		fmt.Printf("Encoding %v record(s).\n", len(records))

		encoder.Encode(records)

		fmt.Printf("Starting %s upload to S3 bucket.\n", tempFile.Name())

		s3FileKey, err := utils.UploadToS3(name, tempFile)

		if err != nil {
			panic(err)
		}

		textMessage := fmt.Sprintf("%s collector uploaded %s to S3", name, s3FileKey)

		if err := utils.Notify(textMessage); err != nil {
			log.Fatal(err)
		}
	})

	err := root.Visit("https://sirena.do/site/home")

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
