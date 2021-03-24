package flowers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gocolly/colly"
	"github.com/zolamk/colly-mongo-storage/colly/mongo"
)

type record struct {
	description string
	unit        string
	vendor      string
	price       string
	currency    string
	publishedAt string
	updatedAt   string
}

func Run() error {

	name := "flowers"
	mongoURI := os.Getenv("mongoURI")

	storage := &mongo.Storage{
		Database: name,
		URI:      mongoURI,
	}

	root := colly.NewCollector(colly.AllowURLRevisit())

	if err := root.SetStorage(storage); err != nil {
		return err
	}

	root.OnHTML("body.container", func(e *colly.HTMLElement) {
		records := []record{}

		publishedAt := e.ChildText("div.container div.impre p")

		e.ForEach("div#productos div.impre center table.table-striped tr", func(_ int, el *colly.HTMLElement) {
			r := record{
				publishedAt: publishedAt,
				updatedAt:   time.Now().UTC().String(),
			}

			r.vendor = el.ChildText("td:nth-child(1)")
			r.description = el.ChildText("td:nth-child(2)")
			r.unit = el.ChildText("td:nth-child(3)")
			r.price = el.ChildText("td:nth-child(4)")
			r.currency = "DOP"

			records = append(records, r)
		})

		tempFile, err := ioutil.TempFile(name, fmt.Sprintf("%s-%s.json", name, time.Now().UTC().Format("2006-01-02")))

		if err != nil {
			panic(err)
		}

		defer os.Remove(tempFile.Name())

		encoder := json.NewEncoder(tempFile)
		encoder.Encode(records)
	})

	root.Visit("http://proconsumidor.gob.do/precios-de-flores.php")

	return nil
}
