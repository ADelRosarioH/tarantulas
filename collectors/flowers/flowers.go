package flowers

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

	storage := utils.Storage(name)

	root := colly.NewCollector(colly.AllowURLRevisit())

	root.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	extensions.RandomUserAgent(root)

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

	err := root.Visit("https://proconsumidor.gob.do/precios-floristerias/")

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
