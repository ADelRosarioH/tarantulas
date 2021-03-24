package transportationfees

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
	province       string
	town           string
	route          string
	company        string
	phoneNumber    string
	representative string
	stop           string
	line           string
	price          string
	currency       string
	publishedAt    string
	updatedAt      string
}

func Run() error {

	name := "transportationfees"
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

		publishedAt := e.ChildText("div.container div.impre div#fecha")

		e.ForEach("div#productos div.impre center table.table-striped tr", func(_ int, el *colly.HTMLElement) {
			r := record{
				publishedAt: publishedAt,
				updatedAt:   time.Now().UTC().String(),
			}

			r.province = el.ChildText("td:nth-child(1)")
			r.town = el.ChildText("td:nth-child(2)")
			r.route = el.ChildText("td:nth-child(3)")
			r.company = el.ChildText("td:nth-child(4)")
			r.phoneNumber = el.ChildText("td:nth-child(5)")
			r.representative = el.ChildText("td:nth-child(6)")
			r.stop = el.ChildText("td:nth-child(7)")
			r.line = el.ChildText("td:nth-child(8)")
			r.price = el.ChildText("td:nth-child(9)")
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

	root.Visit("https://proconsumidor.gob.do/precio-de-pasajes-autobus.php")

	return nil
}
