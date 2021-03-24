package hardware

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gocolly/colly"
	"github.com/zolamk/colly-mongo-storage/colly/mongo"
)

func Run() error {

	name := "hardware"
	mongoURI := os.Getenv("mongoURI")

	storage := &mongo.Storage{
		Database: name,
		URI:      mongoURI,
	}

	root := colly.NewCollector(colly.AllowURLRevisit())

	if err := root.SetStorage(storage); err != nil {
		return err
	}

	page := root.Clone()
	page.AllowURLRevisit = false

	doc := page.Clone()

	root.OnHTML("#menu-item-9868 ul li a[href]", func(e *colly.HTMLElement) {
		fmt.Println("Root", e.Attr("href"))
		page.Visit(e.Attr("href"))
	})

	page.OnHTML(".filetitle a[href]", func(e *colly.HTMLElement) {
		fmt.Println("Doc", e.Attr("href"))
		doc.Visit(e.Attr("href"))
	})

	doc.OnResponse(func(r *colly.Response) {
		fmt.Println("Downloading", r.FileName())

		tempFile, err := ioutil.TempFile(name, r.FileName())

		if err != nil {
			panic(err)
		}

		defer os.Remove(tempFile.Name())

		r.Save(tempFile.Name())
	})

	root.Visit("https://proconsumidor.gob.do/")

	return nil
}
