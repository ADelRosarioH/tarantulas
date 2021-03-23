package collectors

import (
	"fmt"
	"os"

	"github.com/gocolly/colly"
	"github.com/zolamk/colly-mongo-storage/colly/mongo"
)

type BasicBasket struct{}

func (b *BasicBasket) Run() error {

	name := "basicbasket"
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

	root.OnHTML("#menu-item-7437 ul li a[href]", func(e *colly.HTMLElement) {
		fmt.Println("Root", e.Attr("href"))
		page.Visit(e.Attr("href"))
	})

	page.OnHTML(".filetitle a[href]", func(e *colly.HTMLElement) {
		fmt.Println("Doc", e.Attr("href"))
		doc.Visit(e.Attr("href"))
	})

	doc.OnResponse(func(r *colly.Response) {
		fmt.Println("Downloading", r.FileName())
		r.Save(r.FileName())
	})

	root.Visit("https://proconsumidor.gob.do/")

	return nil
}
