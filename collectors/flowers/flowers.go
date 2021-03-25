package flowers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/adelrosarioh/tarantulas/utils"
	"github.com/gocolly/colly"
	"github.com/zolamk/colly-mongo-storage/colly/mongo"
)

type record struct {
	Description string
	Unit        string
	Vendor      string
	Price       string
	Currency    string
	PublishedAt string
	UpdatedAt   string
}

func Run() error {

	name := "flowers"
	AWS_DOCUMENTDB_URI := os.Getenv("AWS_DOCUMENTDB_URI")

	storage := &mongo.Storage{
		Database: name,
		URI:      AWS_DOCUMENTDB_URI,
	}

	root := colly.NewCollector(colly.AllowURLRevisit())

	if err := root.SetStorage(storage); err != nil {
		return err
	}

	root.OnHTML(".jumbotron", func(e *colly.HTMLElement) {
		records := []record{}

		publishedAt := e.ChildText("div.container div.impre p")

		e.ForEach("div#productos div.impre center table.table-striped tr", func(_ int, el *colly.HTMLElement) {
			r := record{
				PublishedAt: publishedAt,
				UpdatedAt:   time.Now().UTC().String(),
			}

			r.Vendor = el.ChildText("td:nth-child(1)")
			r.Description = el.ChildText("td:nth-child(2)")
			r.Unit = el.ChildText("td:nth-child(3)")
			r.Price = el.ChildText("td:nth-child(4)")
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
		encoder.Encode(records)

		if err := utils.UploadToS3(name, tempFile); err != nil {
			panic(err)
		}
	})

	root.Visit("http://proconsumidor.gob.do/precios-de-flores.php")

	return nil
}
