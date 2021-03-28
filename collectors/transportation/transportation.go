package transportation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/adelrosarioh/tarantulas/utils"
	"github.com/gocolly/colly"
)

type record struct {
	Province       string
	Town           string
	Route          string
	Company        string
	PhoneNumber    string
	Representative string
	Stop           string
	Line           string
	Price          string
	Currency       string
	PublishedAt    string
	UpdatedAt      string
}

func Run() error {

	name := "transportation"

	storage := utils.Storage(name)

	root := colly.NewCollector(colly.AllowURLRevisit())

	if err := root.SetStorage(storage); err != nil {
		return err
	}

	root.OnHTML("body", func(e *colly.HTMLElement) {
		records := []record{}

		publishedAt := e.ChildText("div.container p")

		e.ForEach("div#productos div.impre center table.table-striped tr", func(_ int, el *colly.HTMLElement) {
			r := record{
				PublishedAt: publishedAt,
				UpdatedAt:   time.Now().UTC().String(),
			}

			r.Province = el.ChildText("td:nth-child(1)")
			r.Town = el.ChildText("td:nth-child(2)")
			r.Route = el.ChildText("td:nth-child(3)")
			r.Company = el.ChildText("td:nth-child(4)")
			r.PhoneNumber = el.ChildText("td:nth-child(5)")
			r.Representative = el.ChildText("td:nth-child(6)")
			r.Stop = el.ChildText("td:nth-child(7)")
			r.Line = el.ChildText("td:nth-child(8)")
			r.Price = el.ChildText("td:nth-child(9)")
			r.Currency = "DOP"

			records = append(records, r)
		})

		tempFile, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("%s-%s.*.json", name, time.Now().UTC().Format("2006-01-02")))

		if err != nil {
			panic(err)
		}

		defer tempFile.Close()
		defer os.Remove(tempFile.Name())

		encoder := json.NewEncoder(tempFile)

		fmt.Printf("Encoding %v record(s).", len(records))

		encoder.Encode(records)

		fmt.Printf("Starting %s upload to S3 bucket.\n", tempFile.Name())

		if err := utils.UploadToS3(name, tempFile); err != nil {
			panic(err)
		}
	})

	root.Visit("https://proconsumidor.gob.do/precios-pasajes/")

	return nil
}
