package textbooks

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adelrosarioh/tarantulas/utils"
	"github.com/gocolly/colly"
	"github.com/zolamk/colly-mongo-storage/colly/mongo"
)

func Run() error {

	name := "textbooks"
	mongoURI := os.Getenv("mongoURI")

	storage := &mongo.Storage{
		Database: name,
		URI:      mongoURI,
	}

	root := colly.NewCollector(colly.AllowURLRevisit())

	if err := root.SetStorage(storage); err != nil {
		return err
	}

	doc := root.Clone()
	doc.AllowURLRevisit = false

	root.OnHTML("section.content", func(e *colly.HTMLElement) {
		e.ForEach(".entry .themeform .filetitle a[href]", func(_ int, el *colly.HTMLElement) {
			doc.Visit(el.Attr("href"))
		})
	})

	doc.OnResponse(func(r *colly.Response) {
		fmt.Println("Downloading", r.FileName())

		tempFile, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("*_%s", r.FileName()))

		if err != nil {
			panic(err)
		}

		defer tempFile.Close()
		defer os.Remove(tempFile.Name())

		r.Save(tempFile.Name())

		utils.UploadToS3(name, tempFile)
	})

	root.Visit("https://proconsumidor.gob.do/monitoreo-de-libros-de-textos/")

	return nil
}
